package users

import (
	"errors"
	"net/http"
	"time"

	"github.com/cbhcbhcbh/Quantum/internal/apiserver/store"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/date"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/enum"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/helpers"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	v1 "github.com/cbhcbhcbh/Quantum/pkg/api/v1"
	"github.com/cbhcbhcbh/Quantum/pkg/auth"
	"github.com/cbhcbhcbh/Quantum/pkg/jwt"
	"github.com/cbhcbhcbh/Quantum/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// TODO: Add comments to the UserBiz interface and its methods.
type UserBiz interface {
	Login(ctx *gin.Context, r *v1.LoginRequest) (*v1.LoginResponse, error)
	Registered(ctx *gin.Context, r *v1.RegisterUserRequest) error
}

type userBiz struct {
	ds store.IStore
}

var _ UserBiz = (*userBiz)(nil)

func New(ds store.IStore) UserBiz {
	return &userBiz{ds: ds}
}

// TODO: 统一 Context
func (b *userBiz) Login(ctx *gin.Context, r *v1.LoginRequest) (*v1.LoginResponse, error) {
	c := ctx.Request.Context()

	user, err := b.ds.Users().Get(c, r.Name)
	if err != nil {
		return nil, err
	}

	if !auth.BcryptCheck(user.Password, r.Password) {
		response.FailResponse(http.StatusInternalServerError, "").ToJson(ctx)
		return nil, errors.New("password was incorrect")
	}

	ttl := viper.GetInt64("jwt.ttl")
	expireAtTime := time.Now().Unix() + ttl
	token := jwt.NewJWT().IssueToken(
		user.ModelID.ID,
		user.Uid,
		user.Name,
		user.Email,
		expireAtTime,
	)

	return &v1.LoginResponse{
		ID:         user.ID,
		UID:        user.Uid,
		Name:       user.Name,
		Avatar:     user.Avatar,
		Email:      user.Email,
		ExpireTime: expireAtTime,
		Token:      token,
		Ttl:        ttl,
	}, nil
}

func (b *userBiz) Registered(ctx *gin.Context, r *v1.RegisterUserRequest) error {
	c := ctx.Request.Context()

	userm := model.UsersM{
		Name:  r.Name,
		Email: r.Email,
	}

	ok, err := b.ds.Users().CheckUserExist(c, userm.Name, userm.Email)
	if ok {
		response.FailResponse(enum.ParamError, err.Error()).WriteTo(ctx)
		return err
	}

	// TODO: Check 邮件验证码

	userm.Password = auth.BcryptHash(userm.Password)
	userm.LastLoginTime = date.NewDate()
	userm.Uid = helpers.GetUuid()

	err = b.ds.Users().Create(c, &userm)
	if err != nil {
		response.FailResponse(enum.DBError, err.Error()).ToJson(ctx)
		return err
	}

	response.SuccessResponse().ToJson(ctx)
	return nil
}
