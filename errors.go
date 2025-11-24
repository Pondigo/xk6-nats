package nats

import (
	"fmt"
)

type NatsError struct {
	Code    int
	Message string
	Err     error
}

func (e *NatsError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("nats error [%d]: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("nats error [%d]: %s", e.Code, e.Message)
}

func (e *NatsError) Unwrap() error {
	return e.Err
}

func NewNatsError(code int, message string, err error) *NatsError {
	return &NatsError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func NewConnectionError(message string, err error) *NatsError {
	return NewNatsError(1002, message, err)
}

var (
	ErrNoVUState        = NewNatsError(1001, "no VU state available", nil)
	ErrConnectionClosed = NewNatsError(1002, "connection is closed", nil)
	ErrInvalidConfig    = NewNatsError(1003, "invalid configuration", nil)
	ErrStreamNotFound   = NewNatsError(1004, "stream not found", nil)
	ErrConsumerNotFound = NewNatsError(1005, "consumer not found", nil)
	ErrTimeout          = NewNatsError(1006, "operation timed out", nil)
	ErrNoMessage        = NewNatsError(1007, "no message available", nil)
)
