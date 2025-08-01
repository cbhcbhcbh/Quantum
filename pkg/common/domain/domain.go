package domain

import (
	"encoding/json"
	"strconv"
)

type Message struct {
	MessageID uint64 `json:"message_id"`
	Event     int    `json:"event"`
	ChannelID uint64 `json:"channel_id"`
	UserID    uint64 `json:"user_id"`
	Payload   string `json:"payload"`
	Seen      bool   `json:"seen"`
	Time      int64  `json:"time"`
}

type MessagePresenter struct {
	MessageID string `json:"message_id"`
	Event     int    `json:"event"`
	UserID    string `json:"user_id"`
	Payload   string `json:"payload"`
	Seen      bool   `json:"seen"`
	Time      int64  `json:"time"`
}

func (m *Message) ToPresenter() *MessagePresenter {
	return &MessagePresenter{
		MessageID: strconv.FormatUint(m.MessageID, 10),
		Event:     m.Event,
		UserID:    strconv.FormatUint(m.UserID, 10),
		Payload:   m.Payload,
		Seen:      m.Seen,
		Time:      m.Time,
	}
}

func (m *MessagePresenter) Encode() []byte {
	result, _ := json.Marshal(m)
	return result
}
