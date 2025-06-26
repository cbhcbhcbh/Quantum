package store

import (
	"context"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	"gorm.io/gorm"
)

type GroupUserStore interface {
	Create(ctx context.Context, message *model.GroupUserM) error
	List(ctx context.Context, groupID int64) ([]*model.GroupUserM, error)
}

type groupUser struct {
	db *gorm.DB
}

var _ GroupUserStore = (*groupUser)(nil)

func NewGroupUser(db *gorm.DB) GroupUserStore {
	return &groupUser{
		db: db,
	}
}

func (m *groupUser) Create(ctx context.Context, message *model.GroupUserM) error {
	return m.db.WithContext(ctx).Create(message).Error
}

func (m *groupUser) List(ctx context.Context, groupID int64) ([]*model.GroupUserM, error) {
	var groupUsers []*model.GroupUserM
	err := m.db.WithContext(ctx).Where("group_id = ?", groupID).Find(&groupUsers).Error
	if err != nil {
		return nil, err
	}
	return groupUsers, nil
}
