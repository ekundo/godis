package resp

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type protocolTest struct {
	name string
	str  string
	msg  Message
}

var serTests = []protocolTest{
	{"simple string", "+OK\r\n", Message{Element: &SimpleString{Data: []byte("OK"), parsed: true}, parsed: true}},
	{"simple string: empty", "+\r\n", Message{Element: &SimpleString{Data: []byte(""), parsed: true}, parsed: true}},
	{"integer", ":0\r\n", Message{Element: &Integer{Data: 0, parsed: true}, parsed: true}},
	{"integer: positive", ":1000\r\n", Message{Element: &Integer{Data: 1000, parsed: true}, parsed: true}},
	{"integer: negative", ":-1000\r\n", Message{Element: &Integer{Data: -1000, parsed: true}, parsed: true}},
	{"bulk string", "$6\r\nfoobar\r\n", Message{Element: &BulkString{Data: []byte("foobar"), parsed: true}, parsed: true}},
	{"bulk string: new line", "$11\r\nfoobar\r\nbar\r\n", Message{Element: &BulkString{Data: []byte("foobar\r\nbar"), parsed: true}, parsed: true}},
	{"bulk string: new line only", "$2\r\n\r\n\r\n", Message{Element: &BulkString{Data: []byte("\r\n"), parsed: true}, parsed: true}},
	{"bulk string: empty", "$0\r\n\r\n", Message{Element: &BulkString{Data: []byte(""), parsed: true}, parsed: true}},
	{"bulk string: null", "$-1\r\n", Message{Element: &BulkString{Data: nil, parsed: true}, parsed: true}},
	{"array", "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n", Message{Element: &Array{Items: []Data{&BulkString{Data: []byte("foo"), parsed: true}, &BulkString{Data: []byte("bar"), parsed: true}}, parsed: true}, parsed: true}},
	{"array: empty", "*0\r\n", Message{Element: &Array{Items: []Data{}, parsed: true}, parsed: true}},
	{"array: integers", "*3\r\n:1\r\n:2\r\n:3\r\n", Message{Element: &Array{Items: []Data{&Integer{Data: 1, parsed: true}, &Integer{Data: 2, parsed: true}, &Integer{Data: 3, parsed: true}}, parsed: true}, parsed: true}},
	{"array: simple", "*3\r\n+foo\r\n+bar\r\n+foo\r\n", Message{Element: &Array{Items: []Data{&SimpleString{Data: []byte("foo"), parsed: true}, &SimpleString{Data: []byte("bar"), parsed: true}, &SimpleString{Data: []byte("foo"), parsed: true}}, parsed: true}, parsed: true}},
	{"array: mixed", "*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$6\r\nfoobar\r\n", Message{Element: &Array{Items: []Data{&Integer{Data: 1, parsed: true}, &Integer{Data: 2, parsed: true}, &Integer{Data: 3, parsed: true}, &Integer{Data: 4, parsed: true}, &BulkString{Data: []byte("foobar"), parsed: true}}, parsed: true}, parsed: true}},
	{"array: null", "*-1\r\n", Message{Element: &Array{Items: nil, parsed: true}, parsed: true}},
	{"array: arrays", "*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+foo\r\n+bar\r\n", Message{Element: &Array{Items: []Data{&Array{Items: []Data{&Integer{Data: 1, parsed: true}, &Integer{Data: 2, parsed: true}, &Integer{Data: 3, parsed: true}}, parsed: true}, &Array{Items: []Data{&SimpleString{Data: []byte("foo"), parsed: true}, &SimpleString{Data: []byte("bar"), parsed: true}}, parsed: true}}, parsed: true}, parsed: true}},
	{"array: null elements", "*3\r\n$3\r\nfoo\r\n$-1\r\n$3\r\nbar\r\n", Message{Element: &Array{Items: []Data{&BulkString{Data: []byte("foo"), parsed: true}, &BulkString{Data: nil, parsed: true}, &BulkString{Data: []byte("bar"), parsed: true}}, parsed: true}, parsed: true}},
	{"error", "-ERR foo bar\r\n", Message{Element: &Error{Kind: []byte("ERR"), Data: []byte("foo bar"), parsed: true}, parsed: true}},
}

var parseTests = append(serTests, []protocolTest{
	{"cmd", "cmd arg1 arg2\r\n", Message{Element: &Array{Items: []Data{&BulkString{Data: []byte("cmd"), parsed: true}, &BulkString{Data: []byte("arg1"), parsed: true}, &BulkString{Data: []byte("arg2"), parsed: true}}, parsed: true}, parsed: true, inline: true}},
}...)

func TestParse(t *testing.T) {
	for _, tt := range parseTests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &Message{}
			r := NewReader(&ErrOnEachReader{Reader: *strings.NewReader(tt.str)})
			_, err := msg.Parse(r)
			for err != nil && err.Error() == "tsterr" {
				_, err = msg.Parse(r)
			}
			assert.Nil(t, err, "parse error")
			assert.Equal(t, tt.msg, *msg, "result doesn't match")
		})
	}
}

func TestSer(t *testing.T) {
	for _, tt := range serTests {
		t.Run(tt.name, func(t *testing.T) {
			res := string(tt.msg.Ser())
			assert.Equal(t, tt.str, res, "result doesn't match")
		})
	}
}

type ErrOnEachReader struct {
	strings.Reader
	err bool
}

func (r *ErrOnEachReader) Read(buf []byte) (n int, err error) {
	for i := range buf {
		r.err = !r.err
		if r.err {
			return n, errors.New("tsterr")
		} else {
			var nn int
			nn, err = r.Reader.Read(buf[i : i+1])
			n += nn
		}
	}
	return n, err
}
