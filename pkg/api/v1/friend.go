package v1

type FriendDetail struct {
	ID     int64   `json:"to_id" validate:"required"`
	Note   *string `json:"note" validate:"omitempty,note"`
	Status int16   `json:"status" validate:"omitempty,oneof=0 1"` // 0 not pinned 1 pinned
	Uid    string  `json:"uid" validate:"required"`
}
