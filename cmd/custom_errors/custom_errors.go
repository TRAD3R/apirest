package custom_errors

import "errors"

var (
	ErrUserNotFound = errors.New("user is not found")
	ErrPostNotFound = errors.New("post is not found")
)
