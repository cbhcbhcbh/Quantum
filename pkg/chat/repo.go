package chat

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/gocql/gocql"

	"github.com/cbhcbhcbh/Quantum/pkg/common/domain"
	forwarderpb "github.com/cbhcbhcbh/Quantum/pkg/proto/forwarder"
	userpb "github.com/cbhcbhcbh/Quantum/pkg/proto/user"
	"github.com/cbhcbhcbh/Quantum/pkg/transport"
)

type UserRepo interface {
	AddUserToChannel(ctx context.Context, channelID uint64, userID uint64) error
	GetUserByID(ctx context.Context, userID uint64) (*domain.User, error)
	GetChannelUserIDs(ctx context.Context, channelID uint64) ([]uint64, error)
}

type ChannelRepo interface {
	CreateChannel(ctx context.Context, channelID uint64) (*domain.Channel, error)
	DeleteChannel(ctx context.Context, channelID uint64) error
}

type ForwardRepo interface {
	RegisterChannelSession(ctx context.Context, channelID, userID uint64, subscriber string) error
	RemoveChannelSession(ctx context.Context, channelID, userID uint64) error
}

type UserRepoImpl struct {
	s       *gocql.Session
	getUser endpoint.Endpoint
}

func NewUserRepoImpl(s *gocql.Session, userConn *UserClientConn) *UserRepoImpl {
	return &UserRepoImpl{
		s: s,
		getUser: transport.NewGrpcEndpoint(
			userConn.Conn,
			"user.UserService",
			"GetUser",
			&userpb.GetUserResponse{},
		),
	}
}

func (repo *UserRepoImpl) AddUserToChannel(ctx context.Context, channelID uint64, userID uint64) error {
	if err := repo.s.Query("INSERT INTO channel_users (channel_id, user_id) VALUES (?, ?)", channelID, userID).WithContext(ctx).Exec(); err != nil {
		return err
	}
	return nil
}

func (repo *UserRepoImpl) GetUserByID(ctx context.Context, userID uint64) (*domain.User, error) {
	res, err := repo.getUser(ctx, &userpb.GetUserRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	pbUser := res.(*userpb.GetUserResponse)
	if !pbUser.Exist {
		return nil, ErrUserNotFound
	}
	return &domain.User{
		ID:   pbUser.User.Id,
		Name: pbUser.User.Name,
	}, nil
}

func (repo *UserRepoImpl) GetChannelUserIDs(ctx context.Context, channelID uint64) ([]uint64, error) {
	iter := repo.s.Query("SELECT user_id FROM channel_users WHERE channel_id = ?", channelID).WithContext(ctx).Idempotent(true).Iter()
	var userIDs []uint64
	var userID uint64
	for iter.Scan(&userID) {
		userIDs = append(userIDs, userID)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return userIDs, nil
}

type ChannelRepoImpl struct {
	s *gocql.Session
}

func NewChannelRepoImpl(s *gocql.Session) *ChannelRepoImpl {
	return &ChannelRepoImpl{s}
}

func (repo *ChannelRepoImpl) CreateChannel(ctx context.Context, channelID uint64) (*domain.Channel, error) {
	if err := repo.s.Query("INSERT INTO channels (id, user_id) VALUES (?, ?)",
		channelID, 0).WithContext(ctx).Exec(); err != nil {
		return nil, err
	}
	// TODO: Implement channel creation logic
	return nil, nil
}

func (repo *ChannelRepoImpl) DeleteChannel(ctx context.Context, channelID uint64) error {
	if err := repo.s.Query("DELETE FROM channels WHERE id = ?", channelID).
		WithContext(ctx).Exec(); err != nil {
		return err
	}
	return nil
}

type ForwardRepoImpl struct {
	registerChannelSession endpoint.Endpoint
	removeChannelSession   endpoint.Endpoint
}

func NewForwardRepoImpl(forwarderConn *ForwarderClientConn) *ForwardRepoImpl {
	return &ForwardRepoImpl{
		registerChannelSession: transport.NewGrpcEndpoint(
			forwarderConn.Conn,
			"forwarder.ForwardService",
			"RegisterChannelSession",
			&forwarderpb.RegisterChannelSessionResponse{},
		),
		removeChannelSession: transport.NewGrpcEndpoint(
			forwarderConn.Conn,
			"forwarder.ForwardService",
			"RemoveChannelSession",
			&forwarderpb.RemoveChannelSessionResponse{},
		),
	}
}

func (repo *ForwardRepoImpl) RegisterChannelSession(ctx context.Context, channelID, userID uint64, subscriber string) error {
	if _, err := repo.registerChannelSession(ctx, &forwarderpb.RegisterChannelSessionRequest{
		ChannelId:  channelID,
		UserId:     userID,
		Subscriber: subscriber,
	}); err != nil {
		return err
	}
	return nil
}
func (repo *ForwardRepoImpl) RemoveChannelSession(ctx context.Context, channelID, userID uint64) error {
	if _, err := repo.removeChannelSession(ctx, &forwarderpb.RemoveChannelSessionRequest{
		ChannelId: channelID,
		UserId:    userID,
	}); err != nil {
		return err
	}
	return nil
}
