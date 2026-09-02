// Package unit interprets human-friendly measurements and conversions.
package unit

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/codedbyhassan/gokit/interpret"
	"github.com/codedbyhassan/gokit/units"
)

// Quantity represents a numeric value paired with a unit.
type Quantity struct {
	Value float64
	Unit  units.Unit
}

// Conversion represents a conversion from one unit to another.
type Conversion struct {
	Value         float64
	From          Quantity
	To            units.Unit
	OriginalInput string
}

var conversionPattern = regexp.MustCompile(`^(.+?)\s+(?:to|in)\s+(.+)$`)
var quantityPattern = regexp.MustCompile(`^([+-]?(?:\d+(?:\.\d+)?|\.\d+))\s*([a-z°/ ]+)$`)

// Parse interprets a quantity or conversion such as "5 kg", "10 kilometers",
// or "20 km to miles".
func Parse(input string) (interpret.Result[any], error) {
	original := input
	clean := strings.ToLower(strings.TrimSpace(input))
	clean = strings.ReplaceAll(clean, "°", "°")
	if clean == "" {
		return interpret.Result[any]{OriginalInput: original}, interpret.ErrEmptyInput
	}

	if matches := conversionPattern.FindStringSubmatch(clean); len(matches) == 3 {
		from, err := parseQuantity(matches[1])
		if err != nil {
			return interpret.Result[any]{OriginalInput: original}, err
		}
		to, err := units.Lookup(strings.TrimSpace(matches[2]))
		if err != nil {
			return interpret.Result[any]{OriginalInput: original}, err
		}
		value, err := units.Convert(from.Value, from.Unit, to)
		if err != nil {
			return interpret.Result[any]{OriginalInput: original}, err
		}
		return interpret.Result[any]{
			Value: Conversion{Value: value, From: from, To: to, OriginalInput: original},
			Format: "conversion",
			Confidence: 0.98,
			OriginalInput: original,
		}, nil
	}

	quantity, err := parseQuantity(clean)
	if err != nil {
		return interpret.Result[any]{OriginalInput: original}, err
	}
	return interpret.Result[any]{
		Value: quantity,
		Format: "quantity",
		Confidence: 0.96,
		OriginalInput: original,
	}, nil
}

func parseQuantity(input string) (Quantity, error) {
	matches := quantityPattern.FindStringSubmatch(strings.TrimSpace(input))
	if len(matches) != 3 {
		return Quantity{}, fmt.Errorf("invalid measurement %q", input)
	}
	value, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return Quantity{}, fmt.Errorf("invalid measurement value %q: %w", matches[1], err)
	}
	unit, err := units.Lookup(strings.TrimSpace(matches[2]))
	if err != nil {
		return Quantity{}, err
	}
	return Quantity{Value: value, Unit: unit}, nil
}
