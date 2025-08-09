package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/cbhcbhcbh/Quantum/pkg/common/known"
	goJwt "github.com/golang-jwt/jwt/v5"
)

var (
	JwtSecret           string
	JwtExpirationSecond int64
)

type JWTClaims struct {
	ChannelID uint64
	goJwt.RegisteredClaims
}

type AuthPayload struct {
	AccessToken string
}

type AuthResponse struct {
	ChannelID uint64
	Expired   bool
}

func Auth(authPayload *AuthPayload) (*AuthResponse, error) {
	token, err := parseToken(authPayload.AccessToken)
	if err != nil {
		if errors.Is(err, goJwt.ErrTokenExpired) {
			return &AuthResponse{
				Expired: true,
			}, nil
		}
		return nil, known.ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, known.ErrInvalidToken
	}

	return &AuthResponse{
		ChannelID: claims.ChannelID,
		Expired:   false,
	}, nil
}

func NewJWT(channelID uint64) (string, error) {
	expiresAt := time.Now().Add(time.Duration(JwtExpirationSecond) * time.Second)
	jwtClaims := &JWTClaims{
		ChannelID: channelID,
		RegisteredClaims: goJwt.RegisteredClaims{
			ExpiresAt: goJwt.NewNumericDate(expiresAt),
		},
	}
	token := goJwt.NewWithClaims(goJwt.SigningMethodHS256, jwtClaims)
	accessToken, err := token.SignedString([]byte(JwtSecret))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func parseToken(accessToken string) (*goJwt.Token, error) {
	return goJwt.ParseWithClaims(accessToken, &JWTClaims{}, func(token *goJwt.Token) (any, error) {
		if _, ok := token.Method.(*goJwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JwtSecret), nil
	})
}
