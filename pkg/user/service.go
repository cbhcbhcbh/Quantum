package user

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cbhcbhcbh/Quantum/pkg/common/domain"
	"github.com/cbhcbhcbh/Quantum/pkg/common/util"
)

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v3/userinfo?access_token="

type UserService interface {
	GetGoogleUser(ctx context.Context, code string) (*domain.GoogleUserPresenter, error)
	SetUserSession(ctx context.Context, uid uint64) (string, error)
	GetUserByID(ctx context.Context, uid uint64) (*User, error)
	GetUserIDBySession(ctx context.Context, sid string) (uint64, error)
}

type UserServiceImpl struct {
	userRepo UserRepo
}

func NewUserServiceImpl(userRepo UserRepo) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

func (svc *UserServiceImpl) GetGoogleUser(ctx context.Context, code string) (*domain.GoogleUserPresenter, error) {
	req, err := http.NewRequest("GET", util.Join(oauthGoogleUrlAPI, code), nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var googleUser domain.GoogleUserPresenter
	if err := json.Unmarshal(contents, &googleUser); err != nil {
		return nil, err
	}
	return &googleUser, nil
}

func (svc *UserServiceImpl) SetUserSession(ctx context.Context, uid uint64) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("eror create sid: %w", err)
	}
	sid := base64.URLEncoding.EncodeToString(b)
	if err := svc.userRepo.SetUserSession(ctx, uid, sid); err != nil {
		return "", fmt.Errorf("eror set sid for user %d: %w", uid, err)
	}
	return sid, nil
}

func (svc *UserServiceImpl) GetUserByID(ctx context.Context, uid uint64) (*User, error) {
	user, err := svc.userRepo.GetUserByID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("error get user %d: %w", uid, err)
	}
	return user, nil
}

func (svc *UserServiceImpl) GetUserIDBySession(ctx context.Context, sid string) (uint64, error) {
	userID, err := svc.userRepo.GetUserIDBySession(ctx, sid)
	if err != nil {
		return 0, fmt.Errorf("error get user id by sid %s: %w", sid, err)
	}
	return userID, nil
}
