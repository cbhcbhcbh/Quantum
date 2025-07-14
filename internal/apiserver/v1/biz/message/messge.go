package message

import (
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/store"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/enum"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/message"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	messageHandle "github.com/cbhcbhcbh/Quantum/internal/service/message"
	v1 "github.com/cbhcbhcbh/Quantum/pkg/api/v1"
	"github.com/gin-gonic/gin"
)

type MessageBiz interface {
	PrivateMessage(ctx *gin.Context, formId, toId int64) (*[]v1.MessageDetail, int64, error)
	SendMessage(c *gin.Context, form v1.PrivateMessageRequest) (int, error)
}

type messageBiz struct {
	ds store.IStore
}

var _ MessageBiz = (*messageBiz)(nil)

func New(ds store.IStore) MessageBiz {
	return &messageBiz{ds: ds}
}

func (m *messageBiz) PrivateMessage(ctx *gin.Context, formId, toId int64) (*[]v1.MessageDetail, int64, error) {
	list, total, err := m.ds.Message().ListByFormIdAndToId(ctx, formId, toId)
	if err != nil {
		return nil, total, err
	}

	user, err := m.ds.Users().GetById(ctx, toId)
	if err != nil {
		return nil, total, err
	}

	var messageDetails []v1.MessageDetail
	for _, message := range list {
		messageDetails = append(messageDetails, v1.MessageDetail{
			ID:       message.ID,
			FormID:   message.FormID,
			ToID:     message.ToID,
			MsgType:  message.MsgType,
			Msg:      message.Msg,
			SendTime: message.CreatedAt,
			Data:     message.Data,
			User: v1.UserMessage{
				ID:     user.ID,
				Name:   user.Name,
				Avatar: user.Avatar,
				Email:  user.Email,
			},
		})
	}

	return &messageDetails, total, nil
}

func (mc *messageBiz) SendMessage(c *gin.Context, form v1.PrivateMessageRequest) (int, error) {

	var messageService messageHandle.MessageService

	messageModel := &model.MessageM{
		FormID:   form.FormID,
		ToID:     form.ToID,
		MsgType:  form.MsgType,
		Msg:      form.Message,
		SendTime: form.SendTime,
		Data:     form.Data,
		IsRead:   enum.IsNotRead,
		Status:   1,
	}

	switch form.ChannelType {
	case int16(message.ChannelTypePrivate):
		_ = mc.ds.Message().Create(c, messageModel)
		user, err := mc.ds.Users().GetById(c, form.ToID)
		if err != nil {
			return 0, err
		}
		if user.UserType == model.BOT_TYPE {
			// TODO：bot chat logic need to develop
			messageService.SendPrivateMessage(form)
			return 100, nil
		} else {
			if _, err := mc.ds.Friends().GetByFormIDAndToID(c, form.FormID, form.ToID); err != nil {
				return 0, err
			}

			messageService.SendPrivateMessage(form)
			return 100, nil
		}
	case int16(message.ChannelTypeGroup):
		if _, err := mc.ds.Group().GetByUserIDAndID(c, form.FormID, form.ToID); err != nil {
			return 0, err
		}
		// TODO: implement send group message method
		messageService.SendGroupMessage(form)
	}

	return 100, nil
}
