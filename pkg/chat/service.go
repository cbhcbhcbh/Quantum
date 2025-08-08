package chat

import (
	"context"
	"fmt"

	"github.com/cbhcbhcbh/Quantum/pkg/common/domain"
)

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

type UserServiceImpl struct {
	userRepo UserRepoCache
}

func NewUserServiceImpl(userRepo UserRepoCache) *UserServiceImpl {
	return &UserServiceImpl{userRepo}
}

func (svc *UserServiceImpl) AddUserToChannel(ctx context.Context, channelID, userID uint64) error {}

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

func (svc *UserServiceImpl) AddOnlineUser(ctx context.Context, channelID, userID uint64) error {}

func (svc *UserServiceImpl) DeleteOnlineUser(ctx context.Context, channelID, userID uint64) error {}

func (svc *UserServiceImpl) GetOnlineUserIDs(ctx context.Context, channelID uint64) ([]uint64, error) {
}
