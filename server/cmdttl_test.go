// +build integration

package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTtl(t *testing.T) {
	key := key(t, "foo")
	mills, err := c.Ttl(key)
	assert.NoError(t, err)
	assert.Equal(t, -2, mills)
	err = c.Set(key, "bar", nil)
	assert.NoError(t, err)
	mills, err = c.Ttl(key)
	assert.NoError(t, err)
	assert.Equal(t, -1, mills)
	err = c.Set(key, "bar", ttl(time.Minute))
	assert.NoError(t, err)
	mills, err = c.Ttl(key)
	assert.NoError(t, err)
	assert.True(t, mills > 0)
}
