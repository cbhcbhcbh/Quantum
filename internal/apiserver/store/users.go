package store

import (
	"context"
	"errors"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/code"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	"gorm.io/gorm"
)

type UsersStore interface {
	Create(ctx context.Context, user *model.UsersM) error
	Get(ctx context.Context, name string) (*model.UsersM, error)
	Update(ctx context.Context, user *model.UsersM) error
	List(ctx context.Context, offset, limit int) (int64, []*model.UsersM, error)
	Delete(ctx context.Context, name string) error
}

type users struct {
	db *gorm.DB
}

var _ UsersStore = (*users)(nil)

func NewUsers(db *gorm.DB) UsersStore {
	return &users{
		db: db,
	}
}

func (u *users) Create(ctx context.Context, user *model.UsersM) error {
	var count int64
	if err := u.db.WithContext(ctx).Model(&model.UsersM{}).Where("name = ?", user.Name).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New(code.ErrUserAlreadyExist.Message)
	}

	return u.db.WithContext(ctx).Create(user).Error
}

func (u *users) Get(ctx context.Context, name string) (*model.UsersM, error) {
	var user model.UsersM
	if err := u.db.WithContext(ctx).Where("name = ?", name).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *users) Update(ctx context.Context, user *model.UsersM) error {
	if err := u.db.WithContext(ctx).Model(&model.UsersM{}).Where("name = ?", user.Name).Updates(user).Error; err != nil {
		return err
	}

	return nil
}

func (u *users) List(ctx context.Context, offset, limit int) (count int64, ret []*model.UsersM, err error) {
	err = u.db.Offset(offset).Limit(defaultLimit(limit)).Order("id desc").Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error

	return
}

func (u *users) Delete(ctx context.Context, name string) error {
	err := u.db.Where("name = ?", name).Delete(&model.UsersM{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}
