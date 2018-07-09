package server

import (
	"github.com/ekundo/godis/resp"
)

func init() {
	registerCmd(func() cmd { return &listLeftPopCmd{} })
}

type listLeftPopCmd struct {
	listCmd
}

func (cmd *listLeftPopCmd) cmdName() string {
	return "lpop"
}

func (cmd *listLeftPopCmd) getMsg() *resp.Message {
	return &resp.Message{Element: &resp.Array{Items: []resp.Data{
		&resp.BulkString{Data: []byte(cmd.cmdName())},
		&resp.BulkString{Data: []byte(cmd.key)},
	}}}
}

func (cmd *listLeftPopCmd) arity() int {
	return 2
}

func (cmd *listLeftPopCmd) readonly() bool {
	return false
}

func (cmd *listLeftPopCmd) exec(shard *shard) (cmdResult, error) {
	v, err := cmd.listWrite(shard, cmd.key, func(item listItem) (interface{}, error) {
		value, err := item.get(0)
		if err != nil {
			return nil, err
		}
		err = item.remove(0)
		if err != nil {
			return nil, err
		}
		return value, nil
	})
	if err != nil {
		switch err.(type) {
		case keyNotFoundError, indexOutOfRangeError:
			return &nullStringResult{nil}, nil
		default:
			return nil, err
		}
	}
	value := v.(string)
	return &nullStringResult{value: &value}, nil
}

func (cmd *listLeftPopCmd) applyArgs(args []resp.Data) error {
	var err error
	cmd.key, err = cmd.stringArg(args[0])
	if err != nil {
		return err
	}
	return nil
}

var _ cmd = (*listLeftPopCmd)(nil)
