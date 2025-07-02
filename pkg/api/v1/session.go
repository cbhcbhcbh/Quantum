package v1

type SessionDetail struct {
	ID          int64  `json:"id" validate:"required"`
	FormID      int64  `json:"form_id" validate:"required"`
	ToID        int64  `json:"to_id" validate:"required"`
	TopStatus   int16  `json:"top_status" validate:"omitempty,oneof=0 1"` // 0 no 1 yes
	TopTime     string `json:"top_time,omitempty"`
	Note        string `json:"note,omitempty"`
	ChannelType int16  `json:"channel_type" validate:"omitempty,oneof=0 1"` // 0 private chat 1 group chat
	Name        string `json:"name,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	Status      int16  `json:"status" validate:"omitempty,oneof=0 1"` // 0 normal 1 disabled
	GroupID     int64  `json:"group_id,omitempty"`
}

type SessionForm struct {
	ToID        int64 `json:"to_id" validate:"required"`
	ChannelType int16 `json:"channel_type" validate:"omitempty,oneof=0 1"`
}

type SessionUpdateForm struct {
	TopStatus int16  `json:"top_status" validate:"required,gte=0,lte=1"`
	Note      string `json:"type"`
}

type SessionInfo struct {
	SessionId int64 `uri:"id" binding:"required"`
}