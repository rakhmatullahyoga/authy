package auth

import "errors"

var (
	// service errors
	ErrUserAlreadyRegistered = errors.New("already same user_id is used")
)
