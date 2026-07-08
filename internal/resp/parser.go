package resp

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

type Request struct {
	Command string
	Args    []string
}

type Data interface{}

type BulkStr struct {
	Data []byte
}

type Arr struct {
	Data []Data
}

func Parse(reader *bufio.Reader) (Data, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	switch line[0] {
	case '$':
		parseBulkStr(line, reader)
	case '*':
		parseArray(line, reader)
	default:
		return nil, ErrMalformedRESP
	}
}

func parseBulkStr(line string, reader *bufio.Reader) (*BulkStr, error) {
	// $<length>\r\n<data>\r\n
	strLen, err := getLenHelper(line)

	if err != nil {
		return nil, err
	}

	data := make([]byte, strLen)
	io.ReadFull(reader, data)
	crlf := make([]byte, 2)
	io.ReadFull(reader, crlf)

	return &BulkStr{
		Data: data,
	}, nil
}

func parseArray(line string, reader *bufio.Reader) (*Arr, error) {
	// *<number-of-elements>\r\n<element-1>...<element-n>
	lengthStr := strings.TrimSuffix(line, "\r\n")[1:]
	arrLen, err := strconv.Atoi(lengthStr)

	if err != nil {
		return nil, err
	}

	// garbage buffer for crlfs
	crlfGarb := make([]byte, 2)
	result := make([]Data, arrLen)
	// each arr el has its own crlf
	for i := 0; i < arrLen; i++ {
		// *3\r\n :1\r\n :2\r\n:3\r\n
		curData, err := Parse(reader)
		if err != nil {
			return nil, err
		}

		result[0] = curData
		io.ReadFull(reader, crlfGarb)
	}

	return &Arr{
		Data: result,
	}, nil

}

func getLenHelper(line string) (int, error) {
	lengthStr := strings.TrimSuffix(line, "\r\n")[1:]
	strLen, err := strconv.Atoi(lengthStr)
	if err != nil {
		return 0, err
	}
	return strLen, nil
}

func ParseSimple(reader *bufio.Reader) (*Request, error) {
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
