package store

import (
	"context"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	"gorm.io/gorm"
)

type GroupOfflineMessageStore interface {
	Create(ctx context.Context, message *model.GroupOfflineMessageM) error
}

type groupOfflineMessage struct {
	db *gorm.DB
}

var _ GroupOfflineMessageStore = (*groupOfflineMessage)(nil)

func NewGroupOfflineMessage(db *gorm.DB) GroupOfflineMessageStore {
	return &groupOfflineMessage{
		db: db,
	}
}

func (m *groupOfflineMessage) Create(ctx context.Context, message *model.GroupOfflineMessageM) error {
	return m.db.WithContext(ctx).Create(message).Error
}
