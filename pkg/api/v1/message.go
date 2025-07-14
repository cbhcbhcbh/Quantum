package v1

import "time"

type MessageDetail struct {
	ID       int64       `json:"id" validate:"required"`
	SendTime time.Time   `json:"created_at" validate:"omitempty"`
	Msg      string      `json:"msg" validate:"omitempty"`
	FormID   int64       `json:"form_id" validate:"required"`
	ToID     int64       `json:"to_id" validate:"required"`
	IsRead   int16       `json:"is_read" validate:"required"`
	MsgType  int16       `json:"msg_type" validate:"required"`
	Status   string      `json:"theme" validate:"required"`
	Data     string      `json:"data" validate:"required"`
	User     UserMessage `json:"user" validate:"required"` // User who sent the message
}

type UserMessage struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Email  string `json:"email"`
	ID     int64  `json:"id"`
}

type PrivateMessageRequest struct {
	MsgId       int64  `json:"msg_id"`                                        // Server-side unique message ID
	MsgClientId int64  `json:"msg_client_id" validate:"required"`             // Client-side unique message ID
	MsgCode     int16  `json:"msg_code" validate:"required"`                  // Defined message code
	FormID      int64  `json:"form_id" validate:"required"`                   // ID of the message sender
	ToID        int64  `json:"to_id" validate:"required"`                     // ID of the message receiver
	MsgType     int16  `json:"msg_type" validate:"required"`                  // Message type: 1.text 2.voice 3.file 5.leave group 6.block
	ChannelType int16  `json:"channel_type" validate:"required,gte=1,lte=3" ` // Channel type: 1.private chat 2.channel 3.broadcast
	Message     string `json:"message" validate:"required"`                   // Message content
	SendTime    string `json:"send_time,omitempty"`                           // Message send time
	Data        string `json:"data"`                                          // Custom data payload
	UserId      int64  `json:"user_id"`                                       // Custom user data
}
