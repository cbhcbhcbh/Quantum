/*
*

	@author:panliang
	@data:2022/7/30
	@note

*
*/
package grpcMessage

import (
	"context"
	"fmt"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/date"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/enum"
	v1 "github.com/cbhcbhcbh/Quantum/pkg/api/v1"
)

type ImGrpcMessage struct {
}

func (ImGrpcMessage) mustEmbedUnimplementedImMessageServer() {}

func (ImGrpcMessage) SendMessageHandler(c context.Context, request *SendMessageRequest) (*SendMessageResponse, error) {
	params := v1.PrivateMessageRequest{
		MsgId:       date.TimeUnixNano(),
		MsgCode:     enum.WsChantMessage,
		MsgClientId: request.MsgClientId,
		FormID:      request.FormId,
		ToID:        request.ToId,
		ChannelType: int16(request.ChannelType),
		MsgType:     int16(request.MsgType),
		Message:     request.Message,
		SendTime:    date.NewDate(),
		Data:        request.Data,
	}

	msgString := GetGrpcPrivateChatMessages(params)

	switch request.ChannelType {
	case 1:
		fmt.Println(msgString)
		//client.ImManager.PrivateChannel <- []byte(msgString)
	case 2:

	}
	return &SendMessageResponse{Code: 200, Message: "Success"}, nil
}
func GetGrpcPrivateChatMessages(message v1.PrivateMessageRequest) string {
	msg := fmt.Sprintf(`{
                "msg_id": %d,
                "msg_client_id": %d,
                "msg_code": %d,
                "form_id": %d,
                "to_id": %d,
                "msg_type": %d,
                "channel_type": %d,
                "message": %s,
                "data": %s
        }`, message.MsgId, message.MsgClientId, message.MsgCode, message.FormID, message.ToID, message.MsgType, message.ChannelType, message.Message, message.Data)

	return msg
}
