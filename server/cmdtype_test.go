// +build integration

package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestType(t *testing.T) {
	key := key(t, "foo")

	kind, err := c.Type(key)
	assert.NoError(t, err)
	assert.Equal(t, "none", kind)

	err = c.Set(key, "foo", nil)
	assert.NoError(t, err)

	kind, err = c.Type(key)
	assert.NoError(t, err)
	assert.Equal(t, "string", kind)

	_, err = c.Del(key)
	assert.NoError(t, err)

	_, err = c.HSet(key, "foo", "bar")
	assert.NoError(t, err)

	kind, err = c.Type(key)
	assert.NoError(t, err)
	assert.Equal(t, "hash", kind)

	_, err = c.Del(key)
	assert.NoError(t, err)

	_, err = c.RPush(key, "foo")
	assert.NoError(t, err)

	kind, err = c.Type(key)
	assert.NoError(t, err)
	assert.Equal(t, "list", kind)
}
