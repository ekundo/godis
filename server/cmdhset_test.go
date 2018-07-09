// +build integration

package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHSet(t *testing.T) {
	key := key(t, "foo")
	isNew, err := c.HSet(key, "foo", "bar")
	assert.NoError(t, err)
	assert.True(t, isNew)

	v, err := c.HGet(key, "foo")
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "bar", *v)

	isNew, err = c.HSet(key, "foo1", "bar1")
	assert.NoError(t, err)
	assert.True(t, isNew)

	v, err = c.HGet(key, "foo1")
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "bar1", *v)

	isNew, err = c.HSet(key, "foo", "bar2")
	assert.NoError(t, err)
	assert.False(t, isNew)

	v, err = c.HGet(key, "foo")
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "bar2", *v)
}

func TestHSetNotHash(t *testing.T) {
	key := key(t, "foo")
	_, err := c.RPush(key, "bar")
	assert.NoError(t, err)

	_, err = c.HSet(key, "foo", "bar")
	assert.EqualError(t, err, incompatibleTypeErrorMsg)
}
