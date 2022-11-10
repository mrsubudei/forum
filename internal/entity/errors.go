package entity

import "errors"

var (
	ErrUserNotFound           = errors.New("user doesn't exist")
	ErrPostNotFound           = errors.New("posts wasn't found")
	ErrUserEmailAlreadyExists = errors.New("user with such email already exists")
	ErrUserNameAlreadyExists  = errors.New("user with such name already exists")
	ErrUserPasswordIncorrect  = errors.New("password is incorrect")
	ErrUserEmailIncorrect     = errors.New("email is incorrect")
)
