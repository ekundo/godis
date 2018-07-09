// +build integration

package server

import (
	"github.com/ekundo/godis/client"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCommand(t *testing.T) {
	cmds, err := c.Command()
	assert.NoError(t, err)
	assert.NotNil(t, cmds)
	assert.Equal(t, 27, len(cmds))

	assert.True(t, containsCmd(cmds, client.Command{"command", 0, []string{"readonly"}, 0, 0, 0}))
	assert.True(t, containsCmd(cmds, client.Command{"keys", 2, []string{"readonly"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"get", 2, []string{"readonly"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"exists", 2, []string{"readonly"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"set", -3, []string{"write"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"del", 2, []string{"write"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"type", 2, []string{"readonly"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"expire", 3, []string{"write"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"expireat", 3, []string{"write"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"pexpire", 3, []string{"write"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"pexpireat", 3, []string{"write"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"persist", 2, []string{"write"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"ttl", 2, []string{"readonly"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"pttl", 2, []string{"readonly"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"hkeys", 2, []string{"readonly"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"hget", 3, []string{"readonly"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"hgetall", 2, []string{"readonly"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"hexists", 3, []string{"readonly"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"hset", 4, []string{"write"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"hdel", 3, []string{"write"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"llen", 2, []string{"readonly"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"lpop", 2, []string{"write"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"rpop", 2, []string{"write"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"lpush", -3, []string{"write"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"rpush", -3, []string{"write"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"lset", 4, []string{"write"}, 1, 1, 1}))
	assert.True(t, containsCmd(cmds, client.Command{"lindex", 3, []string{"readonly"}, 1, 1, 1}))
}

func containsCmd(cmds []client.Command, cmd client.Command) bool {
	for _, c := range cmds {
		if reflect.DeepEqual(c, cmd) {
			return true
		}
	}
	return false
}
