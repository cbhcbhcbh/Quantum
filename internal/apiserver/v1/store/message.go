package store

import (
	"context"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	"gorm.io/gorm"
)

type MessageStore interface {
	ListByFormIdAndToId(ctx context.Context, formId, toId int64) ([]*model.MessageM, error)
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

func (m *message) ListByFormIdAndToId(ctx context.Context, formId, toId int64) ([]*model.MessageM, error) {
	m.db.WithContext(ctx).Where
}

