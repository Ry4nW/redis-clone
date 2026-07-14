package command

import (
	"errors"
	"testing"

	"redis-clone/internal/resp"
)

func mustDispatch(t *testing.T, h *Handlers, req *resp.Request) string {
	t.Helper()
	result, err := h.Dispatch(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return result
}

func dispatchErr(t *testing.T, h *Handlers, req *resp.Request) error {
	t.Helper()
	_, err := h.Dispatch(req)
	return err
}

// --- PING ---

func TestHandlePing_NoArgs(t *testing.T) {
	got := mustDispatch(t, NewHandlers(), &resp.Request{Command: "PING"})
	if want := "+PONG\r\n"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestHandlePing_OneArg(t *testing.T) {
	got := mustDispatch(t, NewHandlers(), &resp.Request{Command: "PING", Args: []string{"hello"}})
	if want := "$5\r\nhello\r\n"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestHandlePing_TooManyArgs(t *testing.T) {
	err := dispatchErr(t, NewHandlers(), &resp.Request{Command: "PING", Args: []string{"a", "b"}})
	if !errors.Is(err, ErrBadArgAmt) {
		t.Fatalf("expected ErrBadArgAmt, got %v", err)
	}
}

// --- ECHO ---

func TestHandleEcho_OneArg(t *testing.T) {
	got := mustDispatch(t, NewHandlers(), &resp.Request{Command: "ECHO", Args: []string{"hello"}})
	if want := "$5\r\nhello\r\n"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestHandleEcho_NoArgs(t *testing.T) {
	err := dispatchErr(t, NewHandlers(), &resp.Request{Command: "ECHO"})
	if !errors.Is(err, ErrBadArgAmt) {
		t.Fatalf("expected ErrBadArgAmt, got %v", err)
	}
}

func TestHandleEcho_TooManyArgs(t *testing.T) {
	err := dispatchErr(t, NewHandlers(), &resp.Request{Command: "ECHO", Args: []string{"a", "b"}})
	if !errors.Is(err, ErrBadArgAmt) {
		t.Fatalf("expected ErrBadArgAmt, got %v", err)
	}
}

// --- GET / SET ---

func TestHandleSetGet_RoundTrip(t *testing.T) {
	h := NewHandlers()

	if got := mustDispatch(t, h, &resp.Request{Command: "SET", Args: []string{"foo", "bar"}}); got != "+OK\r\n" {
		t.Fatalf("SET: got %q", got)
	}
	if got := mustDispatch(t, h, &resp.Request{Command: "GET", Args: []string{"foo"}}); got != "$3\r\nbar\r\n" {
		t.Fatalf("GET: got %q", got)
	}
}

func TestHandleGet_MissingKey(t *testing.T) {
	got := mustDispatch(t, NewHandlers(), &resp.Request{Command: "GET", Args: []string{"missing"}})
	if want := "$-1\r\n"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestHandleSet_BadArgAmt(t *testing.T) {
	err := dispatchErr(t, NewHandlers(), &resp.Request{Command: "SET", Args: []string{"foo"}})
	if !errors.Is(err, ErrBadArgAmt) {
		t.Fatalf("expected ErrBadArgAmt, got %v", err)
	}
}

// --- EXISTS / DEL ---

func TestHandleExists(t *testing.T) {
	h := NewHandlers()
	mustDispatch(t, h, &resp.Request{Command: "SET", Args: []string{"foo", "bar"}})

	if got := mustDispatch(t, h, &resp.Request{Command: "EXISTS", Args: []string{"foo"}}); got != ":1\r\n" {
		t.Fatalf("got %q", got)
	}
	if got := mustDispatch(t, h, &resp.Request{Command: "EXISTS", Args: []string{"missing"}}); got != ":0\r\n" {
		t.Fatalf("got %q", got)
	}
}

func TestHandleDel(t *testing.T) {
	h := NewHandlers()
	mustDispatch(t, h, &resp.Request{Command: "SET", Args: []string{"foo", "bar"}})

	if got := mustDispatch(t, h, &resp.Request{Command: "DEL", Args: []string{"foo", "missing"}}); got != ":1\r\n" {
		t.Fatalf("got %q", got)
	}
	if got := mustDispatch(t, h, &resp.Request{Command: "EXISTS", Args: []string{"foo"}}); got != ":0\r\n" {
		t.Fatalf("expected foo to be deleted, got %q", got)
	}
}

// --- FLUSH ---

func TestHandleFlush(t *testing.T) {
	h := NewHandlers()
	mustDispatch(t, h, &resp.Request{Command: "SET", Args: []string{"foo", "bar"}})
	mustDispatch(t, h, &resp.Request{Command: "FLUSH"})

	if got := mustDispatch(t, h, &resp.Request{Command: "EXISTS", Args: []string{"foo"}}); got != ":0\r\n" {
		t.Fatalf("expected store to be empty after FLUSH, got %q", got)
	}
}

// --- Dispatch ---

func TestDispatch_UnknownCommand(t *testing.T) {
	err := dispatchErr(t, NewHandlers(), &resp.Request{Command: "FOOBAR"})
	if !errors.Is(err, ErrUnknownCommand) {
		t.Fatalf("expected ErrUnknownCommand, got %v", err)
	}
}

func TestDispatch_CommandIsCaseInsensitive(t *testing.T) {
	got := mustDispatch(t, NewHandlers(), &resp.Request{Command: "ping"})
	if want := "+PONG\r\n"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
