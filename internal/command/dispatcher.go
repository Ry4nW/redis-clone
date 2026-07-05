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
	S
}

func HandleEcho(req *resp.Request) (string, error) {
}
