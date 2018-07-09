package server

import "github.com/ekundo/godis/resp"

func init() {
	registerCmd(func() cmd { return &listLengthCmd{listCmd: listCmd{}} })
}

type listLengthCmd struct {
	listCmd
}

func (cmd *listLengthCmd) cmdName() string {
	return "llen"
}

func (cmd *listLengthCmd) getMsg() *resp.Message {
	return &resp.Message{Element: &resp.Array{Items: []resp.Data{
		&resp.BulkString{Data: []byte(cmd.cmdName())},
		&resp.BulkString{Data: []byte(cmd.key)},
	}}}
}

func (cmd *listLengthCmd) arity() int {
	return 2
}

func (cmd *listLengthCmd) readonly() bool {
	return true
}

func (cmd *listLengthCmd) exec(shard *shard) (cmdResult, error) {
	size, err := cmd.listRead(shard, cmd.key, func(item listItem) (interface{}, error) {
		return item.size(), nil
	})
	if err != nil {
		switch err.(type) {
		case keyNotFoundError:
			return &intResult{0}, nil
		default:
			return nil, err
		}
	}
	return &intResult{size.(int)}, nil
}

func (cmd *listLengthCmd) applyArgs(args []resp.Data) error {
	var err error
	cmd.key, err = cmd.stringArg(args[0])
	if err != nil {
		return err
	}
	return nil
}

var _ cmd = (*listLengthCmd)(nil)
