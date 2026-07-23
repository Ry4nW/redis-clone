package command

import (
	"strings"
	"sync"
	"time"

	"redis-clone/internal/resp"
)

type Store map[string]resp.RespValue

type Handlers struct {
	mu    sync.RWMutex
	store Store
	aof   *AOF
}

// no AOF persistence, e.g. for tests
func NewHandlers() *Handlers {
	return &Handlers{store: Store{}}
}

// logs mutating commands to aof
func NewHandlersWithAOF(aof *AOF) *Handlers {
	return &Handlers{store: Store{}, aof: aof}
}

// handling for expiry
func getNowMS() int64 {
	return time.Now().UnixMilli()
}

// Dispatch encodes the result (or error) to wire-ready RESP bytes, so callers
// can always write what comes back straight to the socket.
func (h *Handlers) Dispatch(req *resp.Request) (string, error) {
	result, err := h.route(req)
	if err != nil {
		return resp.Encode(resp.NewError(err.Error())), err
	}
	return resp.Encode(result), nil
}

func (h *Handlers) execute(req *resp.Request) (resp.RespValue, error) {
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
		return resp.RespValue{}, ErrUnknownCommand
	}
}

// operate aof if cmd completed
func (h *Handlers) route(req *resp.Request) (resp.RespValue, error) {
	resp, err := h.execute(req)
	if err != nil {
		return resp, err
	}

	if req.Mutates && h.aof != nil {
		if aofErr := h.aof.AppendRequest(*req); aofErr != nil {
			return resp, aofErr
		}
	}
	return resp, nil
}

func isValidArgLen(req *resp.Request, expectLen int) bool {
	return len(req.Args) == expectLen
}

func (h *Handlers) handlePing(req *resp.Request) (resp.RespValue, error) {
	if len(req.Args) == 0 {
		return resp.NewSimpleString("PONG"), nil
	}

	return h.handleEcho(req)
}

func (h *Handlers) handleEcho(req *resp.Request) (resp.RespValue, error) {
	if !isValidArgLen(req, 1) {
		return resp.RespValue{}, ErrBadArgAmt
	}

	return resp.NewBulkString(req.Args[0]), nil
}

func (h *Handlers) handleGet(req *resp.Request) (resp.RespValue, error) {
	if !isValidArgLen(req, 1) {
		return resp.RespValue{}, ErrBadArgAmt
	}

	key := req.Args[0]

	h.mu.Lock()
	defer h.mu.Unlock()

	entry, ok := h.store[key]
	if !ok {
		return resp.NewNullBulkString(), nil
	}

	if entry.ExpiresAt != 0 && entry.ExpiresAt <= getNowMS() {
		delete(h.store, key)
		return resp.NewNullBulkString(), nil
	}

	return entry, nil
}

func (h *Handlers) handleSet(req *resp.Request) (resp.RespValue, error) {
	if !isValidArgLen(req, 2) {
		return resp.RespValue{}, ErrBadArgAmt
	}

	key, val := req.Args[0], req.Args[1]

	h.mu.Lock()
	defer h.mu.Unlock()

	// TODO: impl PX, EX options
	h.store[key] = resp.NewBulkString(val)

	return resp.NewSimpleString("OK"), nil
}

func (h *Handlers) handleExists(req *resp.Request) (resp.RespValue, error) {
	if !isValidArgLen(req, 1) {
		return resp.RespValue{}, ErrBadArgAmt
	}

	key := req.Args[0]

	h.mu.Lock()
	defer h.mu.Unlock()

	entry, ok := h.store[key]
	if !ok {
		return resp.NewInteger(0), nil
	}

	if entry.ExpiresAt != 0 && entry.ExpiresAt <= getNowMS() {
		delete(h.store, key)
		return resp.NewInteger(0), nil
	}

	return resp.NewInteger(1), nil
}

func (h *Handlers) handleDel(req *resp.Request) (resp.RespValue, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	deleteCount := 0
	for _, key := range req.Args {
		if _, ok := h.store[key]; ok {
			delete(h.store, key)
			deleteCount++
		}
	}
	return resp.NewInteger(int64(deleteCount)), nil
}

func (h *Handlers) handleFlush(req *resp.Request) (resp.RespValue, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.store = Store{}
	return resp.NewSimpleString("OK"), nil
}
