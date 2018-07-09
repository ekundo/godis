package server

import "github.com/ekundo/godis/resp"

func init() {
	registerCmd(func() cmd { return &existsCmd{} })
}

type existsCmd struct {
	baseCmd
}

func (cmd *existsCmd) cmdName() string {
	return "exists"
}

func (cmd *existsCmd) getMsg() *resp.Message {
	return &resp.Message{Element: &resp.Array{Items: []resp.Data{
		&resp.BulkString{Data: []byte(cmd.cmdName())},
		&resp.BulkString{Data: []byte(cmd.key)},
	}}}
}

func (cmd *existsCmd) arity() int {
	return 2
}

func (cmd *existsCmd) readonly() bool {
	return true
}

func (cmd *existsCmd) exec(shard *shard) (cmdResult, error) {
	_, found := shard.items[cmd.key]
	if !found {
		return &boolResult{false}, nil
	}
	return &boolResult{true}, nil
}

func (cmd *existsCmd) applyArgs(args []resp.Data) error {
	var err error
	if cmd.key, err = cmd.stringArg(args[0]); err != nil {
		return err
	}
	return nil
}

var _ cmd = (*existsCmd)(nil)
