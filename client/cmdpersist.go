package client

/*
Persist removes the expiration from a key.

Persist returns "true" if the timeout was removed.
Persist returns "false" if key does not exist or does not have an associated timeout.
*/
func (c *Client) Persist(key string) (bool, error) {
	req := cmd([]string{"persist", key})
	res, err := c.processRequest(req)
	if err != nil {
		return false, err
	}
	return c.bool(res)
}
