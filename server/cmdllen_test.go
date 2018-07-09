// +build integration

package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLLen(t *testing.T) {
	key := key(t, "foo")
	size, err := c.LLen(key)
	assert.NoError(t, err)
	assert.Equal(t, 0, size)

	_, err = c.RPush(key, "foo")
	assert.NoError(t, err)

	size, err = c.LLen(key)
	assert.NoError(t, err)
	assert.Equal(t, 1, size)
}

func TestLLenNotList(t *testing.T) {
	key := key(t, "foo")
	_, err := c.HSet(key, "foo", "bar")
	assert.NoError(t, err)

	_, err = c.LLen(key)
	assert.EqualError(t, err, incompatibleTypeErrorMsg)
}
