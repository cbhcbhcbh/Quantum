package model

type UsersM struct {
	ModelID
	ModelTimestamps
	Name          string `gorm:"column:name" json:"name"`
	Email         string `gorm:"column:email" json:"email"`
	Password      string `gorm:"column:password" json:"password"`
	Avatar        string `gorm:"column:avatar" json:"avatar"`
	Status        int8   `gorm:"column:status" json:"status"`
	Bio           string `gorm:"column:bio" json:"bio"`
	Sex           int8   `gorm:"column:sex" json:"sex"`
	ClientType    int8   `gorm:"column:client_type" json:"client_type"`
	Age           int    `gorm:"column:age" json:"age"`
	LastLoginTime string `gorm:"column:last_login_time" json:"last_login_time"`
	Uid           string `gorm:"column:uid" json:"uid"`
}

func (u *UsersM) Table() string {
	return "users"
}
