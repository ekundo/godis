package client

/*
HExists determines if a has field exists.

When the value at key is not a hash, an error is returned.
*/
func (c *Client) HExists(key string, field string) (bool, error) {
	req := cmd([]string{"hexists", key, field})
	res, err := c.processRequest(req)
	if err != nil {
		return false, err
	}
	return c.bool(res)
}
