package session

import (
	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/store"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/date"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	v1 "github.com/cbhcbhcbh/Quantum/pkg/api/v1"
	"github.com/gin-gonic/gin"
)

type SessionBiz interface {
	GetSessions(ctx *gin.Context, formId int64) (*[]v1.SessionDetail, error)
	CreateSession(ctx *gin.Context, formId, toId int64, channelType int16) (*model.SessionM, error)
	UpdateSession(ctx *gin.Context, sessionId int64, topStatus int16, note string) error
	DeleteSession(ctx *gin.Context, sessionId int64) error
	CreateSessionRelation(ctx *gin.Context, formId, toId int64, channelType int16, user *v1.UserDetails)
}

type sessionBiz struct {
	ds store.IStore
}

var _ SessionBiz = (*sessionBiz)(nil)

func New(ds store.IStore) SessionBiz {
	return &sessionBiz{ds: ds}
}

func (s *sessionBiz) GetSessions(ctx *gin.Context, userId int64) (*[]v1.SessionDetail, error) {
	c := ctx.Request.Context()

	sessions, err := s.ds.Sessions().GetByUserID(c, userId)
	if err != nil {
		return nil, err
	}

	var SessionDetails []v1.SessionDetail
	for _, session := range *sessions {
		SessionDetails = append(SessionDetails, v1.SessionDetail{
			ID:          session.ID,
			FormID:      session.FormID,
			ToID:        session.ToID,
			TopStatus:   session.TopStatus,
			TopTime:     session.TopTime,
			Note:        session.Note,
			ChannelType: session.ChannelType,
			Name:        session.Name,
			Avatar:      session.Avatar,
			Status:      session.Status,
			GroupID:     session.GroupID,
		})
	}

	return &SessionDetails, nil
}

func (s *sessionBiz) CreateSession(ctx *gin.Context, formId, toId int64, channelType int16) (*model.SessionM, error) {
	c := ctx.Request.Context()

	session, err := s.ds.Sessions().GetByUserIDAndType(c, formId, toId, channelType)
	if err != nil {
		return nil, err
	}

	if session != nil {
		return session, nil
	}

	user, _ := s.ds.Users().GetById(c, toId)

	newSession := model.SessionM{
		FormID:      formId,
		ToID:        toId,
		ChannelType: channelType,
		TopStatus:   model.TopStatus,
		TopTime:     date.NewDate(),
		Note:        user.Name,
		Name:        user.Name,
		Avatar:      user.Avatar,
		Status:      model.SessionStatusOk,
	}

	err = s.ds.Sessions().Create(c, &newSession)
	if err != nil {
		return nil, err
	}

	return &newSession, nil
}

func (s *sessionBiz) UpdateSession(ctx *gin.Context, sessionId int64, topStatus int16, note string) error {
	c := ctx.Request.Context()

	return s.ds.Sessions().Update(c, sessionId, topStatus, note)
}

func (s *sessionBiz) DeleteSession(ctx *gin.Context, sessionId int64) error {
	c := ctx.Request.Context()

	return s.ds.Sessions().Delete(c, sessionId)
}

// FIXME: Fix User Info
func (s *sessionBiz) CreateSessionRelation(ctx *gin.Context, formId, toId int64, channelType int16, user *v1.UserDetails) {
	sessionRelation1 := model.SessionM{
		ToID:        toId,
		FormID:      formId,
		TopStatus:   model.TopStatus,
		TopTime:     date.NewDate(),
		Note:        user.Name,
		ChannelType: channelType,
		Name:        user.Name,
		Avatar:      user.Avatar,
		Status:      model.SessionStatusOk,
	}

	sessionRelation2 := model.SessionM{
		ToID:        formId,
		FormID:      toId,
		TopStatus:   model.TopStatus,
		TopTime:     date.NewDate(),
		Note:        user.Name,
		ChannelType: channelType,
		Name:        user.Name,
		Avatar:      user.Avatar,
		Status:      model.SessionStatusOk,
	}

	_ = s.ds.Sessions().Create(ctx, &sessionRelation1)
	_ = s.ds.Sessions().Create(ctx, &sessionRelation2)
}
