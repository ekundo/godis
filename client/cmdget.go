package client

// Get gets the value of a key.
func (c *Client) Get(key string) (*string, error) {
	req := cmd([]string{"get", key})
	res, err := c.processRequest(req)
	if err != nil {
		return nil, err
	}
	return c.bulkString(res)
}
