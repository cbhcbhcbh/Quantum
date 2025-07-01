package store

import (
	"context"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	"gorm.io/gorm"
)

type FriendsStore interface {
	Create(ctx context.Context, message *model.FriendM) error
	GetByFormID(ctx context.Context, fromId int64) (*[]model.FriendM, error)
}

type friends struct {
	db *gorm.DB
}

var _ FriendsStore = (*friends)(nil)

func NewFriends(db *gorm.DB) FriendsStore {
	return &friends{
		db: db,
	}
}

func (f *friends) Create(ctx context.Context, message *model.FriendM) error {
	return f.db.WithContext(ctx).Create(message).Error
}

func (f *friends) GetByFormID(ctx context.Context, fromId int64) (*[]model.FriendM, error) {
	var friends []model.FriendM
	if err := f.db.WithContext(ctx).Where("form_id = ?", fromId).Find(&friends).Error; err != nil {
		return nil, err
	}

	return &friends, nil
}
