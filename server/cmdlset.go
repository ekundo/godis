package server

import (
	"github.com/ekundo/godis/resp"
	"strconv"
)

func init() {
	registerCmd(func() cmd { return &listSetCmd{} })
}

type listSetCmd struct {
	listCmd
	index int
	value string
}

func (cmd *listSetCmd) cmdName() string {
	return "lset"
}

func (cmd *listSetCmd) getMsg() *resp.Message {
	return &resp.Message{Element: &resp.Array{Items: []resp.Data{
		&resp.BulkString{Data: []byte(cmd.cmdName())},
		&resp.BulkString{Data: []byte(cmd.key)},
		&resp.BulkString{Data: []byte(strconv.Itoa(cmd.index))},
		&resp.BulkString{Data: []byte(cmd.value)},
	}}}
}

func (cmd *listSetCmd) arity() int {
	return 4
}

func (cmd *listSetCmd) readonly() bool {
	return false
}

func (cmd *listSetCmd) exec(shard *shard) (cmdResult, error) {
	_, err := cmd.listWrite(shard, cmd.key, func(item listItem) (interface{}, error) {
		i := cmd.index
		if i < 0 {
			i = item.size() + i
		}
		return nil, item.set(i, cmd.value)
	})
	if err != nil {
		return nil, err
	}
	return &okResult{}, nil
}

func (cmd *listSetCmd) applyArgs(args []resp.Data) error {
	var err error
	cmd.key, err = cmd.stringArg(args[0])
	if err != nil {
		return err
	}
	cmd.index, err = cmd.intArg(args[1])
	if err != nil {
		return err
	}
	cmd.value, err = cmd.stringArg(args[2])
	if err != nil {
		return err
	}
	return nil
}

var _ cmd = (*listSetCmd)(nil)
