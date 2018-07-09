package resp

import (
	"io"
	"math"
)

type Reader struct {
	io.Reader
	buf  []byte
	ridx int
	lidx int
	widx int
}

func NewReader(r io.Reader) *Reader {
	return &Reader{Reader: r, buf: make([]byte, 1024), ridx: 0, widx: 0}
}

func (r *Reader) peek() (byte, error) {
	if r.widx > r.ridx {
		b := r.buf[r.ridx]
		r.lidx = r.ridx
		r.ridx++
		return b, nil
	}
	if r.widx >= len(r.buf) {
		return 0, malformedRespMessageError{}
	}
	n, err := r.Read(r.buf[r.widx : r.widx+1])
	r.widx += n
	if err != nil {
		return 0, err
	}
	return r.peek()
}

func (r *Reader) readInt() (int32, error) {
	r.ridx = 0
	b, err := r.peek()
	if err != nil {
		return 0, err
	}

	sign := int(1)
	if b == '-' {
		sign = -1
		b, err = r.peek()
		if err != nil {
			return 0, err
		}
	}

	var i int
	for b >= '0' && b <= '9' {
		i = i*10 + int(b-'0')
		if i > math.MaxInt32+1 {
			return 0, malformedRespMessageError{}
		}
		b, err = r.peek()
		if err != nil {
			return 0, err
		}
	}

	if b != '\r' {
		return 0, malformedRespMessageError{}
	}
	b, err = r.peek()
	if err != nil {
		return 0, err
	}
	if b != '\n' {
		return 0, malformedRespMessageError{}
	}

	r.ridx = 0
	r.widx = 0

	res := i * sign
	if res > math.MaxInt32 || res < math.MinInt32 {
		return 0, malformedRespMessageError{}
	}

	return int32(res), nil
}

func (r *Reader) readLine() ([]byte, error) {
	r.ridx = 0
	for {
		b, err := r.peek()
		if err != nil {
			return nil, err
		}

		if b == '\r' {
			b, err = r.peek()
			if err != nil {
				return nil, err
			}
			if b != '\n' {
				return nil, malformedRespMessageError{}
			}

			res := r.buf[:r.ridx-2]

			r.ridx = 0
			r.widx = 0

			return res, nil
		}
	}
}

func (r *Reader) readCRLF() error {
	r.ridx = 0

	b, err := r.peek()
	if err != nil {
		return err
	}

	if b != '\r' {
		return malformedRespMessageError{}
	}

	b, err = r.peek()
	if err != nil {
		return err
	}
	if b != '\n' {
		return malformedRespMessageError{}
	}

	r.ridx = 0
	r.widx = 0

	return nil
}

func (r *Reader) readByte() (byte, error) {
	r.ridx = 0

	b, err := r.peek()
	if err != nil {
		return 0, err
	}

	r.ridx = 0
	r.widx = 0

	return b, nil
}

func (r *Reader) unreadByte() {
	r.ridx = r.lidx
	r.widx = r.ridx + 1
}
