package store

import (
	"context"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	"gorm.io/gorm"
)

type FriendsStore interface {
	Create(ctx context.Context, message *model.FriendM) error
	GetByFormID(ctx context.Context, formId int64) (*[]model.FriendM, error)
	GetByFormIDAndToID(ctx context.Context, formId, toId int64) (*model.FriendM, error)
	Delete(ctx context.Context, formId, toId int64) error
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

func (f *friends) GetByFormID(ctx context.Context, formId int64) (*[]model.FriendM, error) {
	var friends []model.FriendM
	if err := f.db.WithContext(ctx).Preload("Users").Where("form_id = ?", formId).Find(&friends).Error; err != nil {
		return nil, err
	}

	return &friends, nil
}

func (f *friends) GetByFormIDAndToID(ctx context.Context, formId, toId int64) (*model.FriendM, error) {
	var friend model.FriendM
	if err := f.db.WithContext(ctx).Preload("Users").Where("form_id = ? and to_id = ?", formId, toId).Find(&friend).Error; err != nil {
		return nil, err
	}

	return &friend, nil
}

func (f *friends) Delete(ctx context.Context, formId, toId int64) error {
	return f.db.WithContext(ctx).Where("form_id = ? and to_id = ?", formId, toId).Delete(&model.FriendM{}).Error
}
