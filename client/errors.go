package client

import (
	"fmt"
	"github.com/ekundo/godis/shared"
)

type Error struct {
	code    string
	message string
}

func (e Error) Error() string {
	return e.message
}

func (e Error) ErrorType() string {
	return e.code
}

var _ shared.TypedError = Error{}

type CommunicationError struct {
	cause error
}

func (e CommunicationError) Error() string {
	return fmt.Sprintf("communication error: %s", e.cause)
}

type UnexpectedResponseError struct {
	CommunicationError
}

func (e UnexpectedResponseError) Error() string {
	return "unexpected response received"
}

type NotConnectedError struct {
	CommunicationError
}

func (e NotConnectedError) Error() string {
	return "not connected"
}
