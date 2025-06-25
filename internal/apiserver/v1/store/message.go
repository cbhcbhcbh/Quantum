package store

import (
	"context"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	"gorm.io/gorm"
)

type MessageStore interface {
	Create(ctx context.Context, message *model.MessageM) error
}

type message struct {
	db *gorm.DB
}

var _ MessageStore = (*message)(nil)

func NewMessage(db *gorm.DB) MessageStore {
	return &message{
		db: db,
	}
}

func (m *message) Create(ctx context.Context, message *model.MessageM) error {
	return m.db.WithContext(ctx).Create(message).Error
}
