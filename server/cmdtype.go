package server

import "github.com/ekundo/godis/resp"

func init() {
	registerCmd(func() cmd { return &typeCmd{} })
}

type typeCmd struct {
	baseCmd
}

func (cmd *typeCmd) cmdName() string {
	return "type"
}

func (cmd *typeCmd) getMsg() *resp.Message {
	return &resp.Message{Element: &resp.Array{Items: []resp.Data{
		&resp.BulkString{Data: []byte(cmd.cmdName())},
		&resp.BulkString{Data: []byte(cmd.key)},
	}}}
}

func (cmd *typeCmd) arity() int {
	return 2
}

func (cmd *typeCmd) readonly() bool {
	return true
}

func (cmd *typeCmd) exec(shard *shard) (cmdResult, error) {
	v, found := shard.items[cmd.key]
	if !found {
		return &stringResult{"none"}, nil
	}
	return &stringResult{v.itemType().String()}, nil
}

func (cmd *typeCmd) applyArgs(args []resp.Data) error {
	var err error
	if cmd.key, err = cmd.stringArg(args[0]); err != nil {
		return err
	}
	return nil
}

var _ cmd = (*typeCmd)(nil)
