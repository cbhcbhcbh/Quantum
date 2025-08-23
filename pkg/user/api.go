package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cbhcbhcbh/Quantum/pkg/common/cookie"
	"github.com/cbhcbhcbh/Quantum/pkg/common/domain"
	"github.com/cbhcbhcbh/Quantum/pkg/common/known"
	"github.com/cbhcbhcbh/Quantum/pkg/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *HttpServer) CreateLocalUser(c *gin.Context) {
	var createLocalUserReq domain.CreateLocalUserRequest
	if err := c.ShouldBindJSON(&createLocalUserReq); err != nil {
		response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
		return
	}
	user, err := h.userSvc.CreateUser(c.Request.Context(), &User{
		Name:     createLocalUserReq.Name,
		AuthType: LocalAuth,
	})
	if err != nil {
		h.logger.Error(err.Error())
		response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
		return
	}
	sid, err := h.userSvc.SetUserSession(c.Request.Context(), user.ID)
	if err != nil {
		h.logger.Error(err.Error())
		response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
		return
	}
	cookie.SetAuthCookie(c, known.SessionIdCookieName, sid, h.authCookieConfig.MaxAge, h.authCookieConfig.Path, h.authCookieConfig.Domain)

	c.JSON(http.StatusCreated, &domain.UserPresenter{
		ID:   strconv.FormatUint(user.ID, 10),
		Name: user.Name,
	})
}

func (h *HttpServer) GetUser(c *gin.Context) {
	_, ok := c.Request.Context().Value(known.UserKey).(uint64)
	if !ok {
		response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
		return
	}
	var getUserReq domain.GetUserRequest
	if err := c.ShouldBindQuery(&getUserReq); err != nil {
		response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
		return
	}
	userID, err := strconv.ParseUint(getUserReq.Uid, 10, 64)
	if err != nil {
		response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
		return
	}
	user, err := h.userSvc.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, known.ErrUserNotFound) {
			response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
			return
		}
		h.logger.Error(err.Error())
		response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
		return
	}
	c.JSON(http.StatusOK, &domain.UserPresenter{
		ID:      strconv.FormatUint(userID, 10),
		Name:    user.Name,
		Picture: user.Picture,
	})
}

func (h *HttpServer) GetUserMe(c *gin.Context) {
	userID, ok := c.Request.Context().Value(known.UserKey).(uint64)
	if !ok {
		response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
		return
	}
	user, err := h.userSvc.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, known.ErrUserNotFound) {
			response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
			return
		}
		h.logger.Error(err.Error())
		response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
		return
	}
	c.JSON(http.StatusOK, &domain.UserPresenter{
		ID:      strconv.FormatUint(userID, 10),
		Name:    user.Name,
		Picture: user.Picture,
	})
}

func (h *HttpServer) OAuthGoogleLogin(c *gin.Context) {
	oauthState, err := cookie.GenerateStateOauthCookie(c, h.oauthCookieConfig.MaxAge, h.oauthCookieConfig.Path, h.oauthCookieConfig.Domain)
	if err != nil {
		h.logger.Error("failed to generate oauth state cookie", zap.Error(err))
		response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
		return
	}
	u := h.googleOauthConfig.AuthCodeURL(oauthState)
	c.Redirect(http.StatusTemporaryRedirect, u)
}

func (h *HttpServer) OAuthGoogleCallback(c *gin.Context) {
	oauthState, err := cookie.GetCookie(c, known.OAuthStateCookieName)
	if err != nil {
		h.logger.Error("failed to get oauth state cookie", zap.Error(err))
		response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
		return
	}
	if c.Query("state") != oauthState {
		h.logger.Error("invalid oauth state", zap.String("expected", oauthState), zap.String("received", c.Query("state")))
		response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
		return
	}

	code := c.Request.FormValue("code")
	token, err := h.googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		h.logger.Error("failed to exchange token", zap.Error(err))
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	fmt.Println("Access Token:", token.AccessToken)

	googleUser, err := h.userSvc.GetGoogleUser(c.Request.Context(), token.AccessToken)
	if err != nil {
		h.logger.Error("failed to get google user", zap.Error(err))
		response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
		return
	}

	user, err := h.userSvc.GetOrCreateUserByOAuth(c.Request.Context(), &User{
		Email:    googleUser.Email,
		Name:     googleUser.Name,
		Picture:  googleUser.Picture,
		AuthType: GoogleAuth,
	})
	if err != nil {
		h.logger.Error(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	sid, err := h.userSvc.SetUserSession(c.Request.Context(), user.ID)
	if err != nil {
		h.logger.Error(err.Error())
		response.ErrorResponse(http.StatusInternalServerError, known.ErrServer.Error()).SetHttpCode(http.StatusInternalServerError).WriteTo(c)
		return
	}
	cookie.SetAuthCookie(c, known.SessionIdCookieName, sid, h.authCookieConfig.MaxAge, h.authCookieConfig.Path, h.authCookieConfig.Domain)

	c.Redirect(http.StatusTemporaryRedirect, "/")
}
