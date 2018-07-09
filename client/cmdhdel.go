package client

/*
HDel deletes hash field.

When the value at key is not a hash, an error is returned.
*/
func (c *Client) HDel(key string, field string) (bool, error) {
	req := cmd([]string{"hdel", key, field})
	res, err := c.processRequest(req)
	if err != nil {
		return false, err
	}
	return c.bool(res)
}
