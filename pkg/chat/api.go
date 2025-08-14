package chat

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/cbhcbhcbh/Quantum/pkg/common/domain"
	"github.com/cbhcbhcbh/Quantum/pkg/common/jwt"
	"github.com/cbhcbhcbh/Quantum/pkg/common/known"
	"github.com/cbhcbhcbh/Quantum/pkg/common/response"
	"github.com/gin-gonic/gin"
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
}

func (h *HttpServer) DeleteChannel(c *gin.Context) {

}
