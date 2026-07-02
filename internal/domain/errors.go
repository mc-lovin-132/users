package domain

import "errors"

var (
	ErrNotUniqueEmail  = errors.New("user with that email is already exists")
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrEmptyName       = errors.New("name is empty")
	ErrInvalidPassword = errors.New("invalid password format")
	ErrUserNotFound    = errors.New("user not found")
	ErrToMuchArgs      = errors.New("need only one arg")
	ErrNotEnoughArgs   = errors.New("not enough args")
	ErrInternal        = errors.New("internal error")
)
