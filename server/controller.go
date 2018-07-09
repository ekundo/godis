package server

import (
	"github.com/ekundo/godis/resp"
	"strings"
)

var cmds = map[string]func() cmd{}

type controller struct {
	cache *cache
}

type cmdFunc func() cmd

func registerCmd(fn cmdFunc) {
	cmd := fn()
	cmds[cmd.cmdName()] = fn
}

func newController(wal *wal) *controller {
	return &controller{cache: newCache(wal)}
}

func (ctrl *controller) processRequest(req *resp.Message, writeToWal bool) (*resp.Message, error) {
	respArr, ok := req.Element.(*resp.Array)
	if !ok {
		return nil, nil
	}
	fields := respArr.Items
	if len(fields) < 1 {
		return nil, nil
	}
	respStr, ok := fields[0].(resp.String)
	if !ok {
		return nil, nil
	}
	reqCmd := strings.ToLower(string(respStr.Str()))
	cmdFunc, ok := cmds[reqCmd]
	if !ok {
		return nil, unknownCommandError{reqCmd}
	}
	cmd := cmdFunc()
	cmd.setWriteToWal(writeToWal)

	if (cmd.arity() > 0 && len(fields) != cmd.arity()) || (cmd.arity() < 0 && len(fields) < -cmd.arity()) {
		return nil, wrongNumberOfArgumentsError{cmd: cmd.cmdName()}
	}

	err := cmd.applyArgs(fields[1:])
	if err != nil {
		return nil, err
	}

	var res cmdResult
	res, err = ctrl.cache.execCmd(cmd)
	if err != nil {
		return nil, err
	}
	return res.resp(), nil
}
