package resp

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

func newReader(s string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(s))
}

// --- inline parser (ParseSimple) ---

func mustParseSimple(t *testing.T, input string) *Request {
	t.Helper()
	req, err := ParseSimple(newReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return req
}

func TestParseSimple_EmptyInput(t *testing.T) {
	req, err := ParseSimple(newReader("\n"))
	if req != nil {
		t.Fatalf("expected nil request, got %+v", req)
	}
	if err != ErrEmptyRequest {
		t.Fatalf("expected ErrEmptyRequest, got %v", err)
	}
}

func TestParseSimple_CommandNormalization(t *testing.T) {
	for _, in := range []string{"ping\n", "PING\n", "Ping\n"} {
		req := mustParseSimple(t, in)
		if req.Command != "PING" {
			t.Fatalf("input %q: expected Command %q, got %q", in, "PING", req.Command)
		}
	}
}

func TestParseSimple_SingleCommand(t *testing.T) {
	req := mustParseSimple(t, "PING\n")
	if req.Command != "PING" {
		t.Fatalf("expected Command PING, got %q", req.Command)
	}
	if len(req.Args) != 0 {
		t.Fatalf("expected no args, got %v", req.Args)
	}
}

func TestParseSimple_CommandWithOneArg(t *testing.T) {
	req := mustParseSimple(t, "ECHO hello\n")
	if req.Command != "ECHO" || len(req.Args) != 1 || req.Args[0] != "hello" {
		t.Fatalf("unexpected request: %+v", req)
	}
}

func TestParseSimple_CommandWithMultipleArgs(t *testing.T) {
	req := mustParseSimple(t, "SET foo bar\n")
	if req.Command != "SET" || len(req.Args) != 2 || req.Args[0] != "foo" || req.Args[1] != "bar" {
		t.Fatalf("unexpected request: %+v", req)
	}
}

func TestParseSimple_EOFReturnsError(t *testing.T) {
	_, err := ParseSimple(newReader(""))
	if err != io.EOF {
		t.Fatalf("expected io.EOF, got %v", err)
	}
}

// --- RESP binary protocol ---

func mustParseRESP(t *testing.T, input string) RespValue {
	t.Helper()
	v, err := Parse(newReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return v
}

func TestParseRESP_SimpleString(t *testing.T) {
	v := mustParseRESP(t, "+OK\r\n")
	if v.Type != SimpleString || v.String != "OK" {
		t.Fatalf("unexpected value: %+v", v)
	}
}

func TestParseRESP_Error(t *testing.T) {
	v := mustParseRESP(t, "-ERR unknown command\r\n")
	if v.Type != Error || v.String != "ERR unknown command" {
		t.Fatalf("unexpected value: %+v", v)
	}
}

func TestParseRESP_Integer(t *testing.T) {
	v := mustParseRESP(t, ":1000\r\n")
	if v.Type != Integer || v.Integer != 1000 {
		t.Fatalf("unexpected value: %+v", v)
	}
}

func TestParseRESP_BulkString(t *testing.T) {
	v := mustParseRESP(t, "$6\r\nfoobar\r\n")
	if v.Type != BulkString || v.Null || v.String != "foobar" {
		t.Fatalf("unexpected value: %+v", v)
	}
}

func TestParseRESP_NullBulkString(t *testing.T) {
	v := mustParseRESP(t, "$-1\r\n")
	if v.Type != BulkString || !v.Null {
		t.Fatalf("unexpected value: %+v", v)
	}
}

func TestParseRESP_Array_PingNoArgs(t *testing.T) {
	v := mustParseRESP(t, "*1\r\n$4\r\nPING\r\n")
	if v.Type != Array || len(v.Array) != 1 {
		t.Fatalf("unexpected value: %+v", v)
	}
	if v.Array[0].String != "PING" {
		t.Fatalf("unexpected first element: %+v", v.Array[0])
	}
}

func TestParseRESP_Array_EchoWithArg(t *testing.T) {
	v := mustParseRESP(t, "*2\r\n$4\r\nECHO\r\n$5\r\nhello\r\n")
	if v.Type != Array || len(v.Array) != 2 {
		t.Fatalf("unexpected value: %+v", v)
	}
	if v.Array[0].String != "ECHO" || v.Array[1].String != "hello" {
		t.Fatalf("unexpected elements: %+v", v.Array)
	}
}

func TestParseRESP_Array_SetKeyValue(t *testing.T) {
	v := mustParseRESP(t, "*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n")
	if v.Type != Array || len(v.Array) != 3 {
		t.Fatalf("unexpected value: %+v", v)
	}
	if v.Array[0].String != "SET" || v.Array[1].String != "foo" || v.Array[2].String != "bar" {
		t.Fatalf("unexpected elements: %+v", v.Array)
	}
}

func TestParseRESP_EmptyArray(t *testing.T) {
	v := mustParseRESP(t, "*0\r\n")
	if v.Type != Array || v.Null || len(v.Array) != 0 {
		t.Fatalf("unexpected value: %+v", v)
	}
}

func TestParseRESP_MixedTypesInArray(t *testing.T) {
	v := mustParseRESP(t, "*2\r\n:1\r\n+OK\r\n")
	if len(v.Array) != 2 {
		t.Fatalf("unexpected value: %+v", v)
	}
	if v.Array[0].Type != Integer || v.Array[0].Integer != 1 {
		t.Fatalf("unexpected first element: %+v", v.Array[0])
	}
	if v.Array[1].Type != SimpleString || v.Array[1].String != "OK" {
		t.Fatalf("unexpected second element: %+v", v.Array[1])
	}
}

func TestParseRESP_NestedArray(t *testing.T) {
	v := mustParseRESP(t, "*2\r\n*1\r\n$4\r\nPING\r\n$2\r\nhi\r\n")
	if len(v.Array) != 2 {
		t.Fatalf("unexpected value: %+v", v)
	}
	inner := v.Array[0]
	if inner.Type != Array || len(inner.Array) != 1 || inner.Array[0].String != "PING" {
		t.Fatalf("unexpected nested array: %+v", inner)
	}
	if v.Array[1].String != "hi" {
		t.Fatalf("unexpected second element: %+v", v.Array[1])
	}
}

func TestParseRESP_MalformedMissingCRLF(t *testing.T) {
	_, err := Parse(newReader("$6\r\nfoobar"))
	if err != ErrUnexpectedEOF {
		t.Fatalf("expected ErrUnexpectedEOF, got %v", err)
	}
}

func TestParseRESP_MalformedBadCount(t *testing.T) {
	_, err := Parse(newReader("*abc\r\n"))
	if err != ErrMalformedRESP {
		t.Fatalf("expected ErrMalformedRESP, got %v", err)
	}
}

func TestParseRESP_MalformedUnknownType(t *testing.T) {
	_, err := Parse(newReader("!nope\r\n"))
	if err != ErrMalformedRESP {
		t.Fatalf("expected ErrMalformedRESP, got %v", err)
	}
}

// --- ReadRequest (protocol auto-detection) ---

func TestReadRequest_InlineFallsBackToParseSimple(t *testing.T) {
	req, err := ReadRequest(newReader("SET foo bar\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Command != "SET" || len(req.Args) != 2 || req.Args[0] != "foo" || req.Args[1] != "bar" {
		t.Fatalf("unexpected request: %+v", req)
	}
}

func TestReadRequest_ArrayPingNoArgs(t *testing.T) {
	req, err := ReadRequest(newReader("*1\r\n$4\r\nPING\r\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Command != "PING" || len(req.Args) != 0 {
		t.Fatalf("unexpected request: %+v", req)
	}
}

func TestReadRequest_ArrayEchoWithArg(t *testing.T) {
	req, err := ReadRequest(newReader("*2\r\n$4\r\nECHO\r\n$5\r\nhello\r\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Command != "ECHO" || len(req.Args) != 1 || req.Args[0] != "hello" {
		t.Fatalf("unexpected request: %+v", req)
	}
}

func TestReadRequest_ArraySetKeyValue(t *testing.T) {
	req, err := ReadRequest(newReader("*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Command != "SET" || len(req.Args) != 2 || req.Args[0] != "foo" || req.Args[1] != "bar" {
		t.Fatalf("unexpected request: %+v", req)
	}
}

func TestReadRequest_ArrayCommandNormalization(t *testing.T) {
	req, err := ReadRequest(newReader("*1\r\n$4\r\nping\r\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Command != "PING" {
		t.Fatalf("expected normalized Command PING, got %q", req.Command)
	}
}

func TestReadRequest_EmptyArray(t *testing.T) {
	_, err := ReadRequest(newReader("*0\r\n"))
	if err != ErrEmptyRequest {
		t.Fatalf("expected ErrEmptyRequest, got %v", err)
	}
}

func TestReadRequest_NullArray(t *testing.T) {
	_, err := ReadRequest(newReader("*-1\r\n"))
	if err != ErrEmptyRequest {
		t.Fatalf("expected ErrEmptyRequest, got %v", err)
	}
}

func TestReadRequest_NonBulkStringElement(t *testing.T) {
	_, err := ReadRequest(newReader("*1\r\n:5\r\n"))
	if err != ErrMalformedRESP {
		t.Fatalf("expected ErrMalformedRESP, got %v", err)
	}
}

func TestReadRequest_EOF(t *testing.T) {
	_, err := ReadRequest(newReader(""))
	if err != io.EOF {
		t.Fatalf("expected io.EOF, got %v", err)
	}
}
