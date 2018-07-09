package server

import (
	"github.com/ekundo/godis/resp"
)

func init() {
	registerCmd(func() cmd { return &listRightPopCmd{} })
}

type listRightPopCmd struct {
	listCmd
}

func (cmd *listRightPopCmd) cmdName() string {
	return "rpop"
}

func (cmd *listRightPopCmd) getMsg() *resp.Message {
	return &resp.Message{Element: &resp.Array{Items: []resp.Data{
		&resp.BulkString{Data: []byte(cmd.cmdName())},
		&resp.BulkString{Data: []byte(cmd.key)},
	}}}
}

func (cmd *listRightPopCmd) arity() int {
	return 2
}

func (cmd *listRightPopCmd) readonly() bool {
	return false
}

func (cmd *listRightPopCmd) exec(shard *shard) (cmdResult, error) {
	v, err := cmd.listWrite(shard, cmd.key, func(item listItem) (interface{}, error) {
		i := item.size() - 1
		value, err := item.get(i)
		if err != nil {
			return nil, err
		}
		err = item.remove(i)
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

func (cmd *listRightPopCmd) applyArgs(args []resp.Data) error {
	var err error
	cmd.key, err = cmd.stringArg(args[0])
	if err != nil {
		return err
	}
	return nil
}

var _ cmd = (*listRightPopCmd)(nil)
