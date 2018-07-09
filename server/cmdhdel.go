package server

import (
	"github.com/ekundo/godis/resp"
)

func init() {
	registerCmd(func() cmd { return &dictDelCmd{} })
}

type dictDelCmd struct {
	dictCmd
	field string
}

func (cmd *dictDelCmd) cmdName() string {
	return "hdel"
}

func (cmd *dictDelCmd) getMsg() *resp.Message {
	return &resp.Message{Element: &resp.Array{Items: []resp.Data{
		&resp.BulkString{Data: []byte(cmd.cmdName())},
		&resp.BulkString{Data: []byte(cmd.key)},
		&resp.BulkString{Data: []byte(cmd.field)},
	}}}
}

func (cmd *dictDelCmd) arity() int {
	return 3
}

func (cmd *dictDelCmd) readonly() bool {
	return false
}

func (cmd *dictDelCmd) exec(shard *shard) (cmdResult, error) {
	_, err := cmd.dictWrite(shard, cmd.key, func(item dictItem) (interface{}, error) {
		return nil, item.remove(cmd.field)
	})
	if err != nil {
		switch err.(type) {
		case fieldNotFoundError:
			return &intResult{0}, nil
		default:
			return nil, err
		}
	}
	return &intResult{1}, nil
}

func (cmd *dictDelCmd) applyArgs(args []resp.Data) error {
	var err error
	if cmd.key, err = cmd.stringArg(args[0]); err != nil {
		return err
	}
	if cmd.field, err = cmd.stringArg(args[1]); err != nil {
		return err
	}
	return nil
}

var _ cmd = (*dictDelCmd)(nil)
