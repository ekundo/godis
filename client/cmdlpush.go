package client

/*
LPush prepends the value to a list.

When the value at key is not a list, an error is returned.
LPush returns the length of the list after the push operation.
*/
func (c *Client) LPush(key string, value string) (int, error) {
	req := cmd([]string{"lpush", key, value})
	res, err := c.processRequest(req)
	if err != nil {
		return 0, err
	}
	return c.integer(res)
}
