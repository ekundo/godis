// +build integration

package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExpireAt(t *testing.T) {
	key := key(t, "foo")
	ok, err := c.ExpireAt(key, time.Now().Add(time.Minute))
	assert.NoError(t, err)
	assert.False(t, ok)
	err = c.Set(key, "bar", nil)
	assert.NoError(t, err)
	mills, err := c.Ttl(key)
	assert.NoError(t, err)
	assert.Equal(t, -1, mills)
	ok, err = c.ExpireAt(key, time.Now().Add(time.Minute))
	assert.NoError(t, err)
	assert.True(t, ok)
	mills, err = c.Ttl(key)
	assert.NoError(t, err)
	assert.True(t, mills > 0)
}
