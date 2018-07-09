// +build integration

package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLPop(t *testing.T) {
	key := key(t, "foo")
	v, err := c.LPop(key)
	assert.NoError(t, err)
	assert.Nil(t, v)

	_, err = c.RPush(key, "foo")
	assert.NoError(t, err)

	_, err = c.RPush(key, "bar")
	assert.NoError(t, err)

	v, err = c.LPop(key)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "foo", *v)

	v, err = c.LPop(key)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "bar", *v)
}

func TestLPopNotList(t *testing.T) {
	key := key(t, "foo")
	_, err := c.HSet(key, "foo", "bar")
	assert.NoError(t, err)

	_, err = c.LPop(key)
	assert.EqualError(t, err, incompatibleTypeErrorMsg)
}
