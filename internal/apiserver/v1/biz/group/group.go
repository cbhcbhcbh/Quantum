package group

import (
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/store"
	v1 "github.com/cbhcbhcbh/Quantum/pkg/api/v1"
	"github.com/gin-gonic/gin"
)

type GroupBiz interface {
	GetAllGroups(ctx *gin.Context, formId int64) (*[]v1.GroupDetail, error)
}

type groupBiz struct {
	ds store.IStore
}

var _ GroupBiz = (*groupBiz)(nil)

func New(ds store.IStore) GroupBiz {
	return &groupBiz{ds: ds}
}

func (g *groupBiz) GetAllGroups(ctx *gin.Context, userId int64) (*[]v1.GroupDetail, error) {
	c := ctx.Request.Context()

	groups, err := g.ds.Group().GetByUserID(c, userId)
	if err != nil {
		return nil, err
	}

	var groupDetails []v1.GroupDetail
	for _, group := range *groups {
		groupDetails = append(groupDetails, v1.GroupDetail{
			ID:        group.ID,
			GroupType: group.GroupType,
			Name:      group.Name,
			Info:      group.Info,
			Avatar:    group.Avatar,
		})
	}

	return &groupDetails, nil
}
