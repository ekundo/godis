package server

import (
	"time"
)

func init() {
	registerCmd(func() cmd { return &ttlCmd{} })
}

type ttlCmd struct {
	pTtlCmd
}

func (cmd *ttlCmd) cmdName() string {
	return "ttl"
}

func (cmd *ttlCmd) exec(shard *shard) (cmdResult, error) {
	return cmd.execTtl(shard, func(t *time.Time) int64 {
		return (t.UnixNano() - time.Now().UnixNano()) / int64(time.Second)
	})
}

var _ cmd = (*ttlCmd)(nil)
