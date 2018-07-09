package server

import "github.com/ekundo/godis/resp"

func init() {
	registerCmd(func() cmd { return &pExpireCmd{} })
}

type pExpireCmd struct {
	pExpireAtCmd
}

func (cmd *pExpireCmd) cmdName() string {
	return "pexpire"
}

func (cmd *pExpireCmd) applyArgs(args []resp.Data) error {
	return cmd.applyExpArgs(args, func(i int) {
		cmd.expiresAt = expiresAtNowPlusMillis(i)
	})
}

var _ cmd = (*pExpireCmd)(nil)
