package server

import "github.com/ekundo/godis/resp"

func init() {
	registerCmd(func() cmd { return &listLeftPushCmd{listCmd: listCmd{}} })
}

type listLeftPushCmd struct {
	listCmd
	values []string
}

func (cmd *listLeftPushCmd) cmdName() string {
	return "lpush"
}

func (cmd *listLeftPushCmd) getMsg() *resp.Message {
	data := make([]resp.Data, 0, len(cmd.values)+2)
	data = append(data, &resp.BulkString{Data: []byte(cmd.cmdName())}, &resp.BulkString{Data: []byte(cmd.key)})
	for _, value := range cmd.values {
		data = append(data, &resp.BulkString{Data: []byte(value)})
	}
	return &resp.Message{Element: &resp.Array{Items: data}}
}

func (cmd *listLeftPushCmd) arity() int {
	return -3
}

func (cmd *listLeftPushCmd) readonly() bool {
	return false
}

func (cmd *listLeftPushCmd) exec(shard *shard) (cmdResult, error) {
	size, err := cmd.listWrite(shard, cmd.key, func(item listItem) (interface{}, error) {
		for _, value := range cmd.values {
			_ = item.insert(0, value)
		}
		return item.size(), nil
	})
	if err != nil {
		return nil, err
	}
	return &intResult{size.(int)}, nil
}

func (cmd *listLeftPushCmd) applyArgs(args []resp.Data) error {
	var err error
	cmd.key, err = cmd.stringArg(args[0])
	if err != nil {
		return err
	}
	cmd.values = make([]string, 0, len(args)-1)
	for _, arg := range args[1:] {
		var value string
		value, err = cmd.stringArg(arg)
		if err != nil {
			return err
		}
		cmd.values = append(cmd.values, value)
	}
	return nil
}

var _ cmd = (*listLeftPushCmd)(nil)
