package message

import (
	"github.com/cbhcbhcbh/Quantum/internal/pkg/message"
	"github.com/cbhcbhcbh/Quantum/internal/service/client"
)

type MessageService struct {
}

type IMessageService interface {
	SendFriendActionMessage(msg message.CreateFriendMessage)
}

func (*MessageService) SendFriendActionMessage(msg message.CreateFriendMessage) {
	client.Manager.SendFriendActionMessage(msg)
}
