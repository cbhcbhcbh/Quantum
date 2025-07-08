package friends

import (
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/store"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/date"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	v1 "github.com/cbhcbhcbh/Quantum/pkg/api/v1"
	"github.com/gin-gonic/gin"
)

type FriendBiz interface {
	CreateFriendRelation(ctx *gin.Context, formId, toId int64)
	GetAllFriends(ctx *gin.Context, formId int64) (*[]v1.FriendDetail, error)
	GetFriend(ctx *gin.Context, formId, toId int64) (*v1.FriendDetail, error)
	DeleteFriend(ctx *gin.Context, formId, toId int64) error
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

func (b *friendBiz) GetFriend(ctx *gin.Context, formId, toId int64) (*v1.FriendDetail, error) {
	c := ctx.Request.Context()

	friend, err := b.ds.Friends().GetByFormIDAndToID(c, formId, toId)
	if err != nil {
		return nil, err
	}
	if friend == nil {
		return nil, nil
	}

	friendDetail := &v1.FriendDetail{
		ID:     friend.FormID,
		Note:   friend.Note,
		Status: friend.Status,
		Uid:    friend.Uid,
		Users: v1.UserDetails{
			ID:            friend.Users.ID,
			Name:          friend.Users.Name,
			Email:         friend.Users.Email,
			Avatar:        friend.Users.Avatar,
			Status:        friend.Users.Status,
			Bio:           friend.Users.Bio,
			Sex:           friend.Users.Sex,
			Age:           friend.Users.Age,
			LastLoginTime: friend.Users.LastLoginTime,
			Uid:           friend.Users.Uid,
		},
	}

	return friendDetail, nil
}

func (b *friendBiz) DeleteFriend(ctx *gin.Context, formId, toId int64) error {
	c := ctx.Request.Context()

	if err := b.ds.Friends().Delete(c, formId, toId); err != nil {
		return err
	}

	return nil
}

func (b *friendBiz) CreateFriendRelation(ctx *gin.Context, formId, toId int64) {
	friendRelation1 := model.FriendM{
		FormID:  formId,
		ToID:    toId,
		Status:  model.FriendNotPinned,
		TopTime: date.NewDate(),
		Note:    "",
	}

	friendRelation2 := model.FriendM{
		FormID:  toId,
		ToID:    formId,
		Status:  model.FriendNotPinned,
		TopTime: date.NewDate(),
		Note:    "",
	}

	_ = b.ds.Friends().Create(ctx, &friendRelation1)
	_ = b.ds.Friends().Create(ctx, &friendRelation2)
}
