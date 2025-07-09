package friends

import (
	"net/http"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/date"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/enum"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/known"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
	MsgStruct "github.com/cbhcbhcbh/Quantum/internal/pkg/message"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	"github.com/cbhcbhcbh/Quantum/internal/service/message"
	v1 "github.com/cbhcbhcbh/Quantum/pkg/api/v1"
	"github.com/cbhcbhcbh/Quantum/pkg/response"
	"github.com/gin-gonic/gin"
)

// TODO: Refactoor FriendController All method
func (fc *FriendController) SendFriendRequest(c *gin.Context) {
	log.C(c).Infow("Friend SendFriendRequest function called")

	var form v1.CreateFriendRequest
	if err := c.ShouldBindJSON(&form); err != nil {
		response.FailResponse(http.StatusBadRequest, "Invalid input").ToJson(c)
		return
	}

	id := c.GetInt64(known.XIdKey)
	r := &v1.Person{
		ID: id,
	}
	userInfo, _ := fc.b.Users().GetUserInfo(c, r)
	records, err := fc.b.FriendRecord().SendFriendRequest(c, id, form.ToID, form.Information)
	if err != nil {
		response.FailResponse(http.StatusInternalServerError, err.Error()).ToJson(c)
		return
	}

	// push msg to user
	var messageService message.MessageService

	var msg MsgStruct.CreateFriendMessage

	msg.MsgCode = enum.WsCreate
	msg.ID = records.ID
	msg.ToID = records.ToID
	msg.FormID = records.FormID
	msg.Information = records.Information
	msg.Status = records.Status
	msg.Users.ID = userInfo.ID
	msg.Users.Avatar = userInfo.Avatar
	msg.Users.Name = userInfo.Name

	messageService.SendFriendActionMessage(msg)

	records.Users.Name = userInfo.Name
	records.Users.Id = userInfo.ID
	records.Users.Avatar = userInfo.Avatar
	response.SuccessResponse(records).ToJson(c)

}

func (fc *FriendController) ListFriendRequests(c *gin.Context) {
	log.C(c).Infow("Friend ListFriendRequests function called")

	id := c.GetInt64(known.XIdKey)
	var err error
	var records *[]model.FriendRecordM
	if records, err = fc.b.FriendRecord().ListFriendRequests(c, id, id); err != nil {
		return
	}

	response.SuccessResponse(records).ToJson(c)

}

func (fc *FriendController) AcceptFriendRequest(c *gin.Context) {
	log.C(c).Infow("Friend AcceptFriendRequest function called")

	var err error
	var form v1.UpdateFriendRequest
	if err := c.ShouldBindJSON(&form); err != nil {
		response.FailResponse(http.StatusBadRequest, "Invalid input").ToJson(c)
		return
	}

	id := c.GetInt64(known.XIdKey)

	var record *model.FriendRecordM

	if record, err = fc.b.FriendRecord().GetFriendRequest(c, form.ID); err != nil {
		return
	}

	var friend *v1.FriendDetail
	if friend, err = fc.b.Friends().GetFriend(c, id, record.ToID); err != nil {
		return
	}

	var user *v1.UserDetails
	if user, err = fc.b.Users().GetUserInfo(c, &v1.Person{
		ID: id,
	}); err != nil {
		return
	}

	record.Status = form.Status
	if err := fc.b.FriendRecord().UpdateFriendRequest(c, record); err != nil {
		return
	}

	// push msg to user
	var messageService message.MessageService

	var msg MsgStruct.CreateFriendMessage
	var msgCode int

	if form.Status == 1 {
		msgCode = enum.WsFriendOk
		fc.b.Friends().CreateFriendRelation(c, record.FormID, record.ToID)
		fc.b.Sessions().CreateSessionRelation(c, record.FormID, record.ToID, 1, user)
	} else {
		msgCode = enum.WsFriendError
	}

	msg.MsgCode = msgCode
	msg.ID = record.ID
	msg.ToID = record.FormID
	msg.FormID = record.ToID
	msg.Information = record.Information
	msg.CreatedAt = date.TimeToString(record.CreatedAt)
	msg.Status = record.Status
	msg.Users.ID = user.ID
	msg.Users.Avatar = user.Avatar
	msg.Users.Name = user.Name

	messageService.SendFriendActionMessage(msg)
	friend.Status = form.Status
	response.SuccessResponse(friend).WriteTo(c)

}
