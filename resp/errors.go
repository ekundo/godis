package resp

import (
	"fmt"
	"github.com/ekundo/godis/shared"
)

type malformedRespMessageError struct {
}

func (f malformedRespMessageError) Error() string {
	return fmt.Sprint("can't parse RESP message")
}

func (f malformedRespMessageError) ErrorType() string {
	return "BADRESP"
}

var _ shared.TypedError = malformedRespMessageError{}
