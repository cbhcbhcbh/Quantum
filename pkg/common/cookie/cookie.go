package cookie

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/cbhcbhcbh/Quantum/pkg/common/known"
)

func GenerateStateOauthCookie(c *gin.Context, maxAge int, path, domain string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate oauth state cookie error: %w", err)
	}
	state := base64.URLEncoding.EncodeToString(b)
	c.SetCookie(known.OAuthStateCookieName, state, maxAge, path, domain, false, true)
	return state, nil
}

func SetAuthCookie(c *gin.Context, cookieName, cookieValue string, maxAge int, path, domain string) {
	c.SetCookie(cookieName, url.QueryEscape(cookieValue), maxAge, path, domain, false, true)
}

func GetCookie(c *gin.Context, cookieName string) (string, error) {
	cookie, err := c.Request.Cookie(cookieName)
	if err != nil {
		return "", fmt.Errorf("get cookie %s error: %w", cookieName, err)
	}
	unescapedCookie, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return "", fmt.Errorf("unescape oauth state cookie error: %w", err)
	}
	return unescapedCookie, nil
}
