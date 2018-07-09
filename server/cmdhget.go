package server

import "github.com/ekundo/godis/resp"

func init() {
	registerCmd(func() cmd { return &dictGetCmd{} })
}

type dictGetCmd struct {
	dictCmd
	field string
}

type dictGetCmdResult struct {
	value string
}

func (res *dictGetCmdResult) resp() *resp.Message {
	return &resp.Message{Element: &resp.BulkString{Data: []byte(res.value)}}
}

var _ cmdResult = (*dictGetCmdResult)(nil)

func (cmd *dictGetCmd) cmdName() string {
	return "hget"
}

func (cmd *dictGetCmd) getMsg() *resp.Message {
	return &resp.Message{Element: &resp.Array{Items: []resp.Data{
		&resp.BulkString{Data: []byte(cmd.cmdName())},
		&resp.BulkString{Data: []byte(cmd.key)},
		&resp.BulkString{Data: []byte(cmd.field)},
	}}}
}

func (cmd *dictGetCmd) arity() int {
	return 3
}

func (cmd *dictGetCmd) readonly() bool {
	return true
}

func (cmd *dictGetCmd) exec(shard *shard) (cmdResult, error) {
	v, err := cmd.dictRead(shard, cmd.key, func(item dictItem) (interface{}, error) {
		return item.get(cmd.field)
	})
	if err != nil {
		switch err.(type) {
		case keyNotFoundError, fieldNotFoundError:
			return &nullStringResult{nil}, nil
		default:
			return nil, err
		}
	}
	item := v.(string)
	return &nullStringResult{&item}, nil
}

func (cmd *dictGetCmd) applyArgs(args []resp.Data) error {
	var err error
	if cmd.key, err = cmd.stringArg(args[0]); err != nil {
		return err
	}
	if cmd.field, err = cmd.stringArg(args[1]); err != nil {
		return err
	}
	return nil
}

var _ cmd = (*dictGetCmd)(nil)
