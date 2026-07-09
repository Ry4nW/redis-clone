package resp

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"unicode"
)

func checkErr[T any](res T, err error) (T, error) {
	// zero is zero-value of type T
	var zero T
	if err != nil {
		return zero, err
	}
	return res, nil
}

func Parse(reader *bufio.Reader) (Data, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	switch line[0] {
	case '$':
		res, err := parseBulkStr(line, reader)
		return checkErr(res, err)
	case '*':
		res, err := parseArray(line, reader)
		return checkErr(res, err)
	case '-':
		res, err := parseError(line, reader)
		return checkErr(res, err)
	case ':':
		res, err := parseInt(line, reader)
		return checkErr(res, err)
	case '+':
		res, err := parseSimpleStr(line, reader)
		return checkErr(res, err)
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

func removeCharPrefixAndCRLF(str string) string {
	return strings.TrimSuffix(str, "\r\n")[1:]
}

func parseArray(line string, reader *bufio.Reader) (*Arr, error) {
	// *<number-of-elements>\r\n<element-1>...<element-n>
	lengthStr := removeCharPrefixAndCRLF(line)
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

func parseError(line string, reader *bufio.Reader) (*RESPError, error) {
	errorStr := removeCharPrefixAndCRLF(line)
	firstSpaceIdx := strings.Index(errorStr, " ")
	errorType := errorStr[:firstSpaceIdx]
	errorMsg := errorStr[firstSpaceIdx+1:]

	return &RESPError{
		Prefix: errorType,
		Message: errorMsg,
	}, nil
}

func parseInt(line string, reader *bufio.Reader) (string, error) {
	intStr := removeCharPrefixAndCRLF(line)
	for _, c := range intStr {
		if !unicode.IsDigit(rune(c)) {
			return "", ErrTypeIsNotInt
		}
	}
	return intStr, nil
}

func parseSimpleStr(line string, reader *bufio.Reader) (string, error) {
	simpleStr := removeCharPrefixAndCRLF(line)
	return simpleStr, nil
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
