package server

import "github.com/ekundo/godis/resp"

func init() {
	registerCmd(func() cmd { return &dictExistsCmd{} })
}

type dictExistsCmd struct {
	dictCmd
	field string
}

type dictExistsCmdResult struct {
	value string
}

func (res *dictExistsCmdResult) resp() *resp.Message {
	return &resp.Message{Element: &resp.BulkString{Data: []byte(res.value)}}
}

var _ cmdResult = (*dictExistsCmdResult)(nil)

func (cmd *dictExistsCmd) cmdName() string {
	return "hexists"
}

func (cmd *dictExistsCmd) getMsg() *resp.Message {
	return &resp.Message{Element: &resp.Array{Items: []resp.Data{
		&resp.BulkString{Data: []byte(cmd.cmdName())},
		&resp.BulkString{Data: []byte(cmd.key)},
		&resp.BulkString{Data: []byte(cmd.field)},
	}}}
}

func (cmd *dictExistsCmd) arity() int {
	return 3
}

func (cmd *dictExistsCmd) readonly() bool {
	return true
}

func (cmd *dictExistsCmd) exec(shard *shard) (cmdResult, error) {
	_, err := cmd.dictRead(shard, cmd.key, func(item dictItem) (interface{}, error) {
		return item.get(cmd.field)
	})
	if err != nil {
		switch err.(type) {
		case keyNotFoundError, fieldNotFoundError:
			return &intResult{0}, nil
		default:
			return nil, err
		}
	}
	return &intResult{1}, nil
}

func (cmd *dictExistsCmd) applyArgs(args []resp.Data) error {
	var err error
	if cmd.key, err = cmd.stringArg(args[0]); err != nil {
		return err
	}
	if cmd.field, err = cmd.stringArg(args[1]); err != nil {
		return err
	}
	return nil
}

var _ cmd = (*dictExistsCmd)(nil)
