package client

/*
LPop removes and returns the first element in a list.

When the value at key is not a list, an error is returned.
LPop returns nil when key does not exist.
*/
func (c *Client) LPop(key string) (*string, error) {
	req := cmd([]string{"lpop", key})
	res, err := c.processRequest(req)
	if err != nil {
		return nil, err
	}
	return c.bulkString(res)
}
