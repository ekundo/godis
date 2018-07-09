// +build integration

package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHDel(t *testing.T) {
	key := key(t, "foo")
	ok, err := c.HDel(key, "foo")
	assert.NoError(t, err)
	assert.False(t, ok)

	isNew, err := c.HSet(key, "foo", "bar")
	assert.NoError(t, err)
	assert.True(t, isNew)

	ok, err = c.HDel(key, "foo1")
	assert.NoError(t, err)
	assert.False(t, ok)

	ok, err = c.HDel(key, "foo")
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestHDelNotHash(t *testing.T) {
	key := key(t, "foo")
	_, err := c.RPush(key, "bar")
	assert.NoError(t, err)

	_, err = c.HDel(key, "bar")
	assert.EqualError(t, err, incompatibleTypeErrorMsg)
}
