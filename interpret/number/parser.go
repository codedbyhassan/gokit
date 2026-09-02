// Package number interprets human-friendly numeric input into float64 values.
package number

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/codedbyhassan/gokit/interpret"
)

var numeric = regexp.MustCompile(`^[+-]?(?:\d+(?:\.\d+)?|\.\d+)$`)

var scaleWords = map[string]float64{
	"thousand": 1_000,
	"million": 1_000_000,
	"billion": 1_000_000_000,
	"trillion": 1_000_000_000_000,
}

var wordNumbers = map[string]float64{
	"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4,
	"five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9,
	"ten": 10, "eleven": 11, "twelve": 12, "thirteen": 13,
	"fourteen": 14, "fifteen": 15, "sixteen": 16, "seventeen": 17,
	"eighteen": 18, "nineteen": 19, "twenty": 20, "thirty": 30,
	"forty": 40, "fifty": 50, "sixty": 60, "seventy": 70,
	"eighty": 80, "ninety": 90,
}

// Parse interprets common numeric representations such as 1,000, 1.5k,
// 2 million, and English number words.
func Parse(input string) (interpret.Result[float64], error) {
	original := input
	clean := strings.ToLower(strings.TrimSpace(input))
	clean = strings.ReplaceAll(clean, ",", "")
	clean = strings.ReplaceAll(clean, "_", "")
	clean = strings.Join(strings.Fields(clean), " ")

	if clean == "" {
		return interpret.Result[float64]{OriginalInput: original}, interpret.ErrEmptyInput
	}

	if numeric.MatchString(clean) {
		value, err := strconv.ParseFloat(clean, 64)
		if err != nil {
			return interpret.Result[float64]{OriginalInput: original}, fmt.Errorf("%w: %q", interpret.ErrUnrecognizedInput, input)
		}
		return interpret.Result[float64]{Value: value, Format: "numeric", Confidence: interpret.HighConfidence, OriginalInput: original}, nil
	}

	if value, ok := parseSuffixed(clean); ok {
		return interpret.Result[float64]{Value: value, Format: "scaled numeric", Confidence: 0.98, OriginalInput: original}, nil
	}

	if value, format, ok := parseWords(clean); ok {
		return interpret.Result[float64]{Value: value, Format: format, Confidence: 0.96, OriginalInput: original}, nil
	}

	return interpret.Result[float64]{OriginalInput: original}, fmt.Errorf("%w: %q", interpret.ErrUnrecognizedInput, input)
}

func parseSuffixed(input string) (float64, bool) {
	suffixes := []struct {
		suffix string
		scale  float64
	}{
		{"k", 1_000}, {"m", 1_000_000}, {"b", 1_000_000_000},
	}
	for _, item := range suffixes {
		if !strings.HasSuffix(input, item.suffix) {
			continue
		}
		n := strings.TrimSpace(strings.TrimSuffix(input, item.suffix))
		if !numeric.MatchString(n) {
			return 0, false
		}
		value, err := strconv.ParseFloat(n, 64)
		if err != nil {
			return 0, false
		}
		return value * item.scale, true
	}

	parts := strings.Fields(input)
	if len(parts) == 2 {
		if scale, ok := scaleWords[parts[1]]; ok && numeric.MatchString(parts[0]) {
			value, err := strconv.ParseFloat(parts[0], 64)
			if err == nil {
				return value * scale, true
			}
		}
	}
	return 0, false
}

func parseWords(input string) (float64, string, bool) {
	parts := strings.Fields(strings.ReplaceAll(input, "-", " "))
	if len(parts) == 0 {
		return 0, "", false
	}

	var total float64
	var current float64
	used := false
	for _, part := range parts {
		if part == "and" {
			continue
		}
		if value, ok := wordNumbers[part]; ok {
			current += value
			used = true
			continue
		}
		if part == "hundred" {
			if current == 0 {
				current = 1
			}
			current *= 100
			used = true
			continue
		}
		if scale, ok := scaleWords[part]; ok {
			if current == 0 {
				current = 1
			}
			total += current * scale
			current = 0
			used = true
			continue
		}
		return 0, "", false
	}
	if !used {
		return 0, "", false
	}
	return total + current, "number words", true
}
