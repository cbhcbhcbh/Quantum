package store

import (
	"context"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	"gorm.io/gorm"
)

type ImUsersStore interface {
	Create(ctx context.Context, user *model.ImUsersM) error
}

type imUsers struct {
	db *gorm.DB
}

var _ ImUsersStore = (*imUsers)(nil)

func NewImUsers(db *gorm.DB) ImUsersStore {
	return &imUsers{
		db: db,
	}
}

func (iu *imUsers) Create(ctx context.Context, user *model.ImUsersM) error {
	if err := iu.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}
