// +build integration

package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeys(t *testing.T) {
	keys, err := c.Keys(key(t, "*"))
	assert.NoError(t, err)
	assert.Len(t, keys, 0)

	addKey(t, "hello")
	addKey(t, "hallo")
	addKey(t, "hxllo")
	addKey(t, "hllo")
	addKey(t, "heeeello")
	addKey(t, "hillo")
	addKey(t, "hbllo")
	addKey(t, "hella")

	keys, err = c.Keys(key(t, "h?llo"))
	assert.NoError(t, err)
	assert.Len(t, keys, 5)
	assert.True(t, containsKey(keys, key(t, "hello")))
	assert.True(t, containsKey(keys, key(t, "hallo")))
	assert.True(t, containsKey(keys, key(t, "hxllo")))
	assert.True(t, containsKey(keys, key(t, "hillo")))
	assert.True(t, containsKey(keys, key(t, "hbllo")))

	keys, err = c.Keys(key(t, "h*llo"))
	assert.NoError(t, err)
	assert.Len(t, keys, 7)
	assert.True(t, containsKey(keys, key(t, "hello")))
	assert.True(t, containsKey(keys, key(t, "hallo")))
	assert.True(t, containsKey(keys, key(t, "hxllo")))
	assert.True(t, containsKey(keys, key(t, "hllo")))
	assert.True(t, containsKey(keys, key(t, "heeeello")))
	assert.True(t, containsKey(keys, key(t, "hillo")))
	assert.True(t, containsKey(keys, key(t, "hbllo")))

	keys, err = c.Keys(key(t, "h[ae]llo"))
	assert.NoError(t, err)
	assert.Len(t, keys, 2)
	assert.True(t, containsKey(keys, key(t, "hello")))
	assert.True(t, containsKey(keys, key(t, "hallo")))

	keys, err = c.Keys(key(t, "h[^e]llo"))
	assert.NoError(t, err)
	assert.Len(t, keys, 4)
	assert.True(t, containsKey(keys, key(t, "hallo")))
	assert.True(t, containsKey(keys, key(t, "hbllo")))
	assert.True(t, containsKey(keys, key(t, "hxllo")))
	assert.True(t, containsKey(keys, key(t, "hillo")))

	keys, err = c.Keys(key(t, "h[a-b]llo"))
	assert.NoError(t, err)
	assert.Len(t, keys, 2)
	assert.True(t, containsKey(keys, key(t, "hallo")))
	assert.True(t, containsKey(keys, key(t, "hbllo")))

}

func addKey(t *testing.T, k string) {
	err := c.Set(key(t, k), "", nil)
	assert.NoError(t, err)
}
