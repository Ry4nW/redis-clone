package command

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"redis-clone/internal/resp"
)

type Store map[string]resp.RespValue

type Handlers struct {
	mu    sync.RWMutex
	store Store
}

func NewHandlers() *Handlers {
	return &Handlers{store: Store{}}
}

// handling for expiry
func getNowMS() int64 {
	return time.Now().UnixMilli()
}

func (h *Handlers) Dispatch(req *resp.Request) (string, error) {
	switch strings.ToUpper(req.Command) {
	case "PING":
		return h.handlePing(req)
	case "ECHO":
		return h.handleEcho(req)
	case "GET":
		return h.handleGet(req)
	case "SET":
		return h.handleSet(req)
	case "EXISTS":
		return h.handleExists(req)
	case "DEL":
		return h.handleDel(req)
	case "FLUSH":
		return h.handleFlush(req)
	default:
		return "", ErrUnknownCommand
	}
}

func isValidArgLen(req *resp.Request, expectLen int) bool {
	return len(req.Args) == expectLen
}

func (h *Handlers) handlePing(req *resp.Request) (string, error) {
	if len(req.Args) == 0 {
		return "PONG\n", nil
	}

	return h.handleEcho(req)
}

func (h *Handlers) handleEcho(req *resp.Request) (string, error) {
	if !isValidArgLen(req, 1) {
		return "", ErrBadArgAmt
	}

	return req.Args[0] + "\n", nil
}

func (h *Handlers) handleGet(req *resp.Request) (string, error) {
	if !isValidArgLen(req, 1) {
		return "", ErrBadArgAmt
	}

	key := req.Args[0]

	h.mu.Lock()
	defer h.mu.Unlock()

	entry, ok := h.store[key]
	if !ok {
		return "(nil)\n", nil
	}

	if entry.ExpiresAt != 0 && entry.ExpiresAt <= getNowMS() {
		delete(h.store, key)
		return "(nil)\n", nil
	}

	return entry.String + "\n", nil
}

func (h *Handlers) handleSet(req *resp.Request) (string, error) {
	if !isValidArgLen(req, 2) {
		return "", ErrBadArgAmt
	}

	key, val := req.Args[0], req.Args[1]

	h.mu.Lock()
	defer h.mu.Unlock()

	// TODO: impl PX, EX options
	h.store[key] = resp.NewBulkString(val)

	return "OK\n", nil
}

func (h *Handlers) handleExists(req *resp.Request) (string, error) {
	if !isValidArgLen(req, 1) {
		return "", ErrBadArgAmt
	}

	key := req.Args[0]

	h.mu.Lock()
	defer h.mu.Unlock()

	entry, ok := h.store[key]
	if !ok {
		return "0\n", nil
	}

	if entry.ExpiresAt != 0 && entry.ExpiresAt <= getNowMS() {
		delete(h.store, key)
		return "0\n", nil
	}

	return "1\n", nil
}

func (h *Handlers) handleDel(req *resp.Request) (string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	deleteCount := 0
	for _, key := range req.Args {
		if _, ok := h.store[key]; ok {
			delete(h.store, key)
			deleteCount++
		}
	}
	return strconv.Itoa(deleteCount) + "\n", nil
}

func (h *Handlers) handleFlush(req *resp.Request) (string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.store = Store{}
	return "OK\n", nil
}
