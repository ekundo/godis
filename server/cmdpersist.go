package server

import (
	"github.com/ekundo/godis/resp"
)

func init() {
	registerCmd(func() cmd { return &persistCmd{} })
}

type persistCmd struct {
	baseCmd
}

func (cmd *persistCmd) cmdName() string {
	return "persist"
}

func (cmd *persistCmd) getMsg() *resp.Message {
	return &resp.Message{Element: &resp.Array{Items: []resp.Data{
		&resp.BulkString{Data: []byte(cmd.cmdName())},
	}}}
}

func (cmd *persistCmd) arity() int {
	return 2
}

func (cmd *persistCmd) readonly() bool {
	return false
}

func (cmd *persistCmd) exec(shard *shard) (cmdResult, error) {
	it, found := shard.items[cmd.key]
	if !found || it.expired() {
		return &boolResult{false}, nil
	}
	if it.getExpiresAt() == nil {
		return &boolResult{false}, nil
	}
	it.setExpiresAt(nil)
	return &boolResult{true}, nil
}

func (cmd *persistCmd) applyArgs(args []resp.Data) error {
	var err error
	if cmd.key, err = cmd.stringArg(args[0]); err != nil {
		return err
	}
	return nil
}

var _ cmd = (*persistCmd)(nil)
