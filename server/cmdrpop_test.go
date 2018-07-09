// +build integration

package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRPop(t *testing.T) {
	key := key(t, "foo")
	v, err := c.RPop(key)
	assert.NoError(t, err)
	assert.Nil(t, v)

	_, err = c.RPush(key, "foo")
	assert.NoError(t, err)

	_, err = c.RPush(key, "bar")
	assert.NoError(t, err)

	v, err = c.RPop(key)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "bar", *v)

	v, err = c.RPop(key)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "foo", *v)
}

func TestRPopNotList(t *testing.T) {
	key := key(t, "foo")
	_, err := c.HSet(key, "foo", "bar")
	assert.NoError(t, err)

	_, err = c.RPop(key)
	assert.EqualError(t, err, incompatibleTypeErrorMsg)
}
