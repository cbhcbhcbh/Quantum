package friends

import (
	"net/http"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/enum"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/known"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
	v1 "github.com/cbhcbhcbh/Quantum/pkg/api/v1"
	"github.com/cbhcbhcbh/Quantum/pkg/response"
	"github.com/gin-gonic/gin"
)

func (fc *FriendController) Index(c *gin.Context) {
	log.C(c).Infow("Index function called")

	var friends *[]v1.FriendDetail
	var err error

	id := c.GetInt64(known.XIdKey)

	if friends, err = fc.b.Friends().GetAllFriends(c, id); err != nil {
		response.FailResponse(http.StatusInternalServerError, err.Error()).ToJson(c)
		return
	}

	response.SuccessResponse(friends).ToJson(c)
}

func (fc *FriendController) Show(c *gin.Context) {
	log.C(c).Infow("Friend Show function called")

	id := c.GetInt64(known.XIdKey)

	var r v1.FriendInfo
	var friend *v1.FriendDetail
	var err error
	if err = c.ShouldBindUri(&r); err != nil {
		response.FailResponse(enum.ParamError, err.Error()).ToJson(c)
		return
	}

	if friend, err = fc.b.Friends().GetFriend(c, id, r.FriendId); err != nil {
		response.FailResponse(http.StatusInternalServerError, err.Error()).ToJson(c)
		return
	}

	response.SuccessResponse(friend).ToJson(c)
}

func (fc *FriendController) GetUserStatus(c *gin.Context) {

}

func (fc *FriendController) Delete(c *gin.Context) {
	log.C(c).Infow("Friend Delete function called")

	id := c.GetInt64(known.XIdKey)

	var r v1.FriendInfo
	var err error
	if err = c.ShouldBindUri(&r); err != nil {
		response.FailResponse(enum.ParamError, err.Error()).ToJson(c)
		return
	}

	if err = fc.b.Friends().DeleteFriend(c, id, r.FriendId); err != nil {
		response.FailResponse(http.StatusInternalServerError, err.Error()).ToJson(c)
		return
	}

	response.SuccessResponse().ToJson(c)
}
