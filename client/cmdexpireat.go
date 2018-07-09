package client

import (
	"strconv"
	"time"
)

// ExpireAt sets the expiration for a key as a timestamp.
func (c *Client) ExpireAt(key string, t time.Time) (bool, error) {
	req := cmd([]string{"pexpireat", key, strconv.FormatInt(t.UnixNano()/int64(time.Millisecond), 10)})
	res, err := c.processRequest(req)
	if err != nil {
		return false, err
	}
	return c.bool(res)
}
