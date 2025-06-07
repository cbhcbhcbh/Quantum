package enum

const (
	ParamError = 1000
	ApiError   = 1001
	DBError    = 1002

	WsChantMessage    = 200
	VideoChantMessage = 600
	WsCreate          = 1000
	WsFriendOk        = 1001
	WsFriendError     = 1002
	WsNotFriend       = 1003
	WsPing            = 1004
	WsAck             = 1005

	WsUserOffline  = 2000
	WsUserOnline   = 2001
	WsIsUserStatus = 2002
	WsSession      = 2003

	WsGroupMessage = 3000

	PrivateMessage = 1
	GroupMessage   = 2

	TEXT         = 1
	VOICE        = 2
	FILE         = 3
	IMAGE        = 4
	LOGOUT_GROUP = 5
	JOIN_GROUP   = 6
)
