package client

import (
	"strconv"
	"time"
)

/*
Expire sets a key's time to live.

Expire returns "true" if the timeout was set.
Expire returns "false" if key does not exist.
*/
func (c *Client) Expire(key string, ttl time.Duration) (bool, error) {
	req := cmd([]string{"pexpire", key, strconv.FormatInt(ttl.Nanoseconds()/int64(time.Millisecond), 10)})
	res, err := c.processRequest(req)
	if err != nil {
		return false, err
	}
	return c.bool(res)
}
