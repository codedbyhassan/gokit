// Package expression interprets human-friendly arithmetic expressions into calculator operations.
package expression

import (
	"fmt"
	"strings"

	"github.com/codedbyhassan/gokit/calculator"
	"github.com/codedbyhassan/gokit/interpret"
	"github.com/codedbyhassan/gokit/interpret/number"
)

type Operator string

const (
	Add Operator = "+"
	Subtract Operator = "-"
	Multiply Operator = "*"
	Divide Operator = "/"
	PercentOf Operator = "percent_of"
	PercentIncrease Operator = "percent_increase"
)

type Operation struct {
	Left float64
	Right float64
	Operator Operator
	OriginalInput string
}

func Parse(input string) (interpret.Result[Operation], error) {
	original := input
	clean := normalize(input)
	if clean == "" { return interpret.Result[Operation]{OriginalInput: original}, interpret.ErrEmptyInput }
	clean = strings.TrimSpace(strings.TrimPrefix(clean, "what is "))
	clean = strings.TrimSuffix(clean, "?")
	if result, ok := parsePercent(clean, original); ok { return result, nil }
	if result, ok := parseHalf(clean, original); ok { return result, nil }
	if left, op, right, ok := parseSymbolic(clean); ok { return makeResult(left, right, op, original), nil }
	if leftText, rightText, op, reverse, ok := parseWords(clean); ok {
		left, err := number.Parse(leftText); if err != nil { return interpret.Result[Operation]{OriginalInput: original}, fmt.Errorf("invalid left operand: %w", err) }
		right, err := number.Parse(rightText); if err != nil { return interpret.Result[Operation]{OriginalInput: original}, fmt.Errorf("invalid right operand: %w", err) }
		if reverse { left.Value, right.Value = right.Value, left.Value }
		return makeResult(left.Value, right.Value, op, original), nil
	}
	return interpret.Result[Operation]{OriginalInput: original}, fmt.Errorf("%w: %q", interpret.ErrUnrecognizedInput, input)
}

func Calculate(input string) (calculator.Result, error) {
	parsed, err := Parse(input); if err != nil { return calculator.Result{}, err }
	if parsed.Value.Operator == PercentOf { return calculator.Result{Left: parsed.Value.Left, Right: parsed.Value.Right, Operation: calculator.Multiply, Value: parsed.Value.Left * parsed.Value.Right / 100}, nil }
	if parsed.Value.Operator == PercentIncrease { return calculator.Result{Left: parsed.Value.Left, Right: parsed.Value.Right, Operation: calculator.Add, Value: parsed.Value.Left + parsed.Value.Left*parsed.Value.Right/100}, nil }
	return calculator.Calculate(parsed.Value.Left, parsed.Value.Right, calculator.Operation(parsed.Value.Operator))
}

func normalize(input string) string { return strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(input))), " ") }

func parseSymbolic(input string) (float64, Operator, float64, bool) {
	for _, item := range []struct{ marker string; op Operator }{{"+", Add}, {"-", Subtract}, {"*", Multiply}, {"/", Divide}} {
		for i := 1; i < len(input); i++ {
			if input[i:i+1] != item.marker { continue }
			leftText, rightText := strings.TrimSpace(input[:i]), strings.TrimSpace(input[i+1:])
			left, le := number.Parse(leftText); right, re := number.Parse(rightText)
			if le == nil && re == nil { return left.Value, item.op, right.Value, true }
		}
	}
	return 0, "", 0, false
}

func parseWords(input string) (string, string, Operator, bool, bool) {
	patterns := []struct{ token string; op Operator; reverse bool }{
		{" plus ", Add, false}, {" added to ", Add, false}, {" add to ", Add, false},
		{" minus ", Subtract, false}, {" subtract from ", Subtract, true},
		{" multiplied by ", Multiply, false}, {" times ", Multiply, false}, {" divided by ", Divide, false},
	}
	for _, item := range patterns { if parts := strings.Split(input, item.token); len(parts) == 2 { left, right := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]); if left != "" && right != "" { return left, right, item.op, item.reverse, true } } }
	return "", "", "", false, false
}

func parsePercent(input, original string) (interpret.Result[Operation], bool) {
	parts := strings.Split(input, " of ")
	if len(parts) == 2 && strings.HasSuffix(strings.TrimSpace(parts[0]), "%") { percent := strings.TrimSuffix(strings.TrimSpace(parts[0]), "%"); p, e1 := number.Parse(percent); b, e2 := number.Parse(strings.TrimSpace(parts[1])); if e1 == nil && e2 == nil { return makeResult(b.Value, p.Value, PercentOf, original), true } }
	for _, marker := range []string{"% increase on ", " percent increase on "} { if parts := strings.Split(input, marker); len(parts) == 2 { p, e1 := number.Parse(strings.TrimSuffix(strings.TrimSpace(parts[0]), "%")); b, e2 := number.Parse(strings.TrimSpace(parts[1])); if e1 == nil && e2 == nil { return makeResult(b.Value, p.Value, PercentIncrease, original), true } } }
	return interpret.Result[Operation]{}, false
}

func parseHalf(input, original string) (interpret.Result[Operation], bool) { if strings.HasPrefix(input, "half of ") { base, err := number.Parse(strings.TrimSpace(strings.TrimPrefix(input, "half of "))); if err == nil { return makeResult(base.Value, 2, Divide, original), true } }; return interpret.Result[Operation]{}, false }

func makeResult(left, right float64, op Operator, original string) interpret.Result[Operation] { return interpret.Result[Operation]{Value: Operation{Left:left, Right:right, Operator:op, OriginalInput:original}, Format:"arithmetic expression", Confidence:interpret.HighConfidence, OriginalInput:original} }
