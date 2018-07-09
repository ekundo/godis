package server

import (
	"github.com/ekundo/godis/resp"
	"strconv"
	"time"
)

func init() {
	registerCmd(func() cmd { return &pExpireAtCmd{} })
}

type pExpireAtCmd struct {
	baseCmd
	expiresAt *time.Time
}

func (cmd *pExpireAtCmd) cmdName() string {
	return "pexpireat"
}

func (cmd *pExpireAtCmd) getMsg() *resp.Message {
	pExpireAtCmd := &pExpireAtCmd{}
	ts := cmd.expiresAt.UnixNano() / int64(time.Millisecond)
	return &resp.Message{Element: &resp.Array{Items: []resp.Data{
		&resp.BulkString{Data: []byte(pExpireAtCmd.cmdName())},
		&resp.BulkString{Data: []byte(cmd.key)},
		&resp.BulkString{Data: []byte(strconv.FormatInt(ts, 10))},
	}}}
}

func (cmd *pExpireAtCmd) arity() int {
	return 3
}

func (cmd *pExpireAtCmd) readonly() bool {
	return false
}

func (cmd *pExpireAtCmd) exec(shard *shard) (cmdResult, error) {
	it, found := shard.items[cmd.key]
	if !found || it.expired() {
		return &boolResult{false}, nil
	}
	shard.exps.PushItem(&expirationItem{key: cmd.key, expiresAt: cmd.expiresAt})
	it.setExpiresAt(cmd.expiresAt)
	return &boolResult{true}, nil
}

func (cmd *pExpireAtCmd) applyArgs(args []resp.Data) error {
	return cmd.applyExpArgs(args, func(i int) {
		cmd.expiresAt = expiresAtFromMillis(i)
	})
}

func (cmd *pExpireAtCmd) applyExpArgs(args []resp.Data, expArgFunc func(i int)) error {
	var err error
	if cmd.key, err = cmd.stringArg(args[0]); err != nil {
		return err
	}
	i, err := cmd.intArg(args[1])
	if err != nil {
		return err
	}
	expArgFunc(i)
	return nil
}

var _ cmd = (*pExpireAtCmd)(nil)
