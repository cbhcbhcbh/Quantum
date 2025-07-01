package friends

import (
	"net/http"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/known"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
	v1 "github.com/cbhcbhcbh/Quantum/pkg/api/v1"
	"github.com/cbhcbhcbh/Quantum/pkg/response"
	"github.com/gin-gonic/gin"
)

func (fc *FriendController) AddressList(c *gin.Context) {
	log.C(c).Infow("AddressList function called")

	var friends *[]v1.FriendDetail
	var err error

	id := c.GetInt64(known.XIdKey)

	if friends, err = fc.b.Friends().GetAllFriends(c, id); err != nil {
		response.FailResponse(http.StatusInternalServerError, err.Error()).ToJson(c)
		return
	}

	response.SuccessResponse(friends).ToJson(c)
}
