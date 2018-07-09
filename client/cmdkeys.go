package client

// Keys finds all keys matching the given pattern.
func (c *Client) Keys(pattern string) ([]string, error) {
	req := cmd([]string{"keys", pattern})
	res, err := c.processRequest(req)
	if err != nil {
		return nil, err
	}
	return c.stringSlice(res)
}
