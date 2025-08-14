package chat

import (
	"context"
	"fmt"

	"github.com/cbhcbhcbh/Quantum/pkg/common/domain"
	"github.com/cbhcbhcbh/Quantum/pkg/common/sonyflake"
)

type MessageService interface {
	BroadcastTextMessage(ctx context.Context, channelID, userID uint64, payload string) error
	BroadcastConnectMessage(ctx context.Context, channelID, userID uint64) error
	BroadcastActionMessage(ctx context.Context, channelID, userID uint64, action domain.Action) error
	BroadcastFileMessage(ctx context.Context, channelID, userID uint64, payload string) error
	MarkMessageSeen(ctx context.Context, channelID, userID, messageID uint64) error
	InsertMessage(ctx context.Context, msg *domain.Message) error
	PublishMessage(ctx context.Context, msg *domain.Message) error
	ListMessages(ctx context.Context, channelID uint64, pageState string) ([]*domain.Message, string, error)
}

type UserService interface {
	AddUserToChannel(ctx context.Context, channelID, userID uint64) error
	GetUser(ctx context.Context, userID uint64) (*domain.User, error)
	IsChannelUserExist(ctx context.Context, channelID, userID uint64) (bool, error)
	GetChannelUserIDs(ctx context.Context, channelID uint64) ([]uint64, error)
	AddOnlineUser(ctx context.Context, channelID, userID uint64) error
	DeleteOnlineUser(ctx context.Context, channelID, userID uint64) error
	GetOnlineUserIDs(ctx context.Context, channelID uint64) ([]uint64, error)
}

type ChannelService interface {
	CreateChannel(ctx context.Context) (*domain.Channel, error)
	DeleteChannel(ctx context.Context, channelID uint64) error
}

type ForwardService interface {
	RegisterChannelSession(ctx context.Context, channelID, userID uint64, subscriber string) error
	RemoveChannelSession(ctx context.Context, channelID, userID uint64) error
}

type UserServiceImpl struct {
	userRepo UserRepoCache
}

func NewUserServiceImpl(userRepo UserRepoCache) *UserServiceImpl {
	return &UserServiceImpl{userRepo}
}

func (svc *UserServiceImpl) AddUserToChannel(ctx context.Context, channelID, userID uint64) error {
	if err := svc.userRepo.AddUserToChannel(ctx, channelID, userID); err != nil {
		return fmt.Errorf("error add user %d to channel %d: %w", userID, channelID, err)
	}
	return nil
}

func (svc *UserServiceImpl) GetUser(ctx context.Context, userID uint64) (*domain.User, error) {
	user, err := svc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error get user %d: %w", userID, err)
	}
	return user, nil
}

func (svc *UserServiceImpl) IsChannelUserExist(ctx context.Context, channelID, userID uint64) (bool, error) {
	exist, err := svc.userRepo.IsChannelUserExist(ctx, channelID, userID)
	if err != nil {
		return false, fmt.Errorf("error check user %d in channel %d: %w", userID, channelID, err)
	}
	return exist, nil
}

func (svc *UserServiceImpl) GetChannelUserIDs(ctx context.Context, channelID uint64) ([]uint64, error) {
	users, err := svc.userRepo.GetChannelUserIDs(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("error get users in channel %d: %w", channelID, err)
	}
	return users, nil
}

func (svc *UserServiceImpl) AddOnlineUser(ctx context.Context, channelID, userID uint64) error {
	if err := svc.userRepo.AddOnlineUser(ctx, channelID, userID); err != nil {
		return fmt.Errorf("error add online user %d to channel %d: %w", userID, channelID, err)
	}
	return nil
}

func (svc *UserServiceImpl) DeleteOnlineUser(ctx context.Context, channelID, userID uint64) error {
	if err := svc.userRepo.DeleteOnlineUser(ctx, channelID, userID); err != nil {
		return fmt.Errorf("error delete online user %d from channel %d: %w", userID, channelID, err)
	}
	return nil
}

func (svc *UserServiceImpl) GetOnlineUserIDs(ctx context.Context, channelID uint64) ([]uint64, error) {
	users, err := svc.userRepo.GetOnlineUserIDs(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("error get online users in channel %d: %w", channelID, err)
	}
	return users, nil
}

type ChannelServiceImpl struct {
	channelRepo ChannelRepoCache
	userRepo    UserRepoCache
	sf          sonyflake.IDGenerator
}

func NewChannelServiceImpl(chanRepo ChannelRepoCache, userRepo UserRepoCache, sf sonyflake.IDGenerator) *ChannelServiceImpl {
	return &ChannelServiceImpl{chanRepo, userRepo, sf}
}

func (svc *ChannelServiceImpl) CreateChannel(ctx context.Context) (*domain.Channel, error) {
	channelID, err := svc.sf.NextID()
	if err != nil {
		return nil, fmt.Errorf("error create snowflake ID for new channel: %w", err)
	}
	channel, err := svc.channelRepo.CreateChannel(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("error create channel %d: %w", channelID, err)
	}
	return channel, nil
}

func (svc *ChannelServiceImpl) DeleteChannel(ctx context.Context, channelID uint64) error {
	if err := svc.channelRepo.DeleteChannel(ctx, channelID); err != nil {
		return fmt.Errorf("error delete channel %d: %w", channelID, err)
	}
	return nil
}

type ForwardServiceImpl struct {
	forwardRepo ForwardRepo
}

func NewForwardServiceImpl(forwardRepo ForwardRepo) *ForwardServiceImpl {
	return &ForwardServiceImpl{forwardRepo}
}

func (svc *ForwardServiceImpl) RegisterChannelSession(ctx context.Context, channelID, userID uint64, subscriber string) error {
	return svc.forwardRepo.RegisterChannelSession(ctx, channelID, userID, subscriber)
}
func (svc *ForwardServiceImpl) RemoveChannelSession(ctx context.Context, channelID, userID uint64) error {
	return svc.forwardRepo.RemoveChannelSession(ctx, channelID, userID)
}
