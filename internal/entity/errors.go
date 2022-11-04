package entity

import "errors"

var (
	ErrUserNotFound           = errors.New("user doesn't exists")
	ErrUserEmailAlreadyExists = errors.New("user with such email already exists")
	ErrUserNameAlreadyExists  = errors.New("user with such name already exists")
	ErrUserPasswordIncorrect  = errors.New("password is incorrect")
)
