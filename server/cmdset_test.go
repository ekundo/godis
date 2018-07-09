// +build integration

package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	key := key(t, "foo")
	err := c.Set(key, "bar", nil)
	assert.NoError(t, err)
	v, err := c.Get(key)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "bar", *v)
}

func TestSetTtl(t *testing.T) {
	key := key(t, "foo")
	err := c.Set(key, "bar", ttl(time.Minute))
	assert.NoError(t, err)
	mills, err := c.Ttl(key)
	assert.NoError(t, err)
	assert.True(t, mills > 0)
}

func TestSetIfNotExists(t *testing.T) {
	key := key(t, "foo")
	ok, err := c.SetIfNotExists(key, "bar", nil)
	assert.NoError(t, err)
	assert.True(t, ok)
	v, err := c.Get(key)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "bar", *v)
	ok, err = c.SetIfNotExists(key, "far", nil)
	assert.NoError(t, err)
	assert.False(t, ok)
	v, err = c.Get(key)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "bar", *v)
}

func TestSetIfExists(t *testing.T) {
	key := key(t, "foo")
	ok, err := c.SetIfExists(key, "bar", nil)
	assert.NoError(t, err)
	assert.False(t, ok)
	v, err := c.Get(key)
	assert.NoError(t, err)
	assert.Nil(t, v)
	err = c.Set(key, "far", nil)
	assert.NoError(t, err)
	v, err = c.Get(key)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "far", *v)
	ok, err = c.SetIfExists(key, "bar", nil)
	assert.NoError(t, err)
	assert.True(t, ok)
	v, err = c.Get(key)
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "bar", *v)
}

func TestSetNotString(t *testing.T) {
	key := key(t, "foo")
	_, err := c.HSet(key, "foo", "bar")
	assert.NoError(t, err)

	err = c.Set(key, "foo", nil)
	assert.NoError(t, err)
}
