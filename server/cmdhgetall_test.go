// +build integration

package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHGetAll(t *testing.T) {
	key := key(t, "foo")

	items, err := c.HGetAll(key)
	assert.NoError(t, err)
	assert.Len(t, items, 0)

	_, err = c.HSet(key, "foo", "bar")
	assert.NoError(t, err)

	items, err = c.HGetAll(key)
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.True(t, containsEntry(items, "foo", "bar"))

	_, err = c.HSet(key, "qwe", "qwe")
	assert.NoError(t, err)

	_, err = c.HSet(key, "asd", "asd")
	assert.NoError(t, err)

	_, err = c.HSet(key, "zxc", "zxc")
	assert.NoError(t, err)

	items, err = c.HGetAll(key)
	assert.NoError(t, err)
	assert.Len(t, items, 4)
	assert.True(t, containsEntry(items, "foo", "bar"))
	assert.True(t, containsEntry(items, "qwe", "qwe"))
	assert.True(t, containsEntry(items, "asd", "asd"))
	assert.True(t, containsEntry(items, "zxc", "zxc"))
}

func TestHGetAllNotHash(t *testing.T) {
	key := key(t, "foo")
	_, err := c.RPush(key, "bar")
	assert.NoError(t, err)

	_, err = c.HGetAll(key)
	assert.EqualError(t, err, incompatibleTypeErrorMsg)
}

func containsEntry(items map[string]string, key, value string) bool {
	v, found := items[key]
	if !found {
		return false
	}
	if v != value {
		return false
	}
	return true
}
