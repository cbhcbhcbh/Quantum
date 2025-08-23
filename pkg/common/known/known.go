package known

import "errors"

const (
	XRequestIDKey = "X-Request-ID"

	XIdKey = "X-ID"

	XUidKey = "X-Uid"

	XUsernameKey = "X-Username"

	XEmailKey = "X-Email"
)

const (
	OAuthStateCookieName string = "oauthstate"
	SessionIdCookieName  string = "sid"
)

type HTTPContextKey string

var (
	JWTAuthHeader                  = "Authorization"
	ChannelIdHeader                = "X-Channel-Id"
	ChannelKey      HTTPContextKey = "channel_key"
	UserKey         HTTPContextKey = "user_key"
)

var (
	ErrInvalidParam = errors.New("invalid parameter")
	ErrServer       = errors.New("server error")
	ErrUnauthorized = errors.New("unauthorized")

	ErrUserNotFound           = errors.New("error user not found")
	ErrSessionNotFound = errors.New("error session not found")
	ErrChannelOrUserNotFound  = errors.New("error channel or user not found")
	ErrExceedMessageNumLimits = errors.New("error exceed max number of messages")

	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")

	ErrInvalidOauth = errors.New("invalid oauth google state")
)

var (
	ChannelUsersPrefix = "rc:chanusers"
	OnlineUsersPrefix  = "rc:onlineusers"
)
