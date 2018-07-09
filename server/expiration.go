package server

import (
	"container/heap"
	"encoding/binary"
	"hash/fnv"
	"time"

	"github.com/timtadh/data-structures/list"
	"github.com/timtadh/data-structures/types"
)

type expirationItem struct {
	key       string
	expiresAt *time.Time
}

func (ei *expirationItem) expired() bool {
	if ei.expiresAt == nil {
		return false
	}
	return ei.expiresAt.Before(time.Now())
}

var _ types.Hashable = (*expirationItem)(nil)

func (ei *expirationItem) Equals(b types.Equatable) bool {
	o, ok := b.(*expirationItem)
	if !ok {
		return false
	}
	if ei.key != o.key {
		return false
	}
	eiExpiresAt := *ei.expiresAt
	oExpiresAt := *o.expiresAt
	return eiExpiresAt.Equal(oExpiresAt)
}

func (ei *expirationItem) Less(b types.Sortable) bool {
	if o, ok := b.(*expirationItem); ok {
		eiExpiresAt := *ei.expiresAt
		oExpiresAt := *o.expiresAt
		return eiExpiresAt.After(oExpiresAt)
	}
	return false
}

func (ei *expirationItem) Hash() int {
	h := fnv.New32a()
	h.Write([]byte(string(ei.key)))

	eiExpiresAt := *ei.expiresAt
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint64(bs, uint64(eiExpiresAt.UnixNano()))
	h.Write(bs)

	return int(h.Sum32())
}

type expirationQueue struct {
	*list.List
}

func (eq *expirationQueue) Len() int {
	return eq.Size()
}

func (eq *expirationQueue) Less(i, j int) bool {
	left, _ := eq.Get(i)
	right, _ := eq.Get(j)
	return left.Less(right)
}

func (eq *expirationQueue) Swap(i, j int) {
	left, _ := eq.Get(i)
	right, _ := eq.Get(j)
	eq.Set(i, right)
	eq.Set(j, left)
}

func (eq *expirationQueue) Push(x interface{}) {
	eq.List.Push(x.(*expirationItem))
}

func (eq *expirationQueue) Pop() interface{} {
	item, _ := eq.List.Pop()
	return item
}

func (eq *expirationQueue) PushItem(item *expirationItem) {
	_ = eq.List.Push(item)
}

func (eq *expirationQueue) PeekItem() (*expirationItem, error) {
	item, err := eq.Get(eq.Size() - 1)
	if err != nil {
		return nil, err
	}
	return item.(*expirationItem), nil
}

func (eq *expirationQueue) RemoveItem() error {
	return eq.Remove(eq.Size() - 1)
}

var _ heap.Interface = (*expirationQueue)(nil)
