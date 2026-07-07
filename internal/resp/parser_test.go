package resp

import (
	"bufio"
	"strings"
	"testing"
)

// --- inline parser (current) ---

func TestParse_EmptyInput(t *testing.T) {
	// TODO: empty string should return ErrEmptyRequest
}

func TestParse_CommandNormalization(t *testing.T) {
	// TODO: "ping\n" and "PING\n" and "Ping\n" should all parse to Command == "PING"
}

func TestParse_SingleCommand(t *testing.T) {
	// TODO: "PING\n" → Request{Command: "PING", Args: []}
}

func TestParse_CommandWithOneArg(t *testing.T) {
	// TODO: "ECHO hello\n" → Request{Command: "ECHO", Args: ["hello"]}
}

func TestParse_CommandWithMultipleArgs(t *testing.T) {
	// TODO: "SET foo bar\n" → Request{Command: "SET", Args: ["foo", "bar"]}
}

func TestParse_EOFReturnsError(t *testing.T) {
	// TODO: reader with no data should return io.EOF
}

// helper to reduce boilerplate in tests above once implemented
func newReader(s string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(s))
}

// --- RESP binary protocol (next up) ---
// These tests will fail until the parser is rewritten to handle real RESP framing.

func TestParseRESP_SimpleString(t *testing.T) {
	// TODO: "+OK\r\n" → Command: "OK", type: simple string
}

func TestParseRESP_Error(t *testing.T) {
	// TODO: "-ERR unknown command\r\n" → error type
}

func TestParseRESP_Integer(t *testing.T) {
	// TODO: ":1000\r\n" → integer 1000
}

func TestParseRESP_BulkString(t *testing.T) {
	// TODO: "$6\r\nfoobar\r\n" → bulk string "foobar"
}

func TestParseRESP_NullBulkString(t *testing.T) {
	// TODO: "$-1\r\n" → null bulk string (nil)
}

func TestParseRESP_Array_PingNoArgs(t *testing.T) {
	// TODO: "*1\r\n$4\r\nPING\r\n" → Request{Command: "PING", Args: []}
}

func TestParseRESP_Array_EchoWithArg(t *testing.T) {
	// TODO: "*2\r\n$4\r\nECHO\r\n$5\r\nhello\r\n" → Request{Command: "ECHO", Args: ["hello"]}
}

func TestParseRESP_Array_SetKeyValue(t *testing.T) {
	// TODO: "*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n" → Request{Command: "SET", Args: ["foo", "bar"]}
}

func TestParseRESP_EmptyArray(t *testing.T) {
	// TODO: "*0\r\n" → should return ErrEmptyRequest or similar
}

func TestParseRESP_MalformedMissingCRLF(t *testing.T) {
	// TODO: "$6\r\nfoobar" (no trailing \r\n) → ErrMalformedRESP
}

func TestParseRESP_MalformedBadCount(t *testing.T) {
	// TODO: "*abc\r\n" → ErrMalformedRESP
}
