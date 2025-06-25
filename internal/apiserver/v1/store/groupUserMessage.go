package store

import (
	"context"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	"gorm.io/gorm"
)

type GroupUserMessageStore interface {
	Create(ctx context.Context, message *model.GroupUserMessageM) error
}

type groupUserMessage struct {
	db *gorm.DB
}

var _ GroupUserMessageStore = (*groupUserMessage)(nil)

func NewGroupUserMessage(db *gorm.DB) GroupUserMessageStore {
	return &groupUserMessage{
		db: db,
	}
}

func (m *groupUserMessage) Create(ctx context.Context, message *model.GroupUserMessageM) error {
	return m.db.WithContext(ctx).Create(message).Error
}
