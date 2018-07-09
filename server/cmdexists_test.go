// +build integration

package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExists(t *testing.T) {
	key := key(t, "foo")
	ok, err := c.Exists(key)
	assert.NoError(t, err)
	assert.False(t, ok)
	err = c.Set(key, "bar", nil)
	assert.NoError(t, err)
	ok, err = c.Exists(key)
	assert.NoError(t, err)
	assert.True(t, ok)
}
