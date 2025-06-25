package store

import (
	"context"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	"gorm.io/gorm"
)

type OfflineMessageStore interface {
	Create(ctx context.Context, message *model.OfflineMessageM) error
}

type offlineMessage struct {
	db *gorm.DB
}

var _ OfflineMessageStore = (*offlineMessage)(nil)

func NewOfflineMessage(db *gorm.DB) OfflineMessageStore {
	return &offlineMessage{
		db: db,
	}
}

func (m *offlineMessage) Create(ctx context.Context, message *model.OfflineMessageM) error {
	return m.db.WithContext(ctx).Create(message).Error
}
