package v1

type FriendDetail struct {
	ID     int64       `json:"to_id" validate:"required"`
	Note   string      `json:"note" validate:"omitempty,note"`
	Status int16       `json:"status" validate:"omitempty,oneof=0 1"` // 0 not pinned 1 pinned
	Uid    string      `json:"uid" validate:"required"`
	Users  UserDetails `json:"users"`
}

type FriendInfo struct {
	FriendId int64 `uri:"id" binding:"required"`
}

type CreateFriendRequest struct {
	ToID        int64  `json:"to_id" validate:"required"`
	Information string `json:"information" validate:"omitempty,note"` // request info
}

type UpdateFriendRequest struct {
	ID     int64 `json:"id" validate:"required"`
	Status int16 `json:"status" validate:"required,gte=1,lte=2"`
}
