package server

import (
	"github.com/ekundo/godis/resp"
	"hash/fnv"
)

type cache struct {
	shards []storage
}

const shardNum = 256

func newCache(wal *wal) *cache {
	cache := &cache{
		shards: make([]storage, shardNum),
	}
	for i := 0; i < shardNum; i++ {
		cache.shards[i] = newShard(wal)
	}
	return cache
}

func (cache *cache) execCmd(cmd cmd) (cmdResult, error) {
	if cmd.distributed() {
		return cache.execCmdOnAllShards(cmd)
	}
	return cache.shard(cmd.getKey()).execCmd(cmd)
}

func (cache *cache) execCmdOnAllShards(cmd cmd) (cmdResult, error) {
	done := make(chan interface{})
	for i := 0; i < shardNum; i++ {
		go func(shard int) {
			res, err := cache.shards[shard].execCmd(cmd)
			if err != nil {
				done <- err
				return
			}
			done <- res
		}(i)
	}
	items := make([]resp.Data, 0)
	for i := 0; i < shardNum; i++ {
		ret := <-done
		switch ret.(type) {
		case error:
			return nil, ret.(error)
		case cmdResult:
			msg := ret.(cmdResult).resp()
			if arr, ok := msg.Element.(*resp.Array); ok {
				items = append(items, arr.Items...)
			}
		}
	}
	return &distributedCmdResult{items: items}, nil
}

type distributedCmdResult struct {
	items []resp.Data
}

func (res *distributedCmdResult) resp() *resp.Message {
	return &resp.Message{Element: &resp.Array{Items: res.items}}
}

func (cache *cache) shard(key string) (shard storage) {
	return cache.shards[hash(key)%shardNum]
}

func hash(key string) (hash uint64) {
	h := fnv.New64()
	h.Write([]byte(key))
	return h.Sum64()
}

var _ storage = (*cache)(nil)
