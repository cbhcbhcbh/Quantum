package v1

type RegisterUserRequest struct {
	Name           string `json:"name" valid:"alphanum,required,stringlength(1|255)"`
	Password       string `json:"password" valid:"required,stringlength(6|18)"`
	PasswordRepeat string `validate:"required,eqcsfield=Password"`
	EmailType      int    `validate:"required,gte=1,lte=2"`
	Email          string `json:"email" valid:"required,email"`
	Phone          string `json:"phone" valid:"required,stringlength(11|11)"`
	Code           string `validate:"required,len=4"`
}

type LoginRequest struct {
	Name     string `json:"name" valid:"alphanum,required,stringlength(1|255)"`
	Password string `json:"password" valid:"required,stringlength(6|18)"`
}

type LoginResponse struct {
	ID         int64  `json:"id"`
	UID        string `json:"uid"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
	Email      string `json:"email"`
	Token      string `json:"token"`
	ExpireTime int64  `json:"expire_time"`
	Ttl        int64  `json:"ttl"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" valid:"required,stringlength(6|18)"`
	NewPassword string `json:"newPassword" valid:"required,stringlength(6|18)"`
}

type GetUserResponse UserInfo

type UserInfo struct {
	Name      string `json:"name"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	PostCount int64  `json:"postCount"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ListUserRequest struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

type ListUserResponse struct {
	TotalCount int64       `json:"totalCount"`
	Users      []*UserInfo `json:"users"`
}

type UpdateUserRequest struct {
	Nickname *string `json:"nickname" valid:"stringlength(1|255)"`
	Email    *string `json:"email" valid:"email"`
	Phone    *string `json:"phone" valid:"stringlength(11|11)"`
}
