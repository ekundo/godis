package client

/*
LLen returns the length of a list.

When the value at key is not a list, an error is returned.
*/
func (c *Client) LLen(key string) (int, error) {
	req := cmd([]string{"llen", key})
	res, err := c.processRequest(req)
	if err != nil {
		return 0, err
	}
	return c.integer(res)
}
