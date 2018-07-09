package server

import "github.com/ekundo/godis/resp"

func init() {
	registerCmd(func() cmd { return &dictGetAllCmd{} })
}

type dictGetAllCmd struct {
	dictCmd
}

type dictGetAllCmdResult struct {
	value string
}

func (res *dictGetAllCmdResult) resp() *resp.Message {
	return &resp.Message{Element: &resp.BulkString{Data: []byte(res.value)}}
}

var _ cmdResult = (*dictGetAllCmdResult)(nil)

func (cmd *dictGetAllCmd) cmdName() string {
	return "hgetall"
}

func (cmd *dictGetAllCmd) getMsg() *resp.Message {
	return &resp.Message{Element: &resp.Array{Items: []resp.Data{
		&resp.BulkString{Data: []byte(cmd.cmdName())},
		&resp.BulkString{Data: []byte(cmd.key)},
	}}}
}

func (cmd *dictGetAllCmd) arity() int {
	return 2
}

func (cmd *dictGetAllCmd) readonly() bool {
	return true
}

func (cmd *dictGetAllCmd) exec(shard *shard) (cmdResult, error) {
	v, err := cmd.dictRead(shard, cmd.key, func(item dictItem) (interface{}, error) {
		return item.entries(), nil
	})
	if err != nil {
		switch err.(type) {
		case keyNotFoundError:
			return &stringsResult{[]string{}}, nil
		default:
			return nil, err
		}
	}
	itemsMap := v.(map[string]string)
	items := make([]string, 0, len(itemsMap)*2)
	for key, value := range itemsMap {
		items = append(items, key, value)
	}
	return &stringsResult{items}, nil
}

func (cmd *dictGetAllCmd) applyArgs(args []resp.Data) error {
	var err error
	if cmd.key, err = cmd.stringArg(args[0]); err != nil {
		return err
	}
	return nil
}

var _ cmd = (*dictGetAllCmd)(nil)
