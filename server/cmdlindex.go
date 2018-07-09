package server

import (
	"github.com/ekundo/godis/resp"
	"strconv"
)

func init() {
	registerCmd(func() cmd { return &listIndexCmd{} })
}

type listIndexCmd struct {
	listCmd
	index int
}

func (cmd *listIndexCmd) cmdName() string {
	return "lindex"
}

func (cmd *listIndexCmd) getMsg() *resp.Message {
	return &resp.Message{Element: &resp.Array{Items: []resp.Data{
		&resp.BulkString{Data: []byte(cmd.cmdName())},
		&resp.BulkString{Data: []byte(cmd.key)},
		&resp.BulkString{Data: []byte(strconv.Itoa(cmd.index))},
	}}}
}

func (cmd *listIndexCmd) arity() int {
	return 3
}

func (cmd *listIndexCmd) readonly() bool {
	return true
}

func (cmd *listIndexCmd) exec(shard *shard) (cmdResult, error) {
	v, err := cmd.listRead(shard, cmd.key, func(item listItem) (interface{}, error) {
		i := cmd.index
		if i < 0 {
			i = item.size() + i
		}
		return item.get(i)
	})
	if err != nil {
		switch err.(type) {
		case keyNotFoundError, indexOutOfRangeError:
			return &nullStringResult{nil}, nil
		default:
			return nil, err
		}
	}
	item := v.(string)
	return &nullStringResult{&item}, nil
}

func (cmd *listIndexCmd) applyArgs(args []resp.Data) error {
	var err error
	cmd.key, err = cmd.stringArg(args[0])
	if err != nil {
		return err
	}
	cmd.index, err = cmd.intArg(args[1])
	if err != nil {
		return err
	}
	return nil
}

var _ cmd = (*listIndexCmd)(nil)
