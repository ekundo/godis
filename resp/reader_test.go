package resp

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestReadInt(t *testing.T) {

	var tests = []struct {
		name string
		str  string
		res  int32
		err  error
	}{
		{"simple int", "123\r\n", 123, nil},
		{"negative int", "-123\r\n", -123, nil},
		{"ending", "-123\n ", 0, malformedRespMessageError{}},
		{"overflow", "2147483648\r\n ", 0, malformedRespMessageError{}},
		{"negative overflow", "-2147483649\r\n ", 0, malformedRespMessageError{}},
		{"float", "123.2\r\n", 0, malformedRespMessageError{}},
		{"not int", "abc\r\n", 0, malformedRespMessageError{}},
		{"continue after error", "12X34X56\r\n", 123456, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(&ErrOnXReader{Reader: *strings.NewReader(tt.str)})
			i, err := r.readInt()
			for err != nil && err.Error() == "tsterr" {
				i, err = r.readInt()
			}
			if err != nil {
				assert.Equal(t, tt.err, err, "error doesn't match")
			}
			assert.Equal(t, tt.res, i, "result doesn't match")
		})
	}
}

func TestReadLine(t *testing.T) {

	var tests = []struct {
		name string
		str  string
		res  string
		err  error
	}{
		{"simple", "asd asd asd\r\nasd", "asd asd asd", nil},
		{"only first", "asd\r\n asd asd\r\nasd", "asd", nil},
		{"new line", "asd asd asd\nasd", "", io.EOF},
		{"large string", strings.Repeat(" ", 1024), "", malformedRespMessageError{}},
		{"continue after error", "asd asdX asdX\r\n", "asd asd asd", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(&ErrOnXReader{Reader: *strings.NewReader(tt.str)})
			str, err := r.readLine()
			for err != nil && err.Error() == "tsterr" {
				str, err = r.readLine()
			}
			if err != nil {
				assert.Equal(t, tt.err, err, "error doesn't match")
			} else {
				assert.Equal(t, []byte(tt.res), str, "result doesn't match")
			}
		})
	}
}

type ErrOnXReader struct {
	strings.Reader
}

func (r *ErrOnXReader) Read(buf []byte) (n int, err error) {
	for i := range buf {
		var nn int
		nn, err = r.Reader.Read(buf[i : i+1])
		if buf[i] == 'X' {
			err = errors.New("tsterr")
			return n, err
		}
		n += nn
	}
	return n, err
}
