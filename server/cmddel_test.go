// +build integration

package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDel(t *testing.T) {
	key := key(t, "foo")
	err := c.Set(key, "bar", nil)
	assert.NoError(t, err)
	ok, err := c.Del(key)
	assert.NoError(t, err)
	assert.True(t, ok)
	v, err := c.Get(key)
	assert.NoError(t, err)
	assert.Nil(t, v)
	ok, err = c.Del(key)
	assert.NoError(t, err)
	assert.False(t, ok)
}
