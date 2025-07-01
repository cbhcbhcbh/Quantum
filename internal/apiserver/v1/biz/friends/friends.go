package friends

import (
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/store"
	v1 "github.com/cbhcbhcbh/Quantum/pkg/api/v1"
	"github.com/gin-gonic/gin"
)

type FriendBiz interface {
	GetAllFriends(ctx *gin.Context, formId int64) (*[]v1.FriendDetail, error)
}

type friendBiz struct {
	ds store.IStore
}

var _ FriendBiz = (*friendBiz)(nil)

func New(ds store.IStore) FriendBiz {
	return &friendBiz{ds: ds}
}

func (b *friendBiz) GetAllFriends(ctx *gin.Context, formId int64) (*[]v1.FriendDetail, error) {
	c := ctx.Request.Context()

	friends, err := b.ds.Friends().GetByFormID(c, formId)
	if err != nil {
		return nil, err
	}

	var friendDetails []v1.FriendDetail
	for _, friend := range *friends {
		friendDetails = append(friendDetails, v1.FriendDetail{
			ID:     friend.FormID,
			Note:   friend.Note,
			Status: friend.Status,
			Uid:    friend.Uid,
		})
	}

	return &friendDetails, nil
}
