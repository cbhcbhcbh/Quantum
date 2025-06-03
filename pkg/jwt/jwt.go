package jwt

import (
	"errors"
	"time"

	goJwt "github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type JWT struct {
	SigningKey []byte
	MaxRefresh *goJwt.NumericDate
}

type CustomClaims struct {
	ID         int64  `json:"id"`
	UID        string `json:"uid"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	ExpireTime int64  `json:"expire_time"`

	goJwt.RegisteredClaims
}

func NewJWT() *JWT {
	return &JWT{
		SigningKey: []byte(viper.GetString("jwt.secret")),
		MaxRefresh: goJwt.NewNumericDate(time.Now().Add(viper.GetDuration("jwt.max-refresh"))),
	}
}

func (j *JWT) createToken(claims CustomClaims) (string, error) {
	token := goJwt.NewWithClaims(goJwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.SigningKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := goJwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *goJwt.Token) (any, error) {
		if _, ok := token.Method.(*goJwt.SigningMethodHMAC); !ok {
			return nil, goJwt.ErrTokenSignatureInvalid
		}
		return []byte(j.SigningKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (j *JWT) IssueToken(ID int64, UID string, Name string, Email string, expireAtTime int64) string {

	claims := CustomClaims{
		ID:         ID,
		UID:        UID,
		Name:       Name,
		Email:      Email,
		ExpireTime: expireAtTime,
		RegisteredClaims: goJwt.RegisteredClaims{
			Issuer:    "Quantum",
			IssuedAt:  goJwt.NewNumericDate(time.Now()),
			ExpiresAt: goJwt.NewNumericDate(time.Now().Add(viper.GetDuration("jwt.timeout"))),
			NotBefore: goJwt.NewNumericDate(time.Now()),
		},
	}

	token, err := j.createToken(claims)
	if err != nil {
		return ""
	}
	return token
}

func (j *JWT) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	if j.MaxRefresh != nil && time.Now().After(j.MaxRefresh.Time) {
		return "", errors.New("refresh period expired")
	}

	claims.IssuedAt = goJwt.NewNumericDate(time.Now())
	claims.ExpiresAt = goJwt.NewNumericDate(time.Now().Add(viper.GetDuration("jwt.timeout")))
	claims.NotBefore = goJwt.NewNumericDate(time.Now())

	return j.createToken(*claims)
}
