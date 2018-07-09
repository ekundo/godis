package client

import "strconv"

/*
LSet sets the value of an element in a list by its index.

When the value at key is not a list, an error is returned.
Also an error is returned for out of range indexes.
*/
func (c *Client) LSet(key string, index int, value string) error {
	req := cmd([]string{"lset", key, strconv.Itoa(index), value})
	res, err := c.processRequest(req)
	if err != nil {
		return err
	}
	return c.assertOk(res)
}
