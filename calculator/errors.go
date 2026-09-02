package calculator

import "errors"

var (
	// ErrDivisionByZero is returned when division or modulo uses zero as the divisor.
	ErrDivisionByZero = errors.New("division by zero")
	// ErrInvalidOperation is returned when an unsupported operation is requested.
	ErrInvalidOperation = errors.New("invalid operation")
)
