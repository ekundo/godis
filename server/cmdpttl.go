package server

import (
	"github.com/ekundo/godis/resp"
	"time"
)

func init() {
	registerCmd(func() cmd { return &pTtlCmd{} })
}

type pTtlCmd struct {
	baseCmd
}

func (cmd *pTtlCmd) cmdName() string {
	return "pttl"
}

func (cmd *pTtlCmd) getMsg() *resp.Message {
	return &resp.Message{Element: &resp.Array{Items: []resp.Data{
		&resp.BulkString{Data: []byte(cmd.cmdName())},
		&resp.BulkString{Data: []byte(cmd.key)},
	}}}
}

func (cmd *pTtlCmd) arity() int {
	return 2
}

func (cmd *pTtlCmd) readonly() bool {
	return true
}

func (cmd *pTtlCmd) exec(shard *shard) (cmdResult, error) {
	return cmd.execTtl(shard, func(t *time.Time) int64 {
		return (t.UnixNano() - time.Now().UnixNano()) / int64(time.Millisecond)
	})
}

func (cmd *pTtlCmd) execTtl(shard *shard, ttlFunc func(*time.Time) int64) (cmdResult, error) {
	it, found := shard.items[cmd.key]
	if !found || it.expired() {
		return &intResult{-2}, nil
	}
	expiresAt := it.getExpiresAt()
	if expiresAt == nil {
		return &intResult{-1}, nil
	}
	count := ttlFunc(expiresAt)
	return &intResult{int(count)}, nil
}

func (cmd *pTtlCmd) applyArgs(args []resp.Data) error {
	var err error
	if cmd.key, err = cmd.stringArg(args[0]); err != nil {
		return err
	}
	return nil
}

var _ cmd = (*pTtlCmd)(nil)
