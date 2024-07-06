package filters

import "time"

type PostFilter struct {
	Offset        int
	Limit         int
	FromCreatedAt time.Time
	ToCreatedAt   time.Time
	Subject       string
	Authors       []int
}

type PostAddRequest struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Author  int    `json:"author"`
}

type PostUpdateRequest struct {
	Subject string `json:"subject,omitempty"`
	Body    string `json:"body,omitempty"`
}
