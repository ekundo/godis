package client

/*
HGet returns the value of a hash field.

When the value at key is not a hash, an error is returned.
*/
func (c *Client) HGet(key string, field string) (*string, error) {
	req := cmd([]string{"hget", key, field})
	res, err := c.processRequest(req)
	if err != nil {
		return nil, err
	}
	return c.bulkString(res)
}
