// Package units provides reusable unit definitions and conversions.
package units

import (
	"fmt"
	"strings"
)

// Dimension identifies the physical dimension of a unit.
type Dimension string

const (
	Length      Dimension = "length"
	Mass        Dimension = "mass"
	Temperature Dimension = "temperature"
	Speed       Dimension = "speed"
)

// Unit identifies a supported unit and its dimension.
type Unit struct {
	Name      string
	Symbol    string
	Dimension Dimension
}

var definitions = map[string]Unit{
	"meter": {"meter", "m", Length}, "m": {"meter", "m", Length}, "meters": {"meter", "m", Length},
	"kilometer": {"kilometer", "km", Length}, "km": {"kilometer", "km", Length}, "kilometers": {"kilometer", "km", Length},
	"mile": {"mile", "mi", Length}, "mi": {"mile", "mi", Length}, "miles": {"mile", "mi", Length},
	"foot": {"foot", "ft", Length}, "ft": {"foot", "ft", Length}, "feet": {"foot", "ft", Length},
	"inch": {"inch", "in", Length}, "in": {"inch", "in", Length}, "inches": {"inch", "in", Length},
	"kilogram": {"kilogram", "kg", Mass}, "kg": {"kilogram", "kg", Mass}, "kilograms": {"kilogram", "kg", Mass},
	"gram": {"gram", "g", Mass}, "g": {"gram", "g", Mass}, "grams": {"gram", "g", Mass},
	"pound": {"pound", "lb", Mass}, "lb": {"pound", "lb", Mass}, "lbs": {"pound", "lb", Mass}, "pounds": {"pound", "lb", Mass},
	"celsius": {"celsius", "°C", Temperature}, "centigrade": {"celsius", "°C", Temperature}, "c": {"celsius", "°C", Temperature},
	"fahrenheit": {"fahrenheit", "°F", Temperature}, "fahrenheit degrees": {"fahrenheit", "°F", Temperature}, "f": {"fahrenheit", "°F", Temperature},
	"kilometers per hour": {"kilometers per hour", "km/h", Speed}, "km/h": {"kilometers per hour", "km/h", Speed},
	"miles per hour": {"miles per hour", "mph", Speed}, "mph": {"miles per hour", "mph", Speed},
}

// Lookup resolves a unit alias case-insensitively.
func Lookup(input string) (Unit, error) {
	key := strings.ToLower(strings.TrimSpace(input))
	unit, ok := definitions[key]
	if !ok {
		return Unit{}, fmt.Errorf("unsupported unit %q", input)
	}
	return unit, nil
}

// Convert converts a value between compatible units.
func Convert(value float64, from, to Unit) (float64, error) {
	if from.Dimension != to.Dimension {
		return 0, fmt.Errorf("cannot convert %s to %s", from.Name, to.Name)
	}
	if from.Name == to.Name {
		return value, nil
	}

	if from.Dimension == Temperature {
		return convertTemperature(value, from, to), nil
	}

	base := toBase(value, from)
	return fromBase(base, to), nil
}

func toBase(value float64, unit Unit) float64 {
	switch unit.Name {
	case "kilometer": return value * 1000
	case "mile": return value * 1609.344
	case "foot": return value * 0.3048
	case "inch": return value * 0.0254
	case "kilogram": return value * 1000
	case "pound": return value * 453.59237
	default: return value
	}
}

func fromBase(value float64, unit Unit) float64 {
	switch unit.Name {
	case "kilometer": return value / 1000
	case "mile": return value / 1609.344
	case "foot": return value / 0.3048
	case "inch": return value / 0.0254
	case "kilogram": return value / 1000
	case "pound": return value / 453.59237
	default: return value
	}
}

func convertTemperature(value float64, from, to Unit) float64 {
	if from.Name == "celsius" && to.Name == "fahrenheit" { return value*9/5 + 32 }
	if from.Name == "fahrenheit" && to.Name == "celsius" { return (value-32)*5/9 }
	return value
}
