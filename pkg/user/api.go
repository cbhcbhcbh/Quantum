package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *HttpServer) CreateLocalUser(c *gin.Context) {
	// TODO：Implement create user logic
}

func (h *HttpServer) OAuthGoogleLogin(c *gin.Context) {
	oauthState := ""
	u := h.googleOauthConfig.AuthCodeURL(oauthState)
	c.Redirect(http.StatusTemporaryRedirect, u)
}

func (h *HttpServer) OAuthGoogleCallback(c *gin.Context) {
	code := c.Request.FormValue("code")
	token, err := h.googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		h.logger.Error("failed to exchange token", zap.Error(err))
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	fmt.Println("Access Token:", token.AccessToken)
}
