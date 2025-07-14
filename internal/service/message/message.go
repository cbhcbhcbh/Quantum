package message

import (
	"github.com/cbhcbhcbh/Quantum/internal/pkg/message"
	"github.com/cbhcbhcbh/Quantum/internal/service/client"
	v1 "github.com/cbhcbhcbh/Quantum/pkg/api/v1"
)

type MessageService struct {
}

type IMessageService interface {
	SendFriendActionMessage(msg message.CreateFriendMessage)
}

func (*MessageService) SendFriendActionMessage(msg message.CreateFriendMessage) {
	client.Manager.SendFriendActionMessage(msg)
}

func (*MessageService) SendPrivateMessage(msg v1.PrivateMessageRequest) (bool, string) {
	isOk, respMessage := client.Manager.SendPrivateMessage(msg)
	return isOk, respMessage
}

func (*MessageService) SendGroupMessage(msg v1.PrivateMessageRequest) (bool, string) {
	isOk, respMessage := client.Manager.SendGroupMessage(msg)
	return isOk, respMessage
}
