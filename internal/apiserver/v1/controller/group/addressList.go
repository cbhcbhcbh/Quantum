package group

import (
	"net/http"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/known"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
	v1 "github.com/cbhcbhcbh/Quantum/pkg/api/v1"
	"github.com/cbhcbhcbh/Quantum/pkg/response"
	"github.com/gin-gonic/gin"
)

func (gc *GroupController) AddressList(c *gin.Context) {
	log.C(c).Infow("AddressList function called")

	var groups *[]v1.GroupDetail
	var err error

	id := c.GetInt64(known.XIdKey)

	if groups, err = gc.b.Groups().GetAllGroups(c, id); err != nil {
		response.FailResponse(http.StatusInternalServerError, err.Error()).ToJson(c)
		return
	}

	response.SuccessResponse(groups).ToJson(c)
}
