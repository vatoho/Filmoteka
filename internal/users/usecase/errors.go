package usecase

import "errors"

var (
	ErrBadCredentials    = errors.New("bad auth data for user")
	ErrUserAlreadyExists = errors.New("user with such username already exists")
	ErrNoUser            = errors.New("user not exists")
)
