package store

import (
	"context"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	"gorm.io/gorm"
)

type OfflineMessageStore interface {
	Create(ctx context.Context, message *model.OfflineMessageM) error
	ListByTimeRangeAndStatus(ctx context.Context, startTime, endTime int64, status int16) ([]*model.OfflineMessageM, error)
	UpdateStatuByID(ctx context.Context, ids []int64, status int16) error
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

func (m *offlineMessage) ListByTimeRangeAndStatus(ctx context.Context, startTime, endTime int64, status int16) ([]*model.OfflineMessageM, error) {
	var messages []*model.OfflineMessageM
	err := m.db.WithContext(ctx).
		Where("send_time >= ? AND send_time <= ? AND status = ?", startTime, endTime, status).
		Find(&messages).Error
	return messages, err
}

func (m *offlineMessage) UpdateStatuByID(ctx context.Context, ids []int64, status int16) error {
	return m.db.
		WithContext(ctx).
		Model(&model.OfflineMessageM{}).
		Where("id IN ?", ids).
		Update("status", status).Error
}
