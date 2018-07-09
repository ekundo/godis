package server

import (
	"time"
)

type itemType int

const (
	typeUnknown itemType = iota
	typeStr
	typeList
	typeDict
)

func (itemType itemType) String() string {
	names := [...]string{
		"unknown",
		"string",
		"list",
		"hash",
	}
	if itemType < typeStr || itemType > typeDict {
		return "unknown"
	}
	return names[itemType]
}

type item interface {
	expired() bool
	itemType() itemType
	getExpiresAt() *time.Time
	setExpiresAt(*time.Time)
}

type baseItem struct {
	expiresAt *time.Time
}

func (item *baseItem) getExpiresAt() *time.Time {
	return item.expiresAt
}

func (item *baseItem) expired() bool {
	if item.expiresAt == nil {
		return false
	}
	return item.expiresAt.Before(time.Now())
}

func (item *baseItem) setExpiresAt(expiresAt *time.Time) {
	item.expiresAt = expiresAt
}

func (item *baseItem) itemType() itemType {
	return typeUnknown
}

var _ item = (*baseItem)(nil)

func expiresAtNowPlusSecs(secs int) *time.Time {
	expiresAt := time.Now().Add(time.Duration(secs) * time.Second)
	return &expiresAt
}

func expiresAtNowPlusMillis(millis int) *time.Time {
	expiresAt := time.Now().Add(time.Duration(millis) * time.Millisecond)
	return &expiresAt
}

func expiresAtFromSecs(secs int) *time.Time {
	expiresAt := time.Unix(0, int64(secs)*int64(time.Second))
	return &expiresAt
}

func expiresAtFromMillis(millis int) *time.Time {
	expiresAt := time.Unix(0, int64(millis)*int64(time.Millisecond))
	return &expiresAt
}
