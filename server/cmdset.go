package server

import (
	"github.com/ekundo/godis/resp"
	"strconv"
	"strings"
	"time"
)

func init() {
	registerCmd(func() cmd { return &setCmd{} })
}

type setCmd struct {
	baseCmd
	value     string
	expiresAt *time.Time
	notExists bool
	ifExists  bool
}

type setCmdResult struct {
	success bool
}

func (res *setCmdResult) resp() *resp.Message {
	var r cmdResult
	r = &okResult{}
	if !res.success {
		r = &nullStringResult{nil}
	}
	return r.resp()
}

var _ cmdResult = (*setCmdResult)(nil)

func (cmd *setCmd) cmdName() string {
	return "set"
}

func (cmd *setCmd) getMsg() *resp.Message {
	data := make([]resp.Data, 0, 6)
	data = append(data, &resp.BulkString{Data: []byte(cmd.cmdName())},
		&resp.BulkString{Data: []byte(cmd.key)}, &resp.BulkString{Data: []byte(cmd.value)})
	if cmd.expiresAt != nil {
		ts := cmd.expiresAt.UnixNano() / int64(time.Millisecond)
		data = append(data, &resp.BulkString{Data: []byte("tx")},
			&resp.BulkString{Data: []byte(strconv.FormatInt(ts, 10))})
	}
	if cmd.notExists {
		data = append(data, &resp.BulkString{Data: []byte("nx")})
	}
	if cmd.ifExists {
		data = append(data, &resp.BulkString{Data: []byte("xx")})
	}
	return &resp.Message{Element: &resp.Array{Items: data}}
}

func (cmd *setCmd) arity() int {
	return -3
}

func (cmd *setCmd) readonly() bool {
	return false
}

func (cmd *setCmd) exec(shard *shard) (cmdResult, error) {
	it, found := shard.items[cmd.key]
	if cmd.notExists && found && !it.expired() {
		return &setCmdResult{success: false}, nil
	}
	if cmd.ifExists && (!found || it.expired()) {
		return &setCmdResult{success: false}, nil
	}
	if cmd.expiresAt != nil {
		shard.exps.PushItem(&expirationItem{key: cmd.key, expiresAt: cmd.expiresAt})
	}
	shard.items[cmd.key] = newStringItem(cmd.value, cmd.expiresAt)
	return &setCmdResult{success: true}, nil
}

func (cmd *setCmd) applyArgs(args []resp.Data) error {
	var err error
	if cmd.key, err = cmd.stringArg(args[0]); err != nil {
		return err
	}
	if cmd.value, err = cmd.stringArg(args[1]); err != nil {
		return err
	}
	if err = cmd.applyExpiration(args[2:]); err != nil {
		return err
	}
	return nil
}

func (cmd *setCmd) applyExpiration(args []resp.Data) error {
	argCnt := len(args)
	if argCnt < 1 {
		return nil
	}
	exp, err := cmd.stringArg(args[0])
	if err != nil {
		return err
	}
	switch strings.ToLower(exp) {
	case "ex":
		if argCnt < 2 {
			return syntaxError{}
		}
		ttl, err := cmd.intArg(args[1])
		if err != nil {
			return err
		}
		if ttl < 0 {
			return invalidExpireTimeError{}
		}
		cmd.expiresAt = expiresAtNowPlusSecs(ttl)
		err = cmd.applyExists(args[2:])
	case "px":
		if argCnt < 2 {
			return syntaxError{}
		}
		ttl, err := cmd.intArg(args[1])
		if err != nil {
			return err
		}
		if ttl < 0 {
			return invalidExpireTimeError{}
		}
		cmd.expiresAt = expiresAtNowPlusMillis(ttl)
		err = cmd.applyExists(args[2:])
	case "tx":
		if argCnt < 2 {
			return syntaxError{}
		}
		ts, err := cmd.intArg(args[1])
		if err != nil {
			return err
		}
		cmd.expiresAt = expiresAtFromMillis(ts)
		err = cmd.applyExists(args[2:])
	default:
		err = cmd.applyExists(args)
	}
	return err
}

func (cmd *setCmd) applyExists(args []resp.Data) error {
	argCnt := len(args)
	if argCnt < 1 {
		return nil
	}
	if argCnt > 1 {
		return syntaxError{}
	}
	exp, err := cmd.stringArg(args[0])
	if err != nil {
		return err
	}
	switch strings.ToLower(exp) {
	case "nx":
		cmd.notExists = true
	case "xx":
		cmd.ifExists = true
	default:
		return syntaxError{}
	}
	return err
}

var _ cmd = (*setCmd)(nil)
