// +build integration

package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHExists(t *testing.T) {
	key := key(t, "foo")

	exists, err := c.HExists(key, "foo")
	assert.NoError(t, err)
	assert.False(t, exists)

	isNew, err := c.HSet(key, "foo", "bar")
	assert.NoError(t, err)
	assert.True(t, isNew)

	exists, err = c.HExists(key, "foo")
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = c.HExists(key, "foo1")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestHExistsNotHash(t *testing.T) {
	key := key(t, "foo")
	_, err := c.RPush(key, "bar")
	assert.NoError(t, err)

	_, err = c.HExists(key, "bar")
	assert.EqualError(t, err, incompatibleTypeErrorMsg)
}
