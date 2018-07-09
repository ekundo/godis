package client

import (
	"strconv"
	"time"
)

/*
Set sets the string value of a key.

If key already holds a value, it is overwritten, regardless of its type.
Any previous time to live associated with the key is discarded on successful operation.
*/
func (c *Client) Set(key string, value string, ttl *time.Duration) error {
	_, err := c.set(key, value, ttl, false, false)
	return err
}

/*
SetIfNotExists sets the string value of a key if it does not already exist.

SetIfNotExists returns "true" if operation was executed correctly.
SetIfNotExists returns "false" if the operation was not performed because the condition was not met.
*/
func (c *Client) SetIfNotExists(key string, value string, ttl *time.Duration) (bool, error) {
	return c.set(key, value, ttl, true, false)
}

/*
SetIfExists sets the string value of a key if it already exists.

SetIfExists returns "true" if operation was executed correctly.
SetIfExists returns "false" if the operation was not performed because the condition was not met.
*/
func (c *Client) SetIfExists(key string, value string, ttl *time.Duration) (bool, error) {
	return c.set(key, value, ttl, false, true)
}

func (c *Client) set(key string, value string, ttl *time.Duration, notExists bool, exists bool) (bool, error) {
	args := make([]string, 0, 6)
	args = append(args, "set", key, value)
	if ttl != nil {
		args = append(args, "px", strconv.FormatInt(ttl.Nanoseconds()/int64(time.Millisecond), 10))
	}
	if notExists {
		args = append(args, "nx")
	} else if exists {
		args = append(args, "xx")
	}
	req := cmd(args)
	res, err := c.processRequest(req)
	if err != nil {
		return false, err
	}
	return c.nullOrOk(res)
}
