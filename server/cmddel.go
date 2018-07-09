package server

import "github.com/ekundo/godis/resp"

func init() {
	registerCmd(func() cmd { return &delCmd{} })
}

type delCmd struct {
	baseCmd
}

func (cmd *delCmd) cmdName() string {
	return "del"
}

func (cmd *delCmd) getMsg() *resp.Message {
	return &resp.Message{Element: &resp.Array{Items: []resp.Data{
		&resp.BulkString{Data: []byte(cmd.cmdName())},
		&resp.BulkString{Data: []byte(cmd.key)},
	}}}
}

func (cmd *delCmd) arity() int {
	return 2
}

func (cmd *delCmd) readonly() bool {
	return false
}

func (cmd *delCmd) exec(shard *shard) (cmdResult, error) {
	_, found := shard.items[cmd.key]
	if found {
		delete(shard.items, cmd.key)
		return &boolResult{true}, nil
	}
	return &boolResult{false}, nil
}

func (cmd *delCmd) applyArgs(args []resp.Data) error {
	var err error
	cmd.key, err = cmd.stringArg(args[0])
	if err != nil {
		return err
	}
	return nil
}

var _ cmd = (*delCmd)(nil)
