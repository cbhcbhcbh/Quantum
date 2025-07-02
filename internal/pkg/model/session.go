package model

const (
	SessionStatusOk = 0

	TopStatus  = 0
	GROUP_TYPE = 2
)

type SessionM struct {
	ModelID
	ModelTimestamps
	FormID      int64  `gorm:"column:form_id" json:"form_id"`             // sender id
	ToID        int64  `gorm:"column:to_id" json:"to_id"`                 // receiver id
	TopStatus   int16  `gorm:"column:top_status" json:"top_status"`       // 0 no 1 yes
	TopTime     string `gorm:"column:top_time" json:"top_time,omitempty"` // top time
	Note        string `gorm:"column:note" json:"note,omitempty"`         // remark
	ChannelType int16  `gorm:"column:channel_type" json:"channel_type"`   // 0 private chat 1 group chat
	Name        string `gorm:"column:name" json:"name,omitempty"`         // session name
	Avatar      string `gorm:"column:avatar" json:"avatar,omitempty"`     // session avatar
	Status      int16  `gorm:"column:status" json:"status"`               // 0 normal 1 disabled
	GroupID     int64  `gorm:"column:group_id" json:"group_id,omitempty"` // group id
}
