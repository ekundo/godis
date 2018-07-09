package server

import "github.com/ekundo/godis/resp"

func init() {
	registerCmd(func() cmd { return &dictKeysCmd{} })
}

type dictKeysCmd struct {
	dictCmd
}

type dictKeysCmdResult struct {
	value string
}

func (res *dictKeysCmdResult) resp() *resp.Message {
	return &resp.Message{Element: &resp.BulkString{Data: []byte(res.value)}}
}

var _ cmdResult = (*dictKeysCmdResult)(nil)

func (cmd *dictKeysCmd) cmdName() string {
	return "hkeys"
}

func (cmd *dictKeysCmd) getMsg() *resp.Message {
	return &resp.Message{Element: &resp.Array{Items: []resp.Data{
		&resp.BulkString{Data: []byte(cmd.cmdName())},
		&resp.BulkString{Data: []byte(cmd.key)},
	}}}
}

func (cmd *dictKeysCmd) arity() int {
	return 2
}

func (cmd *dictKeysCmd) readonly() bool {
	return true
}

func (cmd *dictKeysCmd) exec(shard *shard) (cmdResult, error) {
	v, err := cmd.dictRead(shard, cmd.key, func(item dictItem) (interface{}, error) {
		return item.keys(), nil
	})
	if err != nil {
		switch err.(type) {
		case keyNotFoundError:
			return &stringsResult{[]string{}}, nil
		default:
			return nil, err
		}
		return nil, err
	}
	items := v.([]string)
	return &stringsResult{items}, nil
}

func (cmd *dictKeysCmd) applyArgs(args []resp.Data) error {
	var err error
	if cmd.key, err = cmd.stringArg(args[0]); err != nil {
		return err
	}
	return nil
}

var _ cmd = (*dictKeysCmd)(nil)
