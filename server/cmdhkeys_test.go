// +build integration

package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHKeys(t *testing.T) {
	key := key(t, "foo")

	keys, err := c.HKeys(key)
	assert.NoError(t, err)
	assert.Len(t, keys, 0)

	_, err = c.HSet(key, "foo", "bar")
	assert.NoError(t, err)

	keys, err = c.HKeys(key)
	assert.NoError(t, err)
	assert.Len(t, keys, 1)
	assert.True(t, containsKey(keys, "foo"))

	_, err = c.HSet(key, "qwe", "qwe")
	assert.NoError(t, err)

	_, err = c.HSet(key, "asd", "asd")
	assert.NoError(t, err)

	_, err = c.HSet(key, "zxc", "zxc")
	assert.NoError(t, err)

	keys, err = c.HKeys(key)
	assert.NoError(t, err)
	assert.Len(t, keys, 4)
	assert.True(t, containsKey(keys, "foo"))
	assert.True(t, containsKey(keys, "qwe"))
	assert.True(t, containsKey(keys, "asd"))
	assert.True(t, containsKey(keys, "zxc"))
}

func TestHKeysNotHash(t *testing.T) {
	key := key(t, "foo")
	_, err := c.RPush(key, "bar")
	assert.NoError(t, err)

	_, err = c.HKeys(key)
	assert.EqualError(t, err, incompatibleTypeErrorMsg)
}

func containsKey(keys []string, key string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}
