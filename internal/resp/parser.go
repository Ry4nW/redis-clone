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
	if len(tokens) == 0 {
		return nil, ErrEmptyRequest
	}

	// normalize cmd
	command := strings.ToUpper(tokens[0])
	args := tokens[1:]

	return &Request{
		Command: command,
		Args:    args,
	}, nil
}
