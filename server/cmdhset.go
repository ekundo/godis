package server

import (
	"github.com/ekundo/godis/resp"
)

func init() {
	registerCmd(func() cmd { return &dictSetCmd{} })
}

type dictSetCmd struct {
	dictCmd
	field string
	value string
}

func (cmd *dictSetCmd) cmdName() string {
	return "hset"
}

func (cmd *dictSetCmd) getMsg() *resp.Message {
	return &resp.Message{Element: &resp.Array{Items: []resp.Data{
		&resp.BulkString{Data: []byte(cmd.cmdName())},
		&resp.BulkString{Data: []byte(cmd.key)},
		&resp.BulkString{Data: []byte(cmd.field)},
		&resp.BulkString{Data: []byte(cmd.value)},
	}}}
}

func (cmd *dictSetCmd) arity() int {
	return 4
}

func (cmd *dictSetCmd) readonly() bool {
	return false
}

func (cmd *dictSetCmd) exec(shard *shard) (cmdResult, error) {
	overwrite, err := cmd.dictWrite(shard, cmd.key, func(item dictItem) (interface{}, error) {
		overwrite := item.put(cmd.field, cmd.value)
		return overwrite, nil
	})
	if err != nil {
		return nil, err
	}
	return &boolResult{!overwrite.(bool)}, nil
}

func (cmd *dictSetCmd) applyArgs(args []resp.Data) error {
	var err error
	if cmd.key, err = cmd.stringArg(args[0]); err != nil {
		return err
	}
	if cmd.field, err = cmd.stringArg(args[1]); err != nil {
		return err
	}
	if cmd.value, err = cmd.stringArg(args[2]); err != nil {
		return err
	}
	return nil
}

var _ cmd = (*dictSetCmd)(nil)
