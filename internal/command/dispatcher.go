package command

import (
	"redis-clone/internal/resp"
)

func Dispatch(req *resp.Request) (string, error) {
	switch req.Command {
	case "PING":
		return HandlePing(req)
	case "ECHO":
		return HandleEcho(req)
	default:
		return "", ErrUnknownCommand
	}
}

func HandlePing(req *resp.Request) (string, error) {
	// 0 arg
	lenArgs := len(req.Args)
	if lenArgs == 0 {
		return "PONG\n", nil
	}

	if lenArgs > 1 {
		return "", ErrBadArgAmt
	}

	// 1 arg is echo
	return HandleEcho(req)
}

func HandleEcho(req *resp.Request) (string, error) {
	if lenArgs := len(req.Args); lenArgs > 1 || lenArgs == 0 {
		return "", ErrBadArgAmt
	}

	return req.Args[0] + "\n", nil
}
