package friendrecord

import (
	"errors"

	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/store"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	"github.com/gin-gonic/gin"
)

type FriendRecordBiz interface {
	SendFriendRequest(ctx *gin.Context, formId, toId int64, information string) (*model.FriendRecordM, error)
	ListFriendRequests(ctx *gin.Context, formId, toId int64) (*[]model.FriendRecordM, error)
	GetFriendRequest(ctx *gin.Context, id int64) (*model.FriendRecordM, error)
	UpdateFriendRequest(ctx *gin.Context, record *model.FriendRecordM) error
}

type friendRecordBiz struct {
	ds store.IStore
}

var _ FriendRecordBiz = (*friendRecordBiz)(nil)

func New(ds store.IStore) FriendRecordBiz {
	return &friendRecordBiz{ds: ds}
}

func (b *friendRecordBiz) SendFriendRequest(ctx *gin.Context, formId, toId int64, information string) (*model.FriendRecordM, error) {
	if _, err := b.ds.Users().GetById(ctx, toId); err != nil {
		return nil, err
	}

	friendRecord, err := b.ds.FriendRecord().GetByFormIDAndToID(ctx, formId, toId, model.FriendStatusPending)
	if err != nil {
		return nil, err
	}

	if friendRecord != nil {
		return nil, errors.New("friend request already exists")
	}

	friend, err := b.ds.Friends().GetByFormIDAndToID(ctx, formId, toId)
	if err != nil {
		return nil, err
	}

	if friend != nil {
		return nil, errors.New("existing friend relationship")
	}

	records := &model.FriendRecordM{
		FormID:      formId,
		ToID:        toId,
		Status:      model.FriendStatusPending,
		Information: information,
	}

	if err = b.ds.FriendRecord().Create(ctx, records); err != nil {
		return nil, err
	}

	return records, nil

}

func (b *friendRecordBiz) ListFriendRequests(ctx *gin.Context, formId, toId int64) (*[]model.FriendRecordM, error) {
	records, err := b.ds.FriendRecord().ListByFormIDAndToID(ctx, formId, toId, model.FriendStatusUnchecking)
	if err != nil {
		return nil, err
	}

	if records == nil {
		return nil, errors.New("no friend requests found")
	}

	return records, nil
}

func (b *friendRecordBiz) GetFriendRequest(ctx *gin.Context, id int64) (*model.FriendRecordM, error) {
	var record *model.FriendRecordM
	var err error

	if record, err = b.ds.FriendRecord().GetByID(ctx, id, model.FriendStatusPending); err != nil {
		return nil, err
	}
	return record, nil
}

func (b *friendRecordBiz) UpdateFriendRequest(ctx *gin.Context, record *model.FriendRecordM) error {
	if err := b.ds.FriendRecord().Update(ctx, record); err != nil {
		return err
	}
	return nil
}
