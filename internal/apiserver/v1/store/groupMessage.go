package store

import (
	"context"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	"gorm.io/gorm"
)

type GroupMessageStore interface {
	Create(ctx context.Context, message *model.GroupMessageM) error
}

type groupMessage struct {
	db *gorm.DB
}

var _ MessageStore = (*message)(nil)

func NewGroupMessage(db *gorm.DB) GroupMessageStore {
	return &groupMessage{
		db: db,
	}
}

func (m *groupMessage) Create(ctx context.Context, message *model.GroupMessageM) error {
	return m.db.WithContext(ctx).Create(message).Error
}
