package domain

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInsufficientLimit = errors.New("insufficient daily limit")
	ErrSelfTransfer      = errors.New("cannot transfer to yourself")
	ErrNoUsersInChat     = errors.New("no other users in chat")
	ErrBetAlreadyUsed    = errors.New("bet already used today")
	ErrNotAuthorized     = errors.New("not authorized")
)
