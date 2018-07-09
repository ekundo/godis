// +build integration

package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLIndex(t *testing.T) {
	key := key(t, "foo")
	v, err := c.LIndex(key, 0)
	assert.NoError(t, err)
	assert.Nil(t, v)

	_, err = c.RPush(key, "foo")
	assert.NoError(t, err)

	v, err = c.LIndex(key, 0)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "foo", *v)

	v, err = c.LIndex(key, -1)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "foo", *v)
}

func TestLIndexNotList(t *testing.T) {
	key := key(t, "foo")
	_, err := c.HSet(key, "foo", "bar")
	assert.NoError(t, err)

	_, err = c.LIndex(key, 0)
	assert.EqualError(t, err, incompatibleTypeErrorMsg)
}
