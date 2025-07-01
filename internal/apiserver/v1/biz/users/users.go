package users

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/store"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/date"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/enum"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/helpers"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	"github.com/cbhcbhcbh/Quantum/internal/service/email"
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
	SendEmail(ctx *gin.Context, r *v1.SendEmailRequest) error
	GetUserInfo(ctx *gin.Context, r *v1.Person) (*v1.UserDetails, error)
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

	user, err := b.ds.Users().GetByName(c, r.Name)
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

func (b *userBiz) SendEmail(ctx *gin.Context, r *v1.SendEmailRequest) error {
	ok := b.ds.Users().IsTableFliedExits(ctx, "email", r.Email)

	switch r.EmailType {

	case email.REGISTERED_CODE:
		if ok {
			return errors.New("邮箱已经被注册了")
		}

	case email.RESET_PS_CODE:
		if !ok {
			return errors.New("邮箱未注册了")
		}

	}

	emailService := email.NewEmailService()

	code := helpers.CreateEmailCode()

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Im-Services邮件验证码</title>
</head>
<style>
    .mail{
        margin: 0 auto;
        border-radius: 45px;
        height: 400px;
        padding: 10px;
        background-color: #CC9933;
        background: url("https://img-blog.csdnimg.cn/c32f12dfd48241babd35b15189dc5c78.png") no-repeat;
    }
    .code {
        color: #f6512b;
        font-weight: bold;
        font-size: 30px;
        padding: 2px;
    }
</style>
<body>
<div class="mail">
    <h3>您好 ~ im-services应用账号!</h3>
    <p>下面是您的验证码:</p>
        <p class="code">%s</p>
        <p>请注意查收!谢谢</p>
</div>
<h3>如果可以请给项目点个star～<a target="_blank" href="https://github.com/IM-Tools/Im-Services">项目地址</a> </h3>
</body>
</html>`, code)

	subject := "欢迎使用～👏Im Services,这是一封邮箱验证码的邮件!🎉🎉🎉"

	err := emailService.SendEmail(code, r.EmailType, r.Email, subject, html)
	if err != nil {
		log.C(ctx).Errorw("发送失败邮箱:" + r.Email + "错误日志:" + err.Error())
		response.FailResponse(enum.ApiError, "邮件发送失败,请检查是否是可用邮箱").ToJson(ctx)
		return err
	}

	response.SuccessResponse().ToJson(ctx)
	return nil

}

func (b *userBiz) GetUserInfo(ctx *gin.Context, r *v1.Person) (*v1.UserDetails, error) {
	c := ctx.Request.Context()

	id, _ := strconv.ParseInt(r.ID, 10, 64)
	user, err := b.ds.Users().GetById(c, id)
	if err != nil {
		return nil, err
	}

	return &v1.UserDetails{
		ID:            user.ID,
		Uid:           user.Uid,
		Name:          user.Name,
		Avatar:        user.Avatar,
		Email:         user.Email,
		Status:        user.Status,
		Bio:           user.Bio,
		Sex:           user.Sex,
		Age:           user.Age,
		LastLoginTime: user.LastLoginTime,
	}, nil
}
