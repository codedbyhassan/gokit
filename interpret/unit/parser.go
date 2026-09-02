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

type Quantity struct {
	Value float64
	Unit  units.Unit
}

type Conversion struct {
	Value         float64
	From          Quantity
	To            units.Unit
	OriginalInput string
}

var conversionPattern = regexp.MustCompile(`^(.+?)\s+(?:to|in)\s+(.+)$`)
var quantityPattern = regexp.MustCompile(`^([+-]?(?:\d+(?:\.\d+)?|\.\d+))\s*([a-z°][a-z°/ ]*)$`)

// Parse interprets measurements such as "5kg", "10 kilometers", and
// conversions such as "20 km to miles" or "25 celsius in fahrenheit".
func Parse(input string) (interpret.Result[any], error) {
	original := input
	clean := normalize(input)
	if clean == "" {
		return interpret.Result[any]{OriginalInput: original}, interpret.ErrEmptyInput
	}

	if matches := conversionPattern.FindStringSubmatch(clean); len(matches) == 3 {
		from, err := parseQuantity(matches[1])
		if err != nil { return interpret.Result[any]{OriginalInput: original}, err }
		to, err := units.Lookup(matches[2])
		if err != nil { return interpret.Result[any]{OriginalInput: original}, err }
		value, err := units.Convert(from.Value, from.Unit, to)
		if err != nil { return interpret.Result[any]{OriginalInput: original}, err }
		return interpret.Result[any]{
			Value: Conversion{Value: value, From: from, To: to, OriginalInput: original},
			Format: "unit conversion", Confidence: 0.99, OriginalInput: original,
		}, nil
	}

	quantity, err := parseQuantity(clean)
	if err != nil { return interpret.Result[any]{OriginalInput: original}, err }
	return interpret.Result[any]{Value: quantity, Format: "quantity", Confidence: 0.97, OriginalInput: original}, nil
}

func normalize(input string) string {
	clean := strings.ToLower(strings.TrimSpace(input))
	clean = strings.ReplaceAll(clean, "°c", " celsius")
	clean = strings.ReplaceAll(clean, "°f", " fahrenheit")
	clean = strings.Join(strings.Fields(clean), " ")
	// Undo the intentional spacing when the degree symbol was attached to a number.
	clean = strings.ReplaceAll(clean, " celsius", " celsius")
	return clean
}

func parseQuantity(input string) (Quantity, error) {
	clean := strings.TrimSpace(input)
	matches := quantityPattern.FindStringSubmatch(clean)
	if len(matches) != 3 {
		return Quantity{}, fmt.Errorf("invalid measurement %q", input)
	}
	value, err := strconv.ParseFloat(matches[1], 64)
	if err != nil { return Quantity{}, fmt.Errorf("invalid measurement value %q: %w", matches[1], err) }
	unit, err := units.Lookup(strings.TrimSpace(matches[2]))
	if err != nil { return Quantity{}, err }
	return Quantity{Value: value, Unit: unit}, nil
}
