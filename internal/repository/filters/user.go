package filters

import "time"

type UserFilter struct {
	Offset         uint
	Limit          uint
	FromCreatedAt  *time.Time
	ToCreatedAt    *time.Time
	Name           []string
	TopPostsAmount string
}

type UserAddRequest struct {
	Name        string `json:"name"`
	Phonenumber string `json:"phonenumber"`
}

type UserUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Phonenumber string `json:"phonenumber,omitempty"`
}
