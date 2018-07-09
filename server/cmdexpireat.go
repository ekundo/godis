package server

import (
	"github.com/ekundo/godis/resp"
)

func init() {
	registerCmd(func() cmd { return &expireAtCmd{} })
}

type expireAtCmd struct {
	pExpireAtCmd
}

func (cmd *expireAtCmd) cmdName() string {
	return "expireat"
}

func (cmd *expireAtCmd) applyArgs(args []resp.Data) error {
	return cmd.applyExpArgs(args, func(i int) {
		cmd.expiresAt = expiresAtFromSecs(i)
	})
}

var _ cmd = (*expireAtCmd)(nil)
