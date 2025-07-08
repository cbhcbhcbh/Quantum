package store

import (
	"context"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	"gorm.io/gorm"
)

type FriendRecordStore interface {
	Create(ctx context.Context, message *model.FriendRecordM) error
	GetByID(ctx context.Context, id int64, status int16) (*model.FriendRecordM, error)
	GetByFormID(ctx context.Context, formId int64) (*[]model.FriendRecordM, error)
	GetByFormIDAndToID(ctx context.Context, formId, toId int64, status int16) (*model.FriendRecordM, error)
	ListByFormIDAndToID(ctx context.Context, formId, toId int64, status int16) (*[]model.FriendRecordM, error)
	Delete(ctx context.Context, formId, toId int64) error
	Update(ctx context.Context, record *model.FriendRecordM) error
}

type friendRecord struct {
	db *gorm.DB
}

var _ FriendRecordStore = (*friendRecord)(nil)

func NewFriendRecord(db *gorm.DB) FriendRecordStore {
	return &friendRecord{
		db: db,
	}
}

func (f *friendRecord) Create(ctx context.Context, message *model.FriendRecordM) error {
	return f.db.WithContext(ctx).Create(message).Error
}

func (f *friendRecord) GetByFormID(ctx context.Context, formId int64) (*[]model.FriendRecordM, error) {
	var friends []model.FriendRecordM
	if err := f.db.WithContext(ctx).Where("form_id = ?", formId).Find(&friends).Error; err != nil {
		return nil, err
	}

	return &friends, nil
}

func (f *friendRecord) GetByFormIDAndToID(ctx context.Context, formId, toId int64, status int16) (*model.FriendRecordM, error) {
	var friend model.FriendRecordM
	err := f.db.WithContext(ctx).Where("form_id = ? and to_id = ? and status = ?", formId, toId, status).Preload("Users").First(&friend).Error
	if err == nil {
		return &friend, nil
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

func (f *friendRecord) Delete(ctx context.Context, formId, toId int64) error {
	return f.db.WithContext(ctx).Where("form_id = ? and to_id = ?", formId, toId).Delete(&model.FriendRecordM{}).Error
}

func (f *friendRecord) ListByFormIDAndToID(ctx context.Context, formId, toId int64, status int16) (*[]model.FriendRecordM, error) {
	var friend []model.FriendRecordM
	db := f.db.WithContext(ctx).Where("form_id = ? and to_id = ?", formId, toId).Preload("Users")
	if status != -1 {
		db = db.Where("status = ?", status)
	}
	err := db.Order("created_at desc").Find(&friend).Error
	if err == nil {
		return &friend, nil
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

func (f *friendRecord) GetByID(ctx context.Context, id int64, status int16) (*model.FriendRecordM, error) {
	var friend model.FriendRecordM
	db := f.db.WithContext(ctx).Where("id = ? ", id).Preload("Users")
	if status != -1 {
		db = db.Where("status = ?", status)
	}
	err := db.Order("created_at desc").First(&friend).Error
	if err == nil {
		return &friend, nil
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

func (f *friendRecord) Update(ctx context.Context, record *model.FriendRecordM) error {
	return f.db.WithContext(ctx).Model(&model.FriendRecordM{}).Where("id = ?", record.ID).Updates(record).Error
}