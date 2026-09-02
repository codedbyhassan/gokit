// Package intent provides a deterministic natural-language command layer over GoKit interpreters.
package intent

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/codedbyhassan/gokit/age"
	"github.com/codedbyhassan/gokit/calculator"
	"github.com/codedbyhassan/gokit/interpret"
	"github.com/codedbyhassan/gokit/interpret/date"
	"github.com/codedbyhassan/gokit/interpret/expression"
	"github.com/codedbyhassan/gokit/interpret/unit"
	"github.com/codedbyhassan/gokit/units"
)

// Kind identifies the user's high-level intent.
type Kind string

const (
	Calculate Kind = "calculate"
	Convert   Kind = "convert"
	Age       Kind = "age"
)

// Result is a structured high-level instruction. Value contains a domain-specific
// result: calculator.Result, unit.Conversion, or age.Result.
type Result struct {
	Kind          Kind
	Value         any
	OriginalInput string
	Format        string
	Confidence    interpret.Confidence
	Assumptions   []string
}

var quantityExpressionPattern = regexp.MustCompile(`^(.+?)\s+(plus|minus)\s+(.+)$`)
var agePattern = regexp.MustCompile(`^how old is(?: someone)?(?: born)?\s+(.+)$`)

// Parse interprets a natural-language request using the current date for age calculations.
func Parse(input string) (Result, error) {
	return ParseAt(input, time.Now())
}

// ParseAt is Parse with an explicit reference date, making age interpretation deterministic.
func ParseAt(input string, asOf time.Time) (Result, error) {
	original := input
	clean := normalize(input)
	if clean == "" {
		return Result{OriginalInput: original}, interpret.ErrEmptyInput
	}

	if match := agePattern.FindStringSubmatch(clean); len(match) == 2 {
		birth, err := date.Parse(match[1])
		if err != nil {
			return Result{OriginalInput: original}, fmt.Errorf("invalid birth date: %w", err)
		}
		calculated, err := age.Calculate(birth.Value, asOf)
		if err != nil {
			return Result{OriginalInput: original}, err
		}
		return Result{Kind: Age, Value: calculated, OriginalInput: original, Format: "age request", Confidence: 0.99}, nil
	}

	if looksLikeConversion(clean) {
		parsed, err := unit.Parse(clean)
		if err == nil {
			return Result{Kind: Convert, Value: parsed.Value, OriginalInput: original, Format: "conversion request", Confidence: 0.99}, nil
		}
	}

	if match := quantityExpressionPattern.FindStringSubmatch(clean); len(match) == 4 {
		if result, err := calculateQuantities(match[1], match[2], match[3], original); err == nil {
			return result, nil
		}
	}

	parsed, err := expression.Calculate(clean)
	if err == nil {
		return Result{Kind: Calculate, Value: parsed, OriginalInput: original, Format: "calculation request", Confidence: 0.98}, nil
	}

	return Result{OriginalInput: original}, fmt.Errorf("%w: %q", interpret.ErrUnrecognizedInput, input)
}

func normalize(input string) string {
	clean := strings.ToLower(strings.TrimSpace(input))
	clean = strings.TrimSuffix(clean, "?")
	clean = strings.Replace(clean, "what's ", "what is ", 1)
	return strings.Join(strings.Fields(clean), " ")
}

func looksLikeConversion(input string) bool {
	return strings.HasPrefix(input, "convert ") || strings.HasPrefix(input, "please convert ") ||
		strings.Contains(input, " to ") || strings.Contains(input, " into ") || strings.Contains(input, " in ")
}

func calculateQuantities(leftText, operator, rightText, original string) (Result, error) {
	left, err := parseQuantity(leftText)
	if err != nil {
		return Result{}, err
	}
	right, err := parseQuantity(rightText)
	if err != nil {
		return Result{}, err
	}
	if left.Unit.Dimension != right.Unit.Dimension {
		return Result{}, fmt.Errorf("cannot %s %s and %s", operator, left.Unit.Name, right.Unit.Name)
	}

	rightValue, err := units.Convert(right.Value, right.Unit, left.Unit)
	if err != nil {
		return Result{}, err
	}

	operation := calculator.Add
	if operator == "minus" {
		operation = calculator.Subtract
	}
	calculated, err := calculator.Calculate(left.Value, rightValue, operation)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Kind: Calculate, Value: calculated, OriginalInput: original,
		Format: "quantity calculation", Confidence: 0.99,
		Assumptions: []string{fmt.Sprintf("converted %s to %s before arithmetic", right.Unit.Symbol, left.Unit.Symbol)},
	}, nil
}

func parseQuantity(input string) (unit.Quantity, error) {
	parsed, err := unit.Parse(strings.TrimSpace(input))
	if err != nil {
		return unit.Quantity{}, err
	}
	quantity, ok := parsed.Value.(unit.Quantity)
	if !ok {
		return unit.Quantity{}, fmt.Errorf("expected a quantity, got %T", parsed.Value)
	}
	return quantity, nil
}
