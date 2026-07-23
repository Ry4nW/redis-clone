package resp

import "errors"

var (
	ErrEmptyRequest   = errors.New("empty request")
	ErrMalformedRESP  = errors.New("malformed RESP")
	ErrUnexpectedEOF  = errors.New("unexpected EOF")
	ErrTypeIsNotInt   = errors.New("type not int")
	ErrUnknownCommand = errors.New("unknown command")
)
