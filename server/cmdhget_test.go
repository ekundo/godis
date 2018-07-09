// +build integration

package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHGet(t *testing.T) {
	key := key(t, "foo")

	v, err := c.HGet(key, "foo")
	assert.NoError(t, err)
	assert.Nil(t, v)

	isNew, err := c.HSet(key, "foo", "bar")
	assert.NoError(t, err)
	assert.True(t, isNew)

	v, err = c.HGet(key, "foo")
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "bar", *v)

	v, err = c.HGet(key, "foo1")
	assert.NoError(t, err)
	assert.Nil(t, v)
}

func TestHGetNotHash(t *testing.T) {
	key := key(t, "foo")
	_, err := c.RPush(key, "bar")
	assert.NoError(t, err)

	_, err = c.HGet(key, "bar")
	assert.EqualError(t, err, incompatibleTypeErrorMsg)
}
