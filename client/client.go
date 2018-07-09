package client

import (
	"fmt"
	"github.com/ekundo/godis/resp"
	"net"
	"time"
)

const (
	writeTimeout          = 10 * time.Second
	readTimeout           = 10 * time.Second
	defaultConnectTimeout = 10 * time.Second
)

var noTimeout time.Time

type Client struct {
	conn   net.Conn
	reader *resp.Reader
}

func New() *Client {
	return &Client{}
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
		c.reader = nil
	}
}

func (c *Client) Connect(host string, port uint, timeout time.Duration) error {
	var err error

	if c.conn != nil {
		c.Close()
	}

	if timeout == 0 {
		timeout = defaultConnectTimeout
	}

	if c.conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout); err != nil {
		return CommunicationError{cause: err}
	}

	c.reader = resp.NewReader(c.conn)

	return nil
}

func (c *Client) processMessage(req *resp.Message) (*resp.Message, error) {
	if c.conn == nil {
		return nil, NotConnectedError{}
	}

	defer c.conn.SetWriteDeadline(noTimeout)
	defer c.conn.SetReadDeadline(noTimeout)

	c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	if _, err := c.conn.Write(req.Ser()); err != nil {
		return nil, CommunicationError{cause: err}
	}

	msg := &resp.Message{}
	for {
		c.conn.SetReadDeadline(time.Now().Add(readTimeout))
		parsed, err := msg.Parse(c.reader)
		if err != nil {
			return nil, CommunicationError{cause: err}
		}

		if parsed {
			return msg, nil
		}
	}
}

func (c *Client) processRequest(req *resp.Message) (resp.Data, error) {
	res, err := c.processMessage(req)
	if err != nil {
		return nil, err
	}
	if respErr, ok := res.Element.(*resp.Error); ok {
		return nil, Error{code: string(respErr.Kind), message: string(respErr.Data)}
	}
	return res.Element, nil
}

func (c *Client) bool(res resp.Data) (bool, error) {
	i, err := c.integer(res)
	if err != nil {
		return false, err
	}
	return i != 0, nil
}

func (c *Client) integer(res resp.Data) (int, error) {
	i, ok := res.(*resp.Integer)
	if !ok {
		return 0, UnexpectedResponseError{}
	}
	return i.Data, nil
}

func (c *Client) simpleString(res resp.Data) (string, error) {
	str, ok := res.(*resp.SimpleString)
	if !ok {
		return "", UnexpectedResponseError{}
	}
	return string(str.Str()), nil
}

func (c *Client) stringSlice(res resp.Data) ([]string, error) {
	arr, ok := res.(*resp.Array)
	if !ok {
		return nil, UnexpectedResponseError{}
	}
	result := make([]string, len(arr.Items))
	for i, item := range arr.Items {
		str, ok := item.(*resp.BulkString)
		if !ok {
			return nil, UnexpectedResponseError{}
		}
		result[i] = string(str.Str())
	}
	return result, nil
}

func (c *Client) stringMap(res resp.Data) (map[string]string, error) {
	arr, ok := res.(*resp.Array)
	if !ok {
		return nil, UnexpectedResponseError{}
	}
	items := arr.Items
	result := make(map[string]string)
	if len(items) < 2 {
		return result, nil
	}
	for ; len(items) > 1; items = items[2:] {
		key, ok := items[0].(*resp.BulkString)
		if !ok {
			return nil, UnexpectedResponseError{}
		}
		value, ok := items[1].(*resp.BulkString)
		if !ok {
			return nil, UnexpectedResponseError{}
		}
		result[string(key.Str())] = string(value.Str())
	}
	return result, nil
}

func (c *Client) bulkString(res resp.Data) (*string, error) {
	str, ok := res.(*resp.BulkString)
	if !ok {
		return nil, UnexpectedResponseError{}
	}
	if str.Str() == nil {
		return nil, nil
	}
	v := string(str.Str())
	return &v, nil
}

func (c *Client) nullOrOk(res resp.Data) (bool, error) {
	bulk, ok := res.(*resp.BulkString)
	if ok {
		if bulk.Str() != nil {
			return false, UnexpectedResponseError{}
		}
		return false, nil
	}
	if err := c.assertOk(res); err != nil {
		return false, err
	}
	return true, nil
}

func (c *Client) assertOk(res resp.Data) error {
	str, ok := res.(*resp.SimpleString)
	if !ok {
		return UnexpectedResponseError{}
	}
	if string(str.Str()) != "OK" {
		return UnexpectedResponseError{}
	}
	return nil
}

func cmd(args []string) *resp.Message {
	items := make([]resp.Data, 0, len(args))
	for _, arg := range args {
		items = append(items, &resp.BulkString{Data: []byte(arg)})
	}
	return &resp.Message{Element: &resp.Array{Items: items}}
}
