package server

import "github.com/ekundo/godis/resp"

func init() {
	registerCmd(func() cmd { return &commandCmd{} })
}

type commandCmd struct {
	baseCmd
}

type commandCmdResult struct {
	items []commandCmdResultItem
}

type commandCmdResultItem struct {
	name          string
	flags         []string
	arity         int
	firstKeyIndex int
	lastKeyIndex  int
	keyStep       int
}

func (res *commandCmdResult) resp() *resp.Message {
	items := make([]resp.Data, 0, len(res.items))
	for _, resItem := range res.items {
		desc := make([]resp.Data, 0, 6)
		desc = append(desc, &resp.SimpleString{Data: []byte(resItem.name)})
		desc = append(desc, &resp.Integer{Data: resItem.arity})
		flags := make([]resp.Data, 0, len(resItem.flags))
		for _, resFlag := range resItem.flags {
			flags = append(flags, &resp.SimpleString{Data: []byte(resFlag)})
		}
		desc = append(desc, &resp.Array{Items: flags})
		desc = append(desc, &resp.Integer{Data: resItem.firstKeyIndex})
		desc = append(desc, &resp.Integer{Data: resItem.lastKeyIndex})
		desc = append(desc, &resp.Integer{Data: resItem.keyStep})
		item := &resp.Array{Items: desc}
		items = append(items, item)
	}
	return &resp.Message{Element: &resp.Array{Items: items}}
}

var _ cmdResult = (*commandCmdResult)(nil)

func (cmd *commandCmd) cmdName() string {
	return "command"
}

func (cmd *commandCmd) getMsg() *resp.Message {
	return &resp.Message{Element: &resp.Array{Items: []resp.Data{
		&resp.BulkString{Data: []byte(cmd.cmdName())},
	}}}
}

func (cmd *commandCmd) arity() int {
	return 0
}

func (cmd *commandCmd) readonly() bool {
	return true
}

func (cmd *commandCmd) exec(shard *shard) (cmdResult, error) {
	items := make([]commandCmdResultItem, 0, len(cmds))
	for key, value := range cmds {
		cmd := value()
		arity := cmd.arity()
		ro := cmd.readonly()
		var flags []string
		if ro {
			flags = []string{"readonly"}
		} else {
			flags = []string{"write"}
		}
		var keyIndex int
		if arity != 0 {
			keyIndex = 1
		} else {
			keyIndex = 0
		}
		item := commandCmdResultItem{key, flags, arity, keyIndex, keyIndex, keyIndex}
		items = append(items, item)
	}
	return &commandCmdResult{items: items}, nil
}

func (cmd *commandCmd) applyArgs(args []resp.Data) error {
	cmd.key = ""
	return nil
}

var _ cmd = (*commandCmd)(nil)
