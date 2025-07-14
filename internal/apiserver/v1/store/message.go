package store

import (
	"context"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	"gorm.io/gorm"
)

type MessageStore interface {
	ListByFormIdAndToId(ctx context.Context, formId, toId int64) ([]*model.MessageM, int64, error)
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

func (m *message) ListByFormIdAndToId(ctx context.Context, formId, toId int64) ([]*model.MessageM, int64, error) {
	var total int64
	var list []*model.MessageM

	query := m.db.WithContext(ctx).
		Where("(form_id = ? and to_id = ?) or (form_id = ? and to_id = ?)", formId, toId, formId, toId).
		Order("created_at desc")

	if err := query.Count(&total).Error; err != nil {
		return nil, total, err
	}

	if err := query.Find(&list).Error; err != nil {
		return nil, total, err
	}

	return list, total, nil
}

func (m *message) Create(ctx context.Context, message *model.MessageM) error {
	if err := m.db.WithContext(ctx).Create(message).Error; err != nil {
		return err
	}
	return nil
}