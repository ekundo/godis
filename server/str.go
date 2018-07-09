package server

import "time"

type stringItem interface {
	item
	Value() string
	setValue(string)
}

type baseStringItem struct {
	baseItem
	value string
}

func newStringItem(value string, expiresAt *time.Time) stringItem {
	return &baseStringItem{baseItem: baseItem{expiresAt: expiresAt}, value: value}
}

func (stringItem baseStringItem) Value() string {
	return stringItem.value
}

func (stringItem *baseStringItem) setValue(value string) {
	stringItem.value = value
}

func (stringItem *baseStringItem) itemType() itemType {
	return typeStr
}

var _ item = (*baseStringItem)(nil)
var _ stringItem = (*baseStringItem)(nil)
