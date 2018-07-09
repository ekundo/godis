package server

import (
	"fmt"
	"github.com/ekundo/godis/shared"
)

type keyNotFoundError struct {
	key string
}

func (f keyNotFoundError) Error() string {
	return fmt.Sprintf("item with key '%s' doesn't exist", f.key)
}

type fieldNotFoundError struct {
	name string
}

func (f fieldNotFoundError) Error() string {
	return fmt.Sprintf("dict field '%s' doesn't exist", f.name)
}

func (f fieldNotFoundError) ErrorType() string {
	return "FIELDNOTFOUND"
}

var _ shared.TypedError = fieldNotFoundError{}

type incompatibleTypeError struct {
}

const incompatibleTypeErrorMsg = "operation against a key holding the wrong kind of value"

func (f incompatibleTypeError) Error() string {
	return incompatibleTypeErrorMsg
}

func (f incompatibleTypeError) ErrorType() string {
	return "WRONGTYPE"
}

var _ shared.TypedError = incompatibleTypeError{}

type unknownCommandError struct {
	cmd string
}

func (f unknownCommandError) Error() string {
	return fmt.Sprintf("unknown command '%s'", f.cmd)
}

type wrongNumberOfArgumentsError struct {
	cmd string
}

func (f wrongNumberOfArgumentsError) Error() string {
	return fmt.Sprintf("wrong number of arguments for '%s' command", f.cmd)
}

type wrongArgumentTypeError struct {
}

func (f wrongArgumentTypeError) Error() string {
	return "value is not an integer or out of range"
}

type nullArgumentError struct {
}

func (f nullArgumentError) Error() string {
	return fmt.Sprintf("value is null")
}

type indexOutOfRangeError struct {
}

const indexOutOfRangeErrorMsg = "index out of range"

func (f indexOutOfRangeError) Error() string {
	return indexOutOfRangeErrorMsg
}

type syntaxError struct {
}

func (f syntaxError) Error() string {
	return fmt.Sprint("syntax error")
}

type invalidExpireTimeError struct {
}

func (f invalidExpireTimeError) Error() string {
	return fmt.Sprint("invalid expire time in set")
}
