// +build integration

package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRPush(t *testing.T) {
	key := key(t, "foo")
	size, err := c.RPush(key, "foo")
	assert.NoError(t, err)
	assert.Equal(t, 1, size)

	v, err := c.LIndex(key, 0)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "foo", *v)

	size, err = c.RPush(key, "bar")
	assert.NoError(t, err)
	assert.Equal(t, 2, size)

	v, err = c.LIndex(key, 0)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "foo", *v)

	v, err = c.LIndex(key, 1)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "bar", *v)
}

func TestRPushNotList(t *testing.T) {
	key := key(t, "foo")
	_, err := c.HSet(key, "foo", "bar")
	assert.NoError(t, err)

	_, err = c.RPush(key, "foo")
	assert.EqualError(t, err, incompatibleTypeErrorMsg)
}
