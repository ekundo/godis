package client

/*
RPop removes and returns the last element in a list.

When the value at key is not a list, an error is returned.
RPop returns nil when key does not exist.
*/
func (c *Client) RPop(key string) (*string, error) {
	req := cmd([]string{"rpop", key})
	res, err := c.processRequest(req)
	if err != nil {
		return nil, err
	}
	return c.bulkString(res)
}
