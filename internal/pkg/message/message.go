package message

import (
	"encoding/json"
	"fmt"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/date"
	"github.com/google/uuid"
)

// WsMessage represents the structure of a WebSocket message exchanged between client and server.
type WsMessage struct {
	MsgId       string `json:"msg_id"`        // Unique identifier for the message
	MsgClientId int64  `json:"msg_client_id"` // Client-generated message ID for deduplication or tracking
	MsgCode     int    `json:"msg_code"`      // Message code for business logic or status
	FormID      int64  `json:"form_id"`       // Sender's user ID
	ToID        int64  `json:"to_id"`         // Recipient's user ID
	MsgType     int    `json:"msg_type"`      // Type of the message (e.g., 1.text, 2.voice, 3.file)
	ChannelType int    `json:"channel_type"`  // Channel type (e.g., 1.private, 2.group, 3.broadcast)
	Message     string `json:"message"`       // Main message content
	SendTime    int64  `json:"send_time"`     // Timestamp when the message was sent (Unix time)
	Data        any    `json:"data"`          // Additional data or payload (can be any type)
}

type AckMessage struct {
	Message     string `json:"message"`       // Acknowledgment message content
	MsgId       string `json:"msg_id"`        // ID of the message being acknowledged
	MsgCode     int    `json:"msg_code"`      // Code indicating the status of the acknowledgment
	MsgClientId int64  `json:"msg_client_id"` // Client-generated message ID for tracking
}

func ValidationMsg(msg []byte) (string, string, error) {
	var wsMsg WsMessage
	var ackMsg AckMessage
	var err error

	if len(msg) == 0 {
		ackMsg = AckMessage{
			MsgId:       "",
			MsgClientId: 0,
			MsgCode:     500,
			Message:     "empty message",
		}
		ackMsgByte, _ := json.Marshal(ackMsg)
		return `{"code":500,"message":"请勿发送空消息"}`, string(ackMsgByte), fmt.Errorf("empty message")
	}

	if err = json.Unmarshal(msg, &wsMsg); err != nil {
		ackMsg = AckMessage{
			MsgId:       "",
			MsgClientId: 0,
			MsgCode:     500,
			Message:     "invalid message format",
		}
		ackMsgByte, _ := json.Marshal(ackMsg)
		return `{"code":500,"message":"消息格式错误"}`, string(ackMsgByte), fmt.Errorf("message unmarshal error: %v", err)
	}

	wsMsg.MsgId = uuid.New().String()
	wsMsg.SendTime = date.TimeUnix()

	ackMsg = AckMessage{
		MsgId:       wsMsg.MsgId,
		MsgClientId: wsMsg.MsgClientId,
		MsgCode:     200,
		Message:     "ack",
	}

	msgByte, err1 := json.Marshal(wsMsg)
	ackMsgByte, err2 := json.Marshal(ackMsg)
	if err1 != nil || err2 != nil {
		ackMsg.MsgCode = 500
		ackMsg.Message = "json marshal error"
		ackMsgByte, _ = json.Marshal(ackMsg)
		return "", string(ackMsgByte), fmt.Errorf("json marshal error")
	}

	return string(msgByte), string(ackMsgByte), nil
}
