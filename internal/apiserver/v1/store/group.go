package store

import (
	"context"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	"gorm.io/gorm"
)

type GroupStore interface {
	Create(ctx context.Context, message *model.GroupM) error
}

type group struct {
	db *gorm.DB
}

var _ GroupStore = (*group)(nil)

func NewGroup(db *gorm.DB) GroupStore {
	return &group{
		db: db,
	}
}

func (m *group) Create(ctx context.Context, message *model.GroupM) error {
	return m.db.WithContext(ctx).Create(message).Error
}
