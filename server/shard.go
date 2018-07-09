package server

import (
	"sync"
	"time"

	"github.com/timtadh/data-structures/list"
)

const sweepDelay = 5 * time.Second

type shard struct {
	sync.RWMutex
	items map[string]item
	exps  *expirationQueue
	wal   *wal
}

func newShard(wal *wal) *shard {
	shard := &shard{
		items: make(map[string]item),
		exps:  &expirationQueue{List: list.New(0)},
		wal:   wal,
	}
	shard.startSweepLoop()
	return shard
}

func (shard *shard) item(key string, itemType itemType) (item, error) {
	item, found := shard.items[key]
	if !found || item.expired() {
		return nil, keyNotFoundError{key}
	}
	qwe := item.itemType()
	if qwe != itemType {
		return nil, incompatibleTypeError{}
	}
	return item, nil
}

func (shard *shard) execCmd(cmd cmd) (cmdResult, error) {
	if cmd.readonly() {
		shard.RLock()
		defer shard.RUnlock()
	} else {
		shard.Lock()
		if cmd.getWriteToWal() {
			shard.writeToWal(cmd)
		}
		defer shard.Unlock()
	}
	return cmd.exec(shard)
}

func (shard *shard) writeToWal(cmd cmd) {
	if !cmd.readonly() {
		shard.wal.write(cmd.getMsg())
	}
}

func (shard *shard) startSweepLoop() {
	go func() {
		for {
			time.Sleep(sweepDelay)
			shard.sweep()
		}
	}()
}

func (shard *shard) sweep() {
	defer func() {
		recover()
	}()
	shard.Lock()
	defer shard.Unlock()
	for expItem, err := shard.exps.PeekItem(); err == nil && expItem.expired(); expItem, err = shard.exps.PeekItem() {
		shard.exps.RemoveItem()
		if item, found := shard.items[expItem.key]; found && item.expired() {
			delete(shard.items, expItem.key)
		}
	}
}

var _ storage = (*shard)(nil)
