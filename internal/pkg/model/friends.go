package model

const (
	FriendStatusPending  = 0 // pending
	FriendStatusAccepted = 1 // accepted
	FriendStatusRejected = 2 // rejected

	FriendStatusUnchecking = -1 // uncheck
)

const (
	FriendNotPinned = 0 // not pinned
	FriendPinned    = 1 // pinned
)

type FriendRecordM struct {
	ModelID
	ModelTimestamps
	FormID      int64  `gorm:"column:form_id" json:"form_id"`
	ToID        int64  `gorm:"column:to_id" json:"to_id"`
	Status      int16  `gorm:"column:status" json:"status"`                     // 0 pending 1 accepted 2 rejected
	Information string `gorm:"column:information" json:"information,omitempty"` // request info
	Users       IUsers `gorm:"foreignKey:FormId;references:Id" json:"users"`
}

type IUsers struct {
	Id     int64  `gorm:"column:id;primaryKey" json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func (f *FriendRecordM) TableName() string {
	return "friend_records"
}

type FriendM struct {
	ModelID
	ModelTimestamps
	FormID  int64  `gorm:"column:form_id" json:"form_id"`
	ToID    int64  `gorm:"column:to_id" json:"to_id"`
	Note    string `gorm:"column:note" json:"note,omitempty"`
	TopTime string `gorm:"column:top_time" json:"top_time,omitempty"`
	Status  int16  `gorm:"column:status" json:"status"` // 0 not pinned 1 pinned
	Uid     string `gorm:"column:uid" json:"uid"`
	Users   UsersM `gorm:"foreignKey:ID;references:ToID" json:"users,omitempty"` // User information
}

func (f *FriendM) TableName() string {
	return "friends"
}
