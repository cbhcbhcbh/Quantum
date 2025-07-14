package store

import (
	"context"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	"gorm.io/gorm"
)

type GroupStore interface {
	Create(ctx context.Context, message *model.GroupM) error
	GetByUserID(ctx context.Context, userId int64) (*[]model.GroupM, error)
	GetByUserIDAndID(ctx context.Context, userId int64, id int64) (*[]model.GroupM, error)
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

func (g *group) Create(ctx context.Context, message *model.GroupM) error {
	return g.db.WithContext(ctx).Create(message).Error
}

func (g *group) GetByUserID(ctx context.Context, userId int64) (*[]model.GroupM, error) {
	var groups []model.GroupM
	if err := g.db.WithContext(ctx).Where("user_id = ?", userId).Find(&groups).Error; err != nil {
		return nil, err
	}

	return &groups, nil
}

func (g *group) GetByUserIDAndID(ctx context.Context, userId int64, id int64) (*[]model.GroupM, error) {
	var groups []model.GroupM
	if err := g.db.WithContext(ctx).Where("user_id = ? and id = ?", userId, id).Find(&groups).Error; err != nil {
		return nil, err
	}

	return &groups, nil
}
