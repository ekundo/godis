package shared

type TypedError interface {
	ErrorType() string
}
