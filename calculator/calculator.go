// Package calculator provides reusable arithmetic operations.
package calculator

import "fmt"

// Operation represents a supported arithmetic operation.
type Operation string

const (
	Add      Operation = "+"
	Subtract Operation = "-"
	Multiply Operation = "*"
	Divide   Operation = "/"
	Modulo   Operation = "%"
)

// Result contains the outcome of a calculation.
type Result struct {
	Left      float64
	Right     float64
	Operation Operation
	Value     float64
}

// Calculate performs an arithmetic operation on two numbers.
func Calculate(left, right float64, operation Operation) (Result, error) {
	switch operation {
	case Add:
		return Result{left, right, operation, left + right}, nil
	case Subtract:
		return Result{left, right, operation, left - right}, nil
	case Multiply:
		return Result{left, right, operation, left * right}, nil
	case Divide:
		if right == 0 {
			return Result{}, ErrDivisionByZero
		}
		return Result{left, right, operation, left / right}, nil
	case Modulo:
		if right == 0 {
			return Result{}, ErrDivisionByZero
		}
		return Result{left, right, operation, float64(int64(left) % int64(right))}, nil
	default:
		return Result{}, fmt.Errorf("unsupported operation %q: %w", operation, ErrInvalidOperation)
	}
}

// AddNumbers adds two numbers.
func AddNumbers(left, right float64) float64 { return left + right }

// SubtractNumbers subtracts right from left.
func SubtractNumbers(left, right float64) float64 { return left - right }

// MultiplyNumbers multiplies two numbers.
func MultiplyNumbers(left, right float64) float64 { return left * right }

// DivideNumbers divides left by right.
func DivideNumbers(left, right float64) (float64, error) {
	if right == 0 {
		return 0, ErrDivisionByZero
	}
	return left / right, nil
}
