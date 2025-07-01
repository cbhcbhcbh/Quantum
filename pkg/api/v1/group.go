package v1

type GroupDetail struct {
	ID        int64  `json:"id" validate:"required"`
	GroupType int16  `json:"group_type" validate:"omitempty"`
	Name      string `json:"name" validate:"omitempty"`
	Info      string `json:"info" validate:"required"`
	Avatar    string `json:"avatar" validate:"required"`
	Theme     string `json:"theme" validate:"required"`
}
