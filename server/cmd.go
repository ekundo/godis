package server

import (
	"github.com/ekundo/godis/resp"
	"strconv"
)

type cmdResult interface {
	resp() *resp.Message
}

type baseCmd struct {
	key        string
	writeToWal bool
}

type cmd interface {
	cmdName() string
	getKey() string
	getWriteToWal() bool
	setWriteToWal(bool)
	exec(shard *shard) (cmdResult, error)
	applyArgs(args []resp.Data) error
	arity() int
	readonly() bool
	distributed() bool
	getMsg() *resp.Message
}

func (cmd *baseCmd) getKey() string {
	return cmd.key
}

func (cmd *baseCmd) getWriteToWal() bool {
	return cmd.writeToWal
}

func (cmd *baseCmd) setWriteToWal(writeToWal bool) {
	cmd.writeToWal = writeToWal
}

func (cmd *baseCmd) distributed() bool {
	return false
}

func (cmd *baseCmd) stringArg(arg resp.Data) (string, error) {
	respStr, ok := arg.(resp.String)
	if !ok {
		return "", wrongArgumentTypeError{}
	}
	str := respStr.Str()
	if str == nil {
		return "", nullArgumentError{}
	}
	return string(str), nil
}

func (cmd *baseCmd) intArg(arg resp.Data) (int, error) {
	respStr, err := cmd.stringArg(arg)
	if err != nil {
		return 0, err
	}
	respInt, err := strconv.Atoi(respStr)
	if err != nil {
		return 0, wrongArgumentTypeError{}
	}
	return respInt, nil
}

type boolResult struct {
	bool
}

func (res *boolResult) resp() *resp.Message {
	var data int
	if res.bool {
		data = 1
	}
	return &resp.Message{Element: &resp.Integer{Data: data}}
}

var _ cmdResult = (*boolResult)(nil)

type nullStringResult struct {
	value *string
}

func (res *nullStringResult) resp() *resp.Message {
	var data []byte = nil
	if res.value != nil {
		data = []byte(*res.value)
	}
	return &resp.Message{Element: &resp.BulkString{Data: data}}
}

var _ cmdResult = (*nullStringResult)(nil)

type okResult struct {
}

func (res *okResult) resp() *resp.Message {
	return &resp.Message{Element: &resp.SimpleString{Data: []byte("OK")}}
}

var _ cmdResult = (*okResult)(nil)

type stringResult struct {
	string
}

func (res *stringResult) resp() *resp.Message {
	return &resp.Message{Element: &resp.SimpleString{Data: []byte(res.string)}}
}

var _ cmdResult = (*stringResult)(nil)

type intResult struct {
	int
}

func (res *intResult) resp() *resp.Message {
	return &resp.Message{Element: &resp.Integer{Data: res.int}}
}

var _ cmdResult = (*intResult)(nil)

type stringsResult struct {
	items []string
}

func (res *stringsResult) resp() *resp.Message {
	items := make([]resp.Data, 0, len(res.items))
	for _, resItem := range res.items {
		item := &resp.BulkString{Data: []byte(resItem)}
		items = append(items, item)
	}
	return &resp.Message{Element: &resp.Array{Items: items}}
}

var _ cmdResult = (*stringsResult)(nil)
