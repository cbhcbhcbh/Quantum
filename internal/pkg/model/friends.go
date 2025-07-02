package model

type FriendRecordM struct {
	ModelID
	ModelTimestamps
	FormID      int64   `gorm:"column:form_id" json:"form_id"`
	ToID        int64   `gorm:"column:to_id" json:"to_id"`
	Status      int16   `gorm:"column:status" json:"status"`                     // 0 pending 1 accepted 2 rejected
	Information *string `gorm:"column:information" json:"information,omitempty"` // request info
}

type FriendM struct {
	ModelID
	ModelTimestamps
	FormID  int64   `gorm:"column:form_id" json:"form_id"`
	ToID    int64   `gorm:"column:to_id" json:"to_id"`
	Note    *string `gorm:"column:note" json:"note,omitempty"`
	TopTime *string `gorm:"column:top_time" json:"top_time,omitempty"`
	Status  int16   `gorm:"column:status" json:"status"` // 0 not pinned 1 pinned
	Uid     string  `gorm:"column:uid" json:"uid"`
	Users   UsersM  `gorm:"foreignKey:ID;references:ToID" json:"users,omitempty"` // User information
}
