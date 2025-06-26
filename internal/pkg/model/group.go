package model

// GroupType
var (
	Group     = 0
	Broadcast = 1
)

type GroupUserM struct {
	ModelID
	ModelTimestamps
	UserID    int64   `gorm:"column:user_id" json:"user_id"`
	GroupID   *int64  `gorm:"column:group_id" json:"group_id,omitempty"`
	GroupType int16   `gorm:"column:group_type" json:"group_type"` // 0 group 1 broadcast
	Remark    *string `gorm:"column:remark" json:"remark,omitempty"`
	Avatar    *string `gorm:"column:avatar" json:"avatar,omitempty"`
	Name      *string `gorm:"column:name" json:"name,omitempty"`
}

type GroupM struct {
	ModelID
	ModelTimestamps
	GroupType int16   `gorm:"column:group_type" json:"group_type"`     // 0 group 1 broadcast
	UserID    *int64  `gorm:"column:user_id" json:"user_id,omitempty"` // creator
	Name      *string `gorm:"column:name" json:"name,omitempty"`       // group name
	Info      *string `gorm:"column:info" json:"info,omitempty"`       // group description
	Avatar    *string `gorm:"column:avatar" json:"avatar,omitempty"`   // group avatar
	Password  *string `gorm:"column:password" json:"password,omitempty"`
	IsPwd     int16   `gorm:"column:is_pwd" json:"is_pwd"`         // is encrypted 0 no 1 yes
	Hot       *int64  `gorm:"column:hot" json:"hot,omitempty"`     // popularity
	Theme     *string `gorm:"column:theme" json:"theme,omitempty"` // group theme
}
