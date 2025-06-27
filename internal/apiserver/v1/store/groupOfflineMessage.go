package store

import (
	"context"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	"gorm.io/gorm"
)

type GroupOfflineMessageStore interface {
	Create(ctx context.Context, message *model.GroupOfflineMessageM) error
	ListByTimeRangeAndStatus(ctx context.Context, startTime, endTime int64, status int16) ([]*model.GroupOfflineMessageM, error)
	UpdateStatuByID(ctx context.Context, ids []int64, status int16) error
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

func (m *groupOfflineMessage) ListByTimeRangeAndStatus(ctx context.Context, startTime, endTime int64, status int16) ([]*model.GroupOfflineMessageM, error) {
	var messages []*model.GroupOfflineMessageM
	err := m.db.WithContext(ctx).
		Where("send_time >= ? AND send_time <= ? AND status = ?", startTime, endTime, status).
		Find(&messages).Error
	return messages, err
}

func (m *groupOfflineMessage) UpdateStatuByID(ctx context.Context, ids []int64, status int16) error {
	return m.db.
		WithContext(ctx).
		Model(&model.GroupOfflineMessageM{}).
		Where("id IN ?", ids).
		Update("status", status).Error
}
