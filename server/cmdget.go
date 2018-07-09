package server

import "github.com/ekundo/godis/resp"

func init() {
	registerCmd(func() cmd { return &getCmd{} })
}

type getCmd struct {
	baseCmd
}

func (cmd *getCmd) cmdName() string {
	return "get"
}

func (cmd *getCmd) getMsg() *resp.Message {
	return &resp.Message{Element: &resp.Array{Items: []resp.Data{
		&resp.BulkString{Data: []byte(cmd.cmdName())},
		&resp.BulkString{Data: []byte(cmd.key)},
	}}}
}

func (cmd *getCmd) arity() int {
	return 2
}

func (cmd *getCmd) readonly() bool {
	return true
}

func (cmd *getCmd) exec(shard *shard) (cmdResult, error) {
	item, err := shard.item(cmd.key, typeStr)
	if err != nil {
		if _, ok := err.(keyNotFoundError); ok {
			return &nullStringResult{nil}, nil
		}
		return nil, err
	}
	value := item.(stringItem).Value()
	return &nullStringResult{&value}, nil
}

func (cmd *getCmd) applyArgs(args []resp.Data) error {
	var err error
	if cmd.key, err = cmd.stringArg(args[0]); err != nil {
		return err
	}
	return nil
}

var _ cmd = (*getCmd)(nil)
