package chat

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/cbhcbhcbh/Quantum/pkg/common/domain"
	"github.com/cbhcbhcbh/Quantum/pkg/common/jwt"
	"github.com/cbhcbhcbh/Quantum/pkg/common/known"
	"github.com/cbhcbhcbh/Quantum/pkg/common/response"
	"github.com/cbhcbhcbh/Quantum/pkg/common/util"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func (h *HttpServer) StartChat(c *gin.Context) {
	uid := c.Query("uid")
	userID, err := strconv.ParseUint(uid, 10, 64)
	if err != nil {
		response.ErrorResponse(http.StatusBadRequest, known.ErrInvalidParam.Error()).SetHttpCode(http.StatusBadRequest).WriteTo(c)
		return
	}

	_, err = h.userSvc.GetUser(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, known.ErrUserNotFound) {
			response.ErrorResponse(http.StatusNotFound, known.ErrUserNotFound.Error()).SetHttpCode(http.StatusNotFound).WriteTo(c)
			return
		}
		h.logger.Error(err.Error())
		response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
		return
	}

	accessToken := c.Query("access_token")
	authResult, err := jwt.Auth(&jwt.AuthPayload{
		AccessToken: accessToken,
	})
	if err != nil {
		response.ErrorResponse(http.StatusUnauthorized, known.ErrUnauthorized.Error()).SetHttpCode(http.StatusUnauthorized).WriteTo(c)
		return
	}
	if authResult.Expired {
		h.logger.Error(known.ErrTokenExpired.Error())
		response.ErrorResponse(http.StatusUnauthorized, known.ErrTokenExpired.Error()).SetHttpCode(http.StatusUnauthorized).WriteTo(c)
	}
	channelID := authResult.ChannelID
	exist, err := h.userSvc.IsChannelUserExist(c.Request.Context(), channelID, userID)
	if err != nil {
		h.logger.Error(err.Error())
		response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
		return
	}
	if !exist {
		response.ErrorResponse(http.StatusNotFound, known.ErrChannelOrUserNotFound.Error()).SetHttpCode(http.StatusNotFound).WriteTo(c)
		return
	}

	if err := h.mc.HandleRequest(c.Writer, c.Request); err != nil {
		h.logger.Error("upgrade websocket error: " + err.Error())
		response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
		return
	}
}

func (h *HttpServer) ForwardAuth(c *gin.Context) {
	channelID, ok := c.Request.Context().Value(known.ChannelKey).(uint64)
	if !ok {
		response.ErrorResponse(http.StatusUnauthorized, known.ErrUnauthorized.Error()).SetHttpCode(http.StatusUnauthorized).WriteTo(c)
		return
	}
	c.Writer.Header().Set(known.ChannelIdHeader, strconv.FormatUint(channelID, 10))
	c.Status(http.StatusOK)
}

func (h *HttpServer) GetChannelUsers(c *gin.Context) {
	channelID, ok := c.Request.Context().Value(known.ChannelKey).(uint64)
	if !ok {
		response.ErrorResponse(http.StatusUnauthorized, known.ErrUnauthorized.Error()).SetHttpCode(http.StatusUnauthorized).WriteTo(c)
		return
	}
	userIDs, err := h.userSvc.GetChannelUserIDs(c.Request.Context(), channelID)
	if err != nil {
		h.logger.Error(err.Error())
		response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
		return
	}
	userIDsPresenter := []string{}
	for _, userID := range userIDs {
		userIDsPresenter = append(userIDsPresenter, strconv.FormatUint(userID, 10))
	}
	c.JSON(http.StatusOK, &domain.UserIDsPresenter{
		UserIDs: userIDsPresenter,
	})
}

func (h *HttpServer) GetOnlineUsers(c *gin.Context) {
	channelID, ok := c.Request.Context().Value(known.ChannelKey).(uint64)
	if !ok {
		response.ErrorResponse(http.StatusUnauthorized, known.ErrUnauthorized.Error()).SetHttpCode(http.StatusUnauthorized).WriteTo(c)
		return
	}
	userIDs, err := h.userSvc.GetOnlineUserIDs(c.Request.Context(), channelID)
	if err != nil {
		h.logger.Error(err.Error())
		response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
		return
	}
	userIDsPresenter := []string{}
	for _, userID := range userIDs {
		userIDsPresenter = append(userIDsPresenter, strconv.FormatUint(userID, 10))
	}
	c.JSON(http.StatusOK, &domain.UserIDsPresenter{
		UserIDs: userIDsPresenter,
	})
}

func (h *HttpServer) ListMessages(c *gin.Context) {
	channelID, ok := c.Request.Context().Value(known.ChannelKey).(uint64)
	if !ok {
		response.ErrorResponse(http.StatusUnauthorized, known.ErrUnauthorized.Error()).SetHttpCode(http.StatusUnauthorized).WriteTo(c)
		return
	}
	pageState := c.Query("ps")
	msgs, nextPageState, err := h.msgSvc.ListMessages(c.Request.Context(), channelID, pageState)
	if err != nil {
		h.logger.Error(err.Error())
		response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
		return
	}
	msgsPresenter := []domain.MessagePresenter{}
	for _, msg := range msgs {
		msgsPresenter = append(msgsPresenter, *msg.ToPresenter())
	}
	c.JSON(http.StatusOK, &domain.MessagesPresenter{
		Messages:      msgsPresenter,
		NextPageState: nextPageState,
	})
}

func (h *HttpServer) DeleteChannel(c *gin.Context) {
	channelID, ok := c.Request.Context().Value(known.ChannelKey).(uint64)
	if !ok {
		response.ErrorResponse(http.StatusUnauthorized, known.ErrUnauthorized.Error()).SetHttpCode(http.StatusUnauthorized).WriteTo(c)
		return
	}
	uid := c.Query("delby")
	userID, err := strconv.ParseUint(uid, 10, 64)
	if err != nil {
		response.ErrorResponse(http.StatusBadRequest, known.ErrInvalidParam.Error()).SetHttpCode(http.StatusBadRequest).WriteTo(c)
		return
	}

	exist, err := h.userSvc.IsChannelUserExist(c.Request.Context(), channelID, userID)
	if err != nil {
		h.logger.Error(err.Error())
		response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
		return
	}
	if !exist {
		response.ErrorResponse(http.StatusNotFound, known.ErrChannelOrUserNotFound.Error()).SetHttpCode(http.StatusNotFound).WriteTo(c)
		return
	}

	err = h.msgSvc.BroadcastActionMessage(c.Request.Context(), channelID, userID, domain.LeavedMessage)
	if err != nil {
		h.logger.Error(err.Error())
		response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
		return
	}
	if err := h.chanSvc.DeleteChannel(c.Request.Context(), channelID); err != nil {
		h.logger.Error(err.Error())
		response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
		return
	}
	c.JSON(http.StatusNoContent, domain.SuccessMessage{
		Message: "ok",
	})
}

func (h *HttpServer) HandleChatOnConnect(sess *melody.Session) {
	userID, err := strconv.ParseUint(sess.Request.URL.Query().Get("uid"), 10, 64)
	if err != nil {
		h.logger.Error(err.Error())
		return
	}
	accessToken := sess.Request.URL.Query().Get("access_token")
	authResult, err := jwt.Auth(&jwt.AuthPayload{
		AccessToken: accessToken,
	})
	if err != nil {
		h.logger.Error(err.Error())
	}
	if authResult.Expired {
		h.logger.Error(known.ErrTokenExpired.Error())
	}
	channelID := authResult.ChannelID
	err = h.initializeChatSession(sess, channelID, userID)
	if err != nil {
		h.logger.Error(err.Error())
		return
	}
	if err := h.msgSvc.BroadcastConnectMessage(context.Background(), channelID, userID); err != nil {
		h.logger.Error(err.Error())
		return
	}
}

func (h *HttpServer) initializeChatSession(sess *melody.Session, channelID, userID uint64) error {
	ctx := context.Background()
	if err := h.userSvc.AddOnlineUser(ctx, channelID, userID); err != nil {
		return err
	}
	if err := h.forwardSvc.RegisterChannelSession(ctx, channelID, userID, h.msgSubscriber.subscriberID); err != nil {
		return err
	}
	sess.Set(sessCidKey, channelID)
	return nil
}

func (h *HttpServer) HandleChatOnMessage(sess *melody.Session, data []byte) {
	msgPresenter, err := util.DecodeToMessagePresenter(data)
	if err != nil {
		h.logger.Error(err.Error())
		return
	}
	msg, err := msgPresenter.ToMessage(sess.Request.URL.Query().Get("access_token"))
	if err != nil {
		h.logger.Error(err.Error())
		return
	}

	switch msg.Event {
	case domain.EventText:
		if err := h.msgSvc.BroadcastTextMessage(context.Background(), msg.ChannelID, msg.UserID, msg.Payload); err != nil {
			h.logger.Error(err.Error())
		}
	case domain.EventAction:
		if err := h.msgSvc.BroadcastActionMessage(context.Background(), msg.ChannelID, msg.UserID, domain.Action(msg.Payload)); err != nil {
			h.logger.Error(err.Error())
		}
	case domain.EventSeen:
		messageID, err := strconv.ParseUint(msg.Payload, 10, 64)
		if err != nil {
			h.logger.Error(err.Error())
			return
		}
		if err := h.msgSvc.MarkMessageSeen(context.Background(), msg.ChannelID, msg.UserID, messageID); err != nil {
			h.logger.Error(err.Error())
		}
	case domain.EventFile:
		if err := h.msgSvc.BroadcastFileMessage(context.Background(), msg.ChannelID, msg.UserID, msg.Payload); err != nil {
			h.logger.Error(err.Error())
		}
	default:
		h.logger.Error("invailid event type: " + strconv.Itoa(msg.Event))
	}
}

func (h *HttpServer) HandleChatOnClose(sess *melody.Session, i int, s string) error {
	userID, err := strconv.ParseUint(sess.Request.URL.Query().Get("uid"), 10, 64)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	accessToken := sess.Request.URL.Query().Get("access_token")
	authResult, err := jwt.Auth(&jwt.AuthPayload{
		AccessToken: accessToken,
	})
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	if authResult.Expired {
		h.logger.Error(known.ErrTokenExpired.Error())
		return known.ErrTokenExpired
	}
	channelID := authResult.ChannelID
	err = h.userSvc.DeleteOnlineUser(context.Background(), channelID, userID)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	err = h.forwardSvc.RemoveChannelSession(context.Background(), channelID, userID)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	return h.msgSvc.BroadcastActionMessage(context.Background(), channelID, userID, domain.OfflineMessage)
}
