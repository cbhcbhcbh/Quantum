package chat

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-kit/kit/endpoint"
	"github.com/gocql/gocql"

	"github.com/cbhcbhcbh/Quantum/pkg/common/domain"
	"github.com/cbhcbhcbh/Quantum/pkg/common/jwt"
	"github.com/cbhcbhcbh/Quantum/pkg/common/known"
	"github.com/cbhcbhcbh/Quantum/pkg/config"
	forwarderpb "github.com/cbhcbhcbh/Quantum/pkg/proto/forwarder"
	userpb "github.com/cbhcbhcbh/Quantum/pkg/proto/user"
	"github.com/cbhcbhcbh/Quantum/pkg/transport"
)

var (
	MessagePubTopic = "rc.msg.pub"
)

type UserRepo interface {
	AddUserToChannel(ctx context.Context, channelID uint64, userID uint64) error
	GetUserByID(ctx context.Context, userID uint64) (*domain.User, error)
	GetChannelUserIDs(ctx context.Context, channelID uint64) ([]uint64, error)
}

type MessageRepo interface {
	InsertMessage(ctx context.Context, msg *domain.Message) error
	MarkMessageSeen(ctx context.Context, channelID, messageID uint64) error
	PublishMessage(ctx context.Context, msg *domain.Message) error
	ListMessages(ctx context.Context, channelID uint64, pageStateBase64 string) ([]*domain.Message, string, error)
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
		return nil, known.ErrUserNotFound
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

type MessageRepoImpl struct {
	s          *gocql.Session
	p          message.Publisher
	maxMessage int64
	pagination int
}

func (repo *MessageRepoImpl) InsertMessage(ctx context.Context, msg *domain.Message) error {
	var messageNum int64
	err := repo.s.Query("SELECT msgnum FROM chanmsg_counters WHERE channel_id = ? LIMIT 1", msg.ChannelID).
		WithContext(ctx).Idempotent(true).Scan(&messageNum)
	if err != nil {
		if err == gocql.ErrNotFound {
			messageNum = 0
		} else {
			return err
		}
	}
	if messageNum >= repo.maxMessage {
		return known.ErrExceedMessageNumLimits
	}
	if err := repo.s.Query("INSERT INTO messages (id, event, channel_id, user_id, payload, seen, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?)",
		msg.MessageID,
		msg.Event,
		msg.ChannelID,
		msg.UserID,
		msg.Payload,
		false,
		msg.Time).WithContext(ctx).Exec(); err != nil {
		return err
	}
	return repo.s.Query("UPDATE chanmsg_counters SET msgnum = msgnum + 1 WHERE channel_id = ?", msg.ChannelID).WithContext(ctx).Exec()
}

func (repo *MessageRepoImpl) MarkMessageSeen(ctx context.Context, channelID, messageID uint64) error {
	if err := repo.s.Query("UPDATE messages SET seen = ? WHERE channel_id = ? AND id = ?", true, channelID, messageID).
		WithContext(ctx).Idempotent(true).Exec(); err != nil {
		return err
	}
	return nil
}

func (repo *MessageRepoImpl) PublishMessage(ctx context.Context, msg *domain.Message) error {
	return repo.p.Publish(MessagePubTopic, message.NewMessage(watermill.NewUUID(), msg.Encode()))
}

func (repo *MessageRepoImpl) ListMessages(ctx context.Context, channelID uint64, pageStateBase64 string) ([]*domain.Message, string, error) {
	var messages []*domain.Message
	pageState, err := base64.URLEncoding.DecodeString(pageStateBase64)
	if err != nil {
		return nil, "", err
	}
	iter := repo.s.Query(`SELECT id, event, channel_id, user_id, payload, seen, timestamp FROM messages WHERE channel_id = ?`, channelID).
		WithContext(ctx).Idempotent(true).PageSize(repo.pagination).PageState(pageState).Iter()
	nextPageStateBase64 := base64.URLEncoding.EncodeToString(iter.PageState())
	scanner := iter.Scanner()

	for scanner.Next() {
		var message domain.Message
		if err = scanner.Scan(
			&message.MessageID,
			&message.Event,
			&message.ChannelID,
			&message.UserID,
			&message.Payload,
			&message.Seen,
			&message.Time); err != nil {
			return nil, "", err
		}
		messages = append(messages, &message)
	}
	err = scanner.Err()
	if err != nil {
		return nil, "", err
	}
	return messages, nextPageStateBase64, nil
}

type ChannelRepoImpl struct {
	s *gocql.Session
}

func NewMessageRepoImpl(config *config.Config, s *gocql.Session, p message.Publisher) *MessageRepoImpl {
	return &MessageRepoImpl{s, p, config.Chat.Message.MaxNum, config.Chat.Message.PaginationNum}
}

func NewChannelRepoImpl(s *gocql.Session) *ChannelRepoImpl {
	return &ChannelRepoImpl{s}
}

func (repo *ChannelRepoImpl) CreateChannel(ctx context.Context, channelID uint64) (*domain.Channel, error) {
	if err := repo.s.Query("INSERT INTO channels (id, user_id) VALUES (?, ?)",
		channelID, 0).WithContext(ctx).Exec(); err != nil {
		return nil, err
	}
	accessToken, err := jwt.NewJWT(channelID)
	if err != nil {
		return nil, fmt.Errorf("error create JWT: %w", err)
	}
	return &domain.Channel{
		ID:          channelID,
		AccessToken: accessToken,
	}, nil
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
