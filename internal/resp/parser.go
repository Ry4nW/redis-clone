package resp

import (
	"bufio"
	"strings"
)

func Parse(reader *bufio.Reader) (*Request, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	tokens := strings.Fields(line)
	command := tokens[0]
	args := tokens[1:]

	newReq := &Request{
		Command: command,
		Args: args,
	}

	return newReq, nil
}
