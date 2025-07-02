package store

import (
	"context"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	"gorm.io/gorm"
)

type SessionStore interface {
	Create(ctx context.Context, message *model.SessionM) error
	GetByUserID(ctx context.Context, userId int64) (*[]model.SessionM, error)
	GetByUserIDAndType(ctx context.Context, userId, toId int64, channelType int16) (*model.SessionM, error)
	Update(ctx context.Context, sessionId int64, topStatus int16, Note string) error
	Delete(ctx context.Context, sessionId int64) error
}

type session struct {
	db *gorm.DB
}

var _ SessionStore = (*session)(nil)

func NewSession(db *gorm.DB) SessionStore {
	return &session{
		db: db,
	}
}

func (s *session) Create(ctx context.Context, message *model.SessionM) error {
	return s.db.WithContext(ctx).Create(message).Error
}

func (s *session) GetByUserID(ctx context.Context, userId int64) (*[]model.SessionM, error) {
	var sessions []model.SessionM
	if err := s.db.WithContext(ctx).Where("form_id = ?", userId).Find(&sessions).Error; err != nil {
		return nil, err
	}

	return &sessions, nil
}

func (s *session) GetByUserIDAndType(ctx context.Context, userId, toId int64, channelType int16) (*model.SessionM, error) {
	var session model.SessionM
	if err := s.db.WithContext(ctx).Where("form_id = ? and to_id = ? and channel_type = ?", userId, toId, channelType).First(&session).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *session) Update(ctx context.Context, sessionId int64, topStatus int16, Note string) error {
	session := &model.SessionM{
		TopStatus: topStatus,
		Note:      Note,
	}

	return s.db.WithContext(ctx).Model(session).Where("id = ?", sessionId).Updates(session).Error
}

func (s *session) Delete(ctx context.Context, sessionId int64) error {
	return s.db.WithContext(ctx).Where("id = ?", sessionId).Delete(&model.SessionM{}).Error
}
