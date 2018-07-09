package client

/*
HGetAll returns all the fields and values in a hash.

When the value at key is not a hash, an error is returned.
*/
func (c *Client) HGetAll(key string) (map[string]string, error) {
	req := cmd([]string{"hgetall", key})
	res, err := c.processRequest(req)
	if err != nil {
		return nil, err
	}
	return c.stringMap(res)
}
