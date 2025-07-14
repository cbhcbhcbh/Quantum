package model

import (
	"encoding/json"
)

// MsgType
var (
	TEXT         = 1
	VOICE        = 2
	FILE         = 3
	IMAGE        = 4
	LOGOUT_GROUP = 5
	JOIN_GROUP   = 6
)

// group_messages
type GroupMessageM struct {
	ModelID
	ModelTimestamps
	Message         json.RawMessage `gorm:"column:message" json:"message"` // message entity (JSON)
	SendTime        int64           `gorm:"column:send_time" json:"send_time"`
	MessageID       int64           `gorm:"column:message_id" json:"message_id"`
	ClientMessageID int64           `gorm:"column:client_message_id" json:"client_message_id"`
	FormID          int64           `gorm:"column:form_id" json:"form_id"`
	GroupID         int64           `gorm:"column:group_id" json:"group_id"`
}

func (m *GroupMessageM) TableName() string {
	return "group_messages"
}

// group_offline_messages
type GroupOfflineMessageM struct {
	ModelID
	ModelTimestamps
	Message   json.RawMessage `gorm:"column:message" json:"message"` // message body (JSON)
	SendTime  int64           `gorm:"column:send_time" json:"send_time"`
	Status    int16           `gorm:"column:status" json:"status"`         // 0 not pushed 1 pushed
	ReceiveID int64           `gorm:"column:receive_id" json:"receive_id"` // receiver id
}

func (m *GroupOfflineMessageM) TableName() string {
	return "group_offline_messages"
}

// group_user_messages
type GroupUserMessageM struct {
	ModelID
	ModelTimestamps
	UserID  int64 `gorm:"column:user_id" json:"user_id"`
	GroupID int64 `gorm:"column:group_id" json:"group_id"`
	Status  int16 `gorm:"column:status" json:"status"` // 0 unread 1 read
}

func (m *GroupUserMessageM) TableName() string {
	return "group_user_messages"
}

// private messages
type MessageM struct {
	ModelID
	ModelTimestamps
	Msg      string `gorm:"column:msg" json:"msg"`
	FormID   int64  `gorm:"column:form_id" json:"form_id"`
	ToID     int64  `gorm:"column:to_id" json:"to_id"`
	IsRead   int16  `gorm:"column:is_read" json:"is_read"`   // 0 unread 1 read
	MsgType  int16  `gorm:"column:msg_type" json:"msg_type"` // default 1
	Status   int16  `gorm:"column:status" json:"status"`
	Data     string `gorm:"column:data" json:"data"`
	SendTime string `json:"send_time,omitempty"`
}

func (m *MessageM) TableName() string {
	return "messages"
}

// offline_messages
type OfflineMessageM struct {
	ModelID
	ModelTimestamps
	Message   json.RawMessage `gorm:"column:message" json:"message"` // message body (JSON)
	SendTime  int64           `gorm:"column:send_time" json:"send_time"`
	Status    int16           `gorm:"column:status" json:"status"`         // 0 not pushed 1 pushed
	ReceiveID int64           `gorm:"column:receive_id" json:"receive_id"` // receiver id
}

func (m *OfflineMessageM) TableName() string {
	return "offline_messages"
}
