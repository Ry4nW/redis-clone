package resp

import "errors"

var (
	ErrEmptyRequest = errors.New("empty request")
    ErrMalformedRESP = errors.New("malformed RESP")
    ErrUnexpectedEOF = errors.New("unexpected EOF")
)

