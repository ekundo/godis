package client

/*
HKeys returns all the fields in a hash.

When the value at key is not a hash, an error is returned.
*/
func (c *Client) HKeys(key string) ([]string, error) {
	req := cmd([]string{"hkeys", key})
	res, err := c.processRequest(req)
	if err != nil {
		return nil, err
	}
	return c.stringSlice(res)
}
