package client

/*
HSet sets the string value of a hash field.

HSet returns "true" if field is a new field in the hash and value was set.
HSet returns "false" if field already exists in the hash and the value was updated.
When the value at key is not a hash, an error is returned.
*/
func (c *Client) HSet(key string, field string, value string) (bool, error) {
	req := cmd([]string{"hset", key, field, value})
	res, err := c.processRequest(req)
	if err != nil {
		return false, err
	}
	return c.bool(res)
}
