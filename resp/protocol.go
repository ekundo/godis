package resp

import (
	"io"
	"math"
	"strconv"
	"strings"
)

const (
	respSimpleStringSymbol = '+'
	respErrorSymbol        = '-'
	respIntegerSymbol      = ':'
	respBulkStringSymbol   = '$'
	respArraySymbol        = '*'
	maxLengthBytes         = 512 * 1024 * 1024
)

var respElementFactory = map[byte]func() Data{
	respSimpleStringSymbol: func() Data { return &SimpleString{} },
	respErrorSymbol:        func() Data { return &Error{} },
	respIntegerSymbol:      func() Data { return &Integer{} },
	respBulkStringSymbol:   func() Data { return &BulkString{} },
	respArraySymbol:        func() Data { return &Array{} },
}

type Data interface {
	parse(r *Reader) (bool, error)
	serLen() int
	ser([]byte) []byte
	fullyParsed() bool
}

type String interface {
	Data
	Str() []byte
}

type SimpleString struct {
	Data   []byte
	parsed bool
}

func (str *SimpleString) Str() []byte {
	return str.Data
}

func (str *SimpleString) parse(r *Reader) (bool, error) {
	if str.parsed {
		return true, nil
	}
	bytes, err := r.readLine()
	if err != nil {
		return false, err
	}
	str.Data = append([]byte{}, bytes...)
	str.parsed = true
	return str.parsed, nil
}

func (str *SimpleString) serLen() int {
	return len(str.Data) + 3
}

func (str *SimpleString) ser(bytes []byte) []byte {
	bytes = append(bytes, respSimpleStringSymbol)
	bytes = append(bytes, str.Data...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (str *SimpleString) fullyParsed() bool {
	return str.parsed
}

type Error struct {
	Kind   []byte
	Data   []byte
	parsed bool
}

func (er *Error) parse(r *Reader) (bool, error) {
	if er.parsed {
		return true, nil
	}
	bytes, err := r.readLine()
	if err != nil {
		return false, err
	}
	splits := strings.SplitN(string(bytes), " ", 2)
	if len(splits) < 2 {
		return false, malformedRespMessageError{}
	}
	er.Kind = []byte(splits[0])
	er.Data = []byte(splits[1])
	er.parsed = true
	return er.parsed, nil
}

func (er *Error) serLen() int {
	return len(er.Kind) + len(er.Data) + 4
}

func (er *Error) ser(bytes []byte) []byte {
	bytes = append(bytes, respErrorSymbol)
	bytes = append(bytes, er.Kind...)
	bytes = append(bytes, ' ')
	bytes = append(bytes, er.Data...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (er *Error) fullyParsed() bool {
	return er.parsed
}

type Integer struct {
	Data   int
	parsed bool
}

func (integer *Integer) Str() []byte {
	return []byte(strconv.Itoa(integer.Data))
}

func (integer *Integer) parse(r *Reader) (bool, error) {
	if integer.parsed {
		return true, nil
	}
	i, err := r.readInt()
	if err != nil {
		return false, err
	}
	integer.Data = int(i)
	integer.parsed = true
	return integer.parsed, nil
}

func (integer *Integer) serLen() int {
	str := strconv.Itoa(integer.Data)
	return len(str) + 3
}

func (integer *Integer) ser(bytes []byte) []byte {
	str := strconv.Itoa(integer.Data)
	bytes = append(bytes, respIntegerSymbol)
	bytes = append(bytes, []byte(str)...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (integer *Integer) fullyParsed() bool {
	return integer.parsed
}

type BulkString struct {
	Data   []byte
	parsed bool
}

func (str *BulkString) Str() []byte {
	return str.Data
}

func (str *BulkString) parse(r *Reader) (bool, error) {
	if str.parsed {
		return str.parsed, nil
	}
	if str.Data == nil {
		var size int32
		var err error
		size, err = r.readInt()
		if err != nil {
			return false, err
		}
		if size < -1 || size > maxLengthBytes {
			return false, malformedRespMessageError{}
		}
		if size == -1 {
			str.parsed = true
			return str.parsed, nil
		}
		str.Data = make([]byte, 0, size)
		if size == 0 {
			str.parsed = true
			return str.parsed, nil
		}
	}
	size := len(str.Data)
	need := cap(str.Data) - size
	n, err := io.ReadFull(r, str.Data[size:size+need])
	str.Data = str.Data[:size+n]
	if err != nil {
		return false, err
	}
	err = r.readCRLF()
	if err != nil {
		return false, err
	}
	str.parsed = true
	return str.parsed, nil
}

func (str *BulkString) serLen() int {
	if str.Data == nil {
		return 5
	} else {
		size := len(str.Data)
		s := strconv.Itoa(size)
		return len(s) + size + 5
	}
}

func (str *BulkString) ser(bytes []byte) []byte {
	var size int
	if str.Data == nil {
		size = -1
	} else {
		size = len(str.Data)
	}
	s := strconv.Itoa(size)
	bytes = append(bytes, respBulkStringSymbol)
	bytes = append(bytes, []byte(s)...)
	bytes = append(bytes, '\r', '\n')
	if str.Data != nil {
		bytes = append(bytes, str.Data...)
		bytes = append(bytes, '\r', '\n')
	}
	return bytes
}

func (str *BulkString) fullyParsed() bool {
	return str.parsed
}

type Array struct {
	Items  []Data
	parsed bool
}

func (array *Array) parse(r *Reader) (bool, error) {
	if array.parsed {
		return array.parsed, nil
	}
	if array.Items == nil {
		size, err := r.readInt()
		if err != nil {
			return false, err
		}
		if size < -1 || size > math.MaxInt32 {
			return false, malformedRespMessageError{}
		}
		if size == -1 {
			array.parsed = true
			return array.parsed, nil
		}
		array.Items = make([]Data, 0, size)
		if size == 0 {
			array.parsed = true
			return array.parsed, nil
		}
	}
	for !array.parsed {
		size := len(array.Items)
		var item Data
		if size == 0 || array.Items[size-1].fullyParsed() {
			var err error
			item, err = parseElement(r)
			if err != nil {
				return false, err
			}
			array.Items = append(array.Items, item)
			size++
		} else {
			item = array.Items[size-1]
		}
		parsed, err := item.parse(r)
		if err != nil {
			return false, err
		}
		array.parsed = parsed && size == cap(array.Items)
	}
	return array.parsed, nil
}

func (array *Array) serLen() int {
	if array.Items == nil {
		return 5
	} else {
		size := len(array.Items)
		s := strconv.Itoa(size)
		res := len(s)
		for _, item := range array.Items {
			res += item.serLen()
		}
		res += 3
		return res
	}
}

func (array *Array) ser(bytes []byte) []byte {
	bytes = append(bytes, respArraySymbol)
	var size int
	if array.Items == nil {
		size = -1
	} else {
		size = len(array.Items)
	}
	s := strconv.Itoa(size)
	bytes = append(bytes, []byte(s)...)
	bytes = append(bytes, '\r', '\n')
	for _, item := range array.Items {
		bytes = item.ser(bytes)
	}
	return bytes
}

func (array *Array) fullyParsed() bool {
	return array.parsed
}

type Message struct {
	Element Data
	inline  bool
	parsed  bool
}

func parseElement(r *Reader) (Data, error) {
	b, err := r.readByte()
	if err != nil {
		return nil, err
	}
	factory, ok := respElementFactory[b]
	if !ok {
		return nil, malformedRespMessageError{}
	}
	return factory(), nil
}

func (msg *Message) Parse(r *Reader) (bool, error) {
	if msg.parsed {
		return true, nil
	}
	if msg.Element == nil {
		if !msg.inline {
			var err error
			msg.Element, err = parseElement(r)
			if err != nil {
				if _, ok := err.(malformedRespMessageError); ok {
					r.unreadByte()
					msg.inline = true
				} else {
					return false, err
				}
			}
		}
		var err error
		if msg.inline {
			msg.Element, err = parseInline(r)
			if err != nil {
				return false, err
			}
			msg.parsed = true
		} else {
			msg.parsed, err = msg.Element.parse(r)
		}
		return msg.parsed, err
	}
	var err error
	msg.parsed, err = msg.Element.parse(r)
	return msg.parsed, err
}

func parseInline(r *Reader) (Data, error) {
	bytes, err := r.readLine()
	if err != nil {
		return nil, err
	}
	args := strings.Fields(string(bytes))
	items := make([]Data, len(args))
	for i, arg := range args {
		items[i] = &BulkString{Data: []byte(arg), parsed: true}
	}
	return &Array{Items: items, parsed: true}, nil
}

func (msg *Message) Ser() []byte {
	size := msg.Element.serLen()
	bytes := make([]byte, 0, size)
	bytes = msg.Element.ser(bytes)
	return bytes
}
