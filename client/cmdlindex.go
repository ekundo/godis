package client

import "strconv"

/*
LIndex returns an element from a list by its index.

The index is zero-based, so 0 means the first element, 1 the second element and so on.
Negative indices can be used to designate elements starting at the tail of the list.
When the value at key is not a list, an error is returned.
LIndex returns nil when index is out of range.
*/
func (c *Client) LIndex(key string, index int) (*string, error) {
	req := cmd([]string{"lindex", key, strconv.Itoa(index)})
	res, err := c.processRequest(req)
	if err != nil {
		return nil, err
	}
	return c.bulkString(res)
}
