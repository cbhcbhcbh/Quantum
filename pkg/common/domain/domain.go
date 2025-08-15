package domain

import (
	"encoding/json"
	"strconv"

	"github.com/cbhcbhcbh/Quantum/pkg/common/jwt"
	"github.com/cbhcbhcbh/Quantum/pkg/common/known"
)

type SuccessMessage struct {
	Message string `json:"msg" example:"ok"`
}

type ErrResponse struct {
	Message string `json:"msg"`
}

type UserIDsPresenter struct {
	UserIDs []string `json:"user_ids"`
}

type User struct {
	ID   uint64
	Name string
}

type Channel struct {
	ID          uint64
	AccessToken string
}

const (
	EventText = iota
	EventAction
	EventSeen
	EventFile
)

type Action string

var (
	WaitingMessage   Action = "waiting"
	JoinedMessage    Action = "joined"
	IsTypingMessage  Action = "istyping"
	EndTypingMessage Action = "endtyping"
	OfflineMessage   Action = "offline"
	LeavedMessage    Action = "leaved"
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

type MessagesPresenter struct {
	NextPageState string             `json:"next_ps"`
	Messages      []MessagePresenter `json:"messages"`
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

func (m *MessagePresenter) ToMessage(accessToken string) (*Message, error) {
	authResult, err := jwt.Auth(&jwt.AuthPayload{
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}
	if authResult.Expired {
		return nil, known.ErrTokenExpired
	}
	channelID := authResult.ChannelID
	userID, err := strconv.ParseUint(m.UserID, 10, 64)
	if err != nil {
		return nil, err
	}
	return &Message{
		Event:     m.Event,
		ChannelID: channelID,
		UserID:    userID,
		Payload:   m.Payload,
		Time:      m.Time,
	}, nil
}

func (m *Message) Encode() []byte {
	result, _ := json.Marshal(m)
	return result
}

func (m *MessagePresenter) Encode() []byte {
	result, _ := json.Marshal(m)
	return result
}
