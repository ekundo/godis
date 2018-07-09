package client

/*
Type determines the type stored at a key.

Type returns type of key, or none when key does not exist.
*/
func (c *Client) Type(key string) (string, error) {
	req := cmd([]string{"type", key})
	res, err := c.processRequest(req)
	if err != nil {
		return "", err
	}
	return c.simpleString(res)
}
