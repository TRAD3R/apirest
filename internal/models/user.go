package models

import "time"

type User struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Phonenumber string     `json:"phonenumber"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	PostCount   int
}
