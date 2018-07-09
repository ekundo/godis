// +build integration

package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGet(t *testing.T) {
	key := key(t, "foo")
	v, err := c.Get(key)
	assert.NoError(t, err)
	assert.Nil(t, v)
	err = c.Set(key, "bar", nil)
	assert.NoError(t, err)
	v, err = c.Get(key)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "bar", *v)
}

func TestGetNotString(t *testing.T) {
	key := key(t, "foo")
	_, err := c.RPush(key, "bar")
	assert.NoError(t, err)

	_, err = c.Get(key)
	assert.EqualError(t, err, incompatibleTypeErrorMsg)
}
