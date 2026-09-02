// Package expression interprets human-friendly arithmetic expressions into calculator operations.
package expression

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/codedbyhassan/gokit/calculator"
	"github.com/codedbyhassan/gokit/interpret"
	"github.com/codedbyhassan/gokit/interpret/number"
)

// Operator identifies the arithmetic operation to perform.
type Operator string

const (
	Add         Operator = "+"
	Subtract    Operator = "-"
	Multiply    Operator = "*"
	Divide      Operator = "/"
	PercentOf   Operator = "percent_of"
	PercentIncrease Operator = "percent_increase"
)

// Operation is a structured arithmetic instruction produced by the interpreter.
type Operation struct {
	Left          float64
	Right         float64
	Operator      Operator
	OriginalInput string
}

var binaryPattern = regexp.MustCompile(`^(.+?)\s+(?:plus|add(?:ed)?\s+to|minus|subtract(?:ed)?\s+from|multiplied\s+by|times|divided\s+by)\s+(.+)$`)

// Parse interprets common arithmetic expressions such as "10 plus 5",
// "what is 10 multiplied by 5", and "100 divided by 4".
func Parse(input string) (interpret.Result[Operation], error) {
	original := input
	clean := normalize(input)
	if clean == "" {
		return interpret.Result[Operation]{OriginalInput: original}, interpret.ErrEmptyInput
	}

	clean = strings.TrimSpace(strings.TrimPrefix(clean, "what is "))
	clean = strings.TrimSuffix(clean, "?")

	if result, ok := parsePercent(clean, original); ok {
		return result, nil
	}
	if result, ok := parseHalf(clean, original); ok {
		return result, nil
	}

	if left, op, right, ok := parseSymbolic(clean); ok {
		return makeResult(left, right, op, original), nil
	}
	if leftText, rightText, op, ok := parseWords(clean); ok {
		left, err := number.Parse(leftText)
		if err != nil {
			return interpret.Result[Operation]{OriginalInput: original}, fmt.Errorf("invalid left operand: %w", err)
		}
		right, err := number.Parse(rightText)
		if err != nil {
			return interpret.Result[Operation]{OriginalInput: original}, fmt.Errorf("invalid right operand: %w", err)
		}
		return makeResult(left.Value, right.Value, op, original), nil
	}

	return interpret.Result[Operation]{OriginalInput: original}, fmt.Errorf("%w: %q", interpret.ErrUnrecognizedInput, input)
}

// Calculate parses an expression and executes the resulting operation.
func Calculate(input string) (calculator.Result, error) {
	parsed, err := Parse(input)
	if err != nil {
		return calculator.Result{}, err
	}

	if parsed.Value.Operator == PercentOf {
		value := parsed.Value.Left * parsed.Value.Right / 100
		return calculator.Result{Left: parsed.Value.Left, Right: parsed.Value.Right, Operation: calculator.Multiply, Value: value}, nil
	}
	if parsed.Value.Operator == PercentIncrease {
		value := parsed.Value.Left + parsed.Value.Left*parsed.Value.Right/100
		return calculator.Result{Left: parsed.Value.Left, Right: parsed.Value.Right, Operation: calculator.Add, Value: value}, nil
	}

	return calculator.Calculate(parsed.Value.Left, parsed.Value.Right, calculator.Operation(parsed.Value.Operator))
}

func normalize(input string) string {
	return strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(input))), " ")
}

func parseSymbolic(input string) (float64, Operator, float64, bool) {
	operators := []struct {
		marker string
		op     Operator
	}{
		{"+", Add}, {"-", Subtract}, {"*", Multiply}, {"/", Divide},
	}
	for _, item := range operators {
		parts := strings.Split(input, item.marker)
		if len(parts) != 2 {
			continue
		}
		left, leftErr := number.Parse(strings.TrimSpace(parts[0]))
		right, rightErr := number.Parse(strings.TrimSpace(parts[1]))
		if leftErr == nil && rightErr == nil {
			return left.Value, item.op, right.Value, true
		}
	}
	return 0, "", 0, false
}

func parseWords(input string) (string, string, Operator, bool) {
	patterns := []struct {
		token string
		op    Operator
	}{
		{" plus ", Add},
		{" added to ", Add},
		{" add to ", Add},
		{" minus ", Subtract},
		{" subtract from ", Subtract},
		{" multiplied by ", Multiply},
		{" times ", Multiply},
		{" divided by ", Divide},
	}
	for _, item := range patterns {
		if parts := strings.Split(input, item.token); len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])
			if left != "" && right != "" {
				return left, right, item.op, true
			}
		}
	}
	return "", "", "", false
}

func parsePercent(input, original string) (interpret.Result[Operation], bool) {
	parts := strings.Split(input, " of ")
	if len(parts) == 2 && strings.HasSuffix(parts[0], "%") {
		percent := strings.TrimSuffix(strings.TrimSpace(parts[0]), "%")
		base := strings.TrimSpace(parts[1])
		p, err1 := number.Parse(percent)
		b, err2 := number.Parse(base)
		if err1 == nil && err2 == nil {
			return makeResult(b.Value, p.Value, PercentOf, original), true
		}
	}

	if strings.HasSuffix(input, "% increase") {
		percent := strings.TrimSpace(strings.TrimSuffix(input, "% increase"))
		parts := strings.Split(percent, " increase on ")
		if len(parts) == 2 {
			p, err1 := number.Parse(strings.TrimSpace(parts[0]))
			b, err2 := number.Parse(strings.TrimSpace(parts[1]))
			if err1 == nil && err2 == nil {
				return makeResult(b.Value, p.Value, PercentIncrease, original), true
			}
		}
	}
	return interpret.Result[Operation]{}, false
}

func parseHalf(input, original string) (interpret.Result[Operation], bool) {
	if strings.HasPrefix(input, "half of ") {
		base, err := number.Parse(strings.TrimSpace(strings.TrimPrefix(input, "half of ")))
		if err == nil {
			return makeResult(base.Value, 2, Divide, original), true
		}
	}
	return interpret.Result[Operation]{}, false
}

func makeResult(left, right float64, op Operator, original string) interpret.Result[Operation] {
	return interpret.Result[Operation]{
		Value: Operation{Left: left, Right: right, Operator: op, OriginalInput: original},
		Format: "arithmetic expression",
		Confidence: interpret.HighConfidence,
		OriginalInput: original,
	}
}
