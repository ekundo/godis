package client

// Del deletes a key.
func (c *Client) Del(key string) (bool, error) {
	req := cmd([]string{"del", key})
	res, err := c.processRequest(req)
	if err != nil {
		return false, err
	}
	return c.bool(res)
}
