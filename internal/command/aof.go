package command

import (
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
func (a *AOF) AppendRequest(req resp.Request) {
	cmd, args := req.Command, req.Args
	// args are arbitrary bytes
	argBulkStrs := make([]resp.RespValue, len(args)+1)
	argBulkStrs[0] = resp.NewBulkString(cmd)

	for i, arg := range args {
		argBulkStrs[i] = resp.NewBulkString(arg)
	}
	newRV := resp.RespValue{
		Type:  resp.Array,
		Array: argBulkStrs,
	}

	encoded := resp.Encode(newRV)

	a.Append([]byte(encoded))
}

// performs aof compaction i.e. writes a SET for each kv pair,
// then replaces appendonly.aof for future appends
func (a *AOF) Rewrite(handler *Handlers) error {
	// define rewrite mechanisms

	// temp for new cmds
	f, err := os.Create("appendonly.aof.temp")
	if err != nil {
		return err
	}
	defer f.Close()

	for k, v := range handler.store {
		newReq := resp.Request{
			Command: "SET",
			Args:    []string{k, v.String},
		}
		a.AppendRequest(newReq)
	}
	os.Rename("appendonly.aof.tmp", "appendonly.aof")
	return nil
}
