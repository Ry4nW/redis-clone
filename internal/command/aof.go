package command

import (
	"bufio"
	"errors"
	"io"
	"os"
	"redis-clone/internal/resp"
)

type SyncPolicy int

const (
	Always SyncPolicy = iota
	EverySecond
	Never
)

type AOF struct {
	file   *os.File
	policy SyncPolicy
}

func New(path string, policy SyncPolicy) (*AOF, error) {
	// logging uses appending
	// owner read + write, read for everyone else (standard perms)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &AOF{
		file:   file,
		policy: policy,
	}, nil
}

func (a *AOF) Append(data []byte) error {
	if _, err := a.file.Write(data); err != nil {
		return err
	}

	if a.policy == Always {
		return a.file.Sync()
	}

	return nil
}

func (a *AOF) Close() error {
	return a.file.Close()
}

func (a *AOF) Clear() error {
	if err := a.file.Truncate(0); err != nil {
		return err
	}

	// prevent future writes at old offset
	_, err := a.file.Seek(0, 0)
	return err
}

// converts resp req to parts to be translated, then written in resp
func (a *AOF) AppendRequest(req resp.Request) error {
	cmd, args := req.Command, req.Args
	// args are arbitrary bytes
	argBulkStrs := make([]resp.RespValue, len(args)+1)
	argBulkStrs[0] = resp.NewBulkString(cmd)

	for i, arg := range args {
		argBulkStrs[i+1] = resp.NewBulkString(arg)
	}
	newRV := resp.RespValue{
		Type:  resp.Array,
		Array: argBulkStrs,
	}

	encoded := resp.Encode(newRV)

	return a.Append([]byte(encoded))
}

// replays the log into handler's store, e.g. on startup
func (a *AOF) Load(handler *Handlers) error {
	if _, err := a.file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	reader := bufio.NewReader(a.file)
	for {
		req, err := resp.ReadRequest(reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		// execute, not route -- don't re-append what we just read
		if _, err := handler.execute(req); err != nil {
			return err
		}
	}

	_, err := a.file.Seek(0, io.SeekEnd)
	return err
}

// performs aof compaction i.e. writes a SET for each kv pair,
// then replaces appendonly.aof for future appends
func (a *AOF) Rewrite(handler *Handlers) error {
	// define rewrite mechanisms

	// temp for new cmds
	f, err := os.Create("appendonly.aof.tmp")
	if err != nil {
		return err
	}

	// write to the temp file, not a.file (the old log)
	tempAOF := &AOF{file: f, policy: a.policy}

	handler.mu.RLock()
	var writeErr error
	for k, v := range handler.store {
		newReq := resp.Request{
			Command: "SET",
			Args:    []string{k, v.String},
		}
		if writeErr = tempAOF.AppendRequest(newReq); writeErr != nil {
			break
		}
	}
	handler.mu.RUnlock()

	if writeErr != nil {
		f.Close()
		return writeErr
	}

	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	f.Close()

	if err := os.Rename("appendonly.aof.tmp", "appendonly.aof"); err != nil {
		return err
	}

	a.file.Close()
	newFile, err := os.OpenFile("appendonly.aof", os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	a.file = newFile

	return nil
}
