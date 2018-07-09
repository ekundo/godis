// +build integration

package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLSet(t *testing.T) {
	key := key(t, "foo")
	err := c.LSet(key, 0, "bar")
	assert.EqualError(t, err, indexOutOfRangeErrorMsg)

	_, err = c.RPush(key, "foo")
	assert.NoError(t, err)

	err = c.LSet(key, 0, "bar")
	assert.NoError(t, err)

	v, err := c.LIndex(key, 0)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "bar", *v)

	err = c.LSet(key, -1, "bak")
	assert.NoError(t, err)

	v, err = c.LIndex(key, 0)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "bak", *v)
}

func TestLSetNotList(t *testing.T) {
	key := key(t, "foo")
	_, err := c.HSet(key, "foo", "bar")
	assert.NoError(t, err)

	err = c.LSet(key, 0, "foo")
	assert.EqualError(t, err, incompatibleTypeErrorMsg)
}
