package known

import "errors"

const (
	XRequestIDKey = "X-Request-ID"

	XIdKey = "X-ID"

	XUidKey = "X-Uid"

	XUsernameKey = "X-Username"

	XEmailKey = "X-Email"
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
)

var (
	ChannelUsersPrefix = "rc:chanusers"
	OnlineUsersPrefix  = "rc:onlineusers"
)
