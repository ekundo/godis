package server

import (
	"github.com/ekundo/godis/resp"
)

func init() {
	registerCmd(func() cmd { return &expireCmd{} })
}

type expireCmd struct {
	pExpireAtCmd
}

func (cmd *expireCmd) cmdName() string {
	return "expire"
}

func (cmd *expireCmd) applyArgs(args []resp.Data) error {
	return cmd.applyExpArgs(args, func(i int) {
		cmd.expiresAt = expiresAtNowPlusSecs(i)
	})
}

var _ cmd = (*expireCmd)(nil)
