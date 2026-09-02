package calculator

import (
	"errors"
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name      string
		left      float64
		right     float64
		operation Operation
		want      float64
	}{
		{"add", 10, 5, Add, 15},
		{"subtract", 10, 5, Subtract, 5},
		{"multiply", 10, 5, Multiply, 50},
		{"divide", 10, 5, Divide, 2},
		{"modulo", 10, 3, Modulo, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Calculate(tt.left, tt.right, tt.operation)
			if err != nil {
				t.Fatalf("Calculate() error = %v", err)
			}
			if result.Value != tt.want {
				t.Fatalf("Calculate() value = %v, want %v", result.Value, tt.want)
			}
		})
	}
}

func TestCalculateDivisionByZero(t *testing.T) {
	_, err := Calculate(10, 0, Divide)
	if !errors.Is(err, ErrDivisionByZero) {
		t.Fatalf("expected ErrDivisionByZero, got %v", err)
	}
}

func TestCalculateInvalidOperation(t *testing.T) {
	_, err := Calculate(10, 5, Operation("invalid"))
	if !errors.Is(err, ErrInvalidOperation) {
		t.Fatalf("expected ErrInvalidOperation, got %v", err)
	}
}
