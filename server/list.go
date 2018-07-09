package server

import (
	"github.com/emirpasic/gods/lists/doublylinkedlist"
)

type listItem interface {
	item
	get(int) (string, error)
	set(int, string) error
	add(string)
	remove(int) error
	insert(int, string) error
	size() int
}

type baseListItem struct {
	baseItem
	list *doublylinkedlist.List
}

func newListItem() listItem {
	return &baseListItem{baseItem: baseItem{}, list: doublylinkedlist.New()}
}

func (listItem *baseListItem) size() int {
	return listItem.list.Size()
}

func (listItem *baseListItem) get(index int) (string, error) {
	value, ok := listItem.list.Get(index)
	if !ok {
		return "", indexOutOfRangeError{}
	}
	return value.(string), nil
}

func (listItem *baseListItem) set(index int, value string) error {
	if index < 0 || index >= listItem.list.Size() {
		return indexOutOfRangeError{}
	}
	listItem.list.Remove(index)
	listItem.list.Insert(index, value)
	return nil
}

func (listItem *baseListItem) add(value string) {
	listItem.list.Add(value)
}

func (listItem *baseListItem) remove(index int) error {
	if index < 0 || index >= listItem.list.Size() {
		return indexOutOfRangeError{}
	}
	listItem.list.Remove(index)
	return nil
}

func (listItem *baseListItem) insert(index int, value string) error {
	if index < 0 || index > listItem.list.Size() {
		return indexOutOfRangeError{}
	}
	listItem.list.Insert(index, value)
	return nil
}

func (listItem *baseListItem) itemType() itemType {
	return typeList
}

var _ item = (*baseListItem)(nil)
var _ listItem = (*baseListItem)(nil)
