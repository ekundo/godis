package client

// Exists determines if a key exists.
func (c *Client) Exists(key string) (bool, error) {
	req := cmd([]string{"exists", key})
	res, err := c.processRequest(req)
	if err != nil {
		return false, err
	}
	return c.bool(res)
}
