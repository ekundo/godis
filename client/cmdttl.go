package client

/*
Ttl returns the time to live for a key in milliseconds.

Ttl returns -2 if the key does not exist.
Ttl returns -1 if the key exists but has no associated expire.
*/
func (c *Client) Ttl(key string) (int, error) {
	req := cmd([]string{"pttl", key})
	res, err := c.processRequest(req)
	if err != nil {
		return 0, err
	}
	return c.integer(res)
}
