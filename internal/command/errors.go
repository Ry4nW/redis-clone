package command

import "errors"

var (
	ErrBadArgAmt      = errors.New("bad arg amount")
	ErrBadArgType     = errors.New("bad arg type")
	ErrUnknownCommand = errors.New("unknown command")
)
