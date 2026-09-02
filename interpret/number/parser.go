// Package number interprets human-friendly numeric input into float64 values.
package number

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/codedbyhassan/gokit/interpret"
)

var (
	numeric = regexp.MustCompile(`^[+-]?(?:\d+(?:\.\d+)?|\.\d+)$`)
	groupedNumeric = regexp.MustCompile(`^[+-]?\d{1,3}(?: \d{3})+$`)
	suffixedNumeric = regexp.MustCompile(`^([+-]?(?:\d+(?:\.\d+)?|\.\d+))([kmbt])$`)
	wordScale = map[string]float64{
		"thousand":    1e3,
		"million":     1e6,
		"billion":     1e9,
		"trillion":    1e12,
		"quadrillion": 1e15,
	}
	wordNumbers = map[string]int64{
		"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5,
		"six": 6, "seven": 7, "eight": 8, "nine": 9, "ten": 10, "eleven": 11,
		"twelve": 12, "thirteen": 13, "fourteen": 14, "fifteen": 15, "sixteen": 16,
		"seventeen": 17, "eighteen": 18, "nineteen": 19, "twenty": 20, "thirty": 30,
		"forty": 40, "fifty": 50, "sixty": 60, "seventy": 70, "eighty": 80, "ninety": 90,
	}
)

func Parse(input string) (interpret.Result[float64], error) {
	original := input
	clean := normalize(input)
	if clean == "" {
		return interpret.Result[float64]{OriginalInput: original}, interpret.ErrEmptyInput
	}

	if groupedNumeric.MatchString(clean) {
		clean = strings.ReplaceAll(clean, " ", "")
	}
	if numeric.MatchString(clean) {
		value, err := strconv.ParseFloat(clean, 64)
		if err == nil && !math.IsNaN(value) && !math.IsInf(value, 0) {
			return result(value, "numeric", interpret.HighConfidence, original), nil
		}
	}

	if value, ok := parseSuffixed(clean); ok {
		return result(value, "scaled numeric", interpret.Confidence(0.98), original), nil
	}
	if value, format, ok := parseWords(clean); ok {
		return result(value, format, interpret.Confidence(0.96), original), nil
	}

	return interpret.Result[float64]{OriginalInput: original}, fmt.Errorf("%w: %q", interpret.ErrUnrecognizedInput, input)
}

func normalize(input string) string {
	s := strings.ToLower(strings.TrimSpace(input))
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "–", "-")
	s = strings.ReplaceAll(s, "—", "-")
	return strings.Join(strings.Fields(s), " ")
}

func result(value float64, format string, confidence interpret.Confidence, original string) interpret.Result[float64] {
	return interpret.Result[float64]{Value: value, Format: format, Confidence: confidence, OriginalInput: original}
}

func parseSuffixed(input string) (float64, bool) {
	if matches := suffixedNumeric.FindStringSubmatch(input); len(matches) == 3 {
		base, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return 0, false
		}
		scales := map[string]float64{"k": 1e3, "m": 1e6, "b": 1e9, "t": 1e12}
		return base * scales[matches[2]], true
	}

	parts := strings.Fields(input)
	if len(parts) < 2 {
		return 0, false
	}

	// A word-scaled number is made of descending scale groups:
	// "two million five hundred thousand" => 2*1e6 + 500*1e3.
	var total float64
	var group []string
	lastScale := math.Inf(1)
	foundScale := false

	for _, token := range parts {
		if scale, ok := wordScale[token]; ok {
			if scale >= lastScale || len(group) == 0 {
				return 0, false
			}
			base, ok := parseBase(strings.Join(group, " "))
			if !ok {
				return 0, false
			}
			total += base * scale
			group = nil
			lastScale = scale
			foundScale = true
			continue
		}
		group = append(group, token)
	}

	if !foundScale {
		return 0, false
	}
	if len(group) > 0 {
		base, ok := parseBase(strings.Join(group, " "))
		if !ok {
			return 0, false
		}
		total += base
	}
	return total, true
}

func parseBase(input string) (float64, bool) {
	input = strings.TrimSpace(input)
	if numeric.MatchString(input) {
		value, err := strconv.ParseFloat(input, 64)
		return value, err == nil
	}
	value, _, ok := parseWords(input)
	return value, ok
}

func parseWords(input string) (float64, string, bool) {
	parts := strings.Fields(strings.ReplaceAll(input, "-", " "))
	if len(parts) == 0 {
		return 0, "", false
	}

	negative := false
	if parts[0] == "minus" || parts[0] == "negative" {
		negative = true
		parts = parts[1:]
	}
	if len(parts) == 0 {
		return 0, "", false
	}

	point := -1
	for i, token := range parts {
		if token == "point" || token == "dot" {
			if point != -1 {
				return 0, "", false
			}
			point = i
		}
	}

	if point >= 0 {
		if point == 0 || point == len(parts)-1 {
			return 0, "", false
		}
		whole, ok := parseIntegerWords(parts[:point])
		if !ok {
			return 0, "", false
		}
		fractionDigits := make([]string, 0, len(parts)-point-1)
		for _, token := range parts[point+1:] {
			if value, ok := wordNumbers[token]; ok && value >= 0 && value <= 9 {
				fractionDigits = append(fractionDigits, strconv.FormatInt(value, 10))
				continue
			}
			if numeric.MatchString(token) && !strings.Contains(token, ".") {
				fractionDigits = append(fractionDigits, token)
				continue
			}
			return 0, "", false
		}
		fraction, err := strconv.ParseFloat("0."+strings.Join(fractionDigits, ""), 64)
		if err != nil {
			return 0, "", false
		}
		value := whole + fraction
		if negative {
			value = -value
		}
		return value, "number words", true
	}

	value, ok := parseIntegerWords(parts)
	if !ok {
		return 0, "", false
	}
	if negative {
		value = -value
	}
	return value, "number words", true
}

func parseIntegerWords(parts []string) (float64, bool) {
	if len(parts) == 0 {
		return 0, false
	}

	var total int64
	var group int64
	used := false
	lastWasHundred := false

	for i, token := range parts {
		if token == "and" {
			if i == 0 || i == len(parts)-1 {
				return 0, false
			}
			continue
		}

		if value, ok := wordNumbers[token]; ok {
			// A tens word may be followed by a ones word, e.g. "twenty five".
			// Reject repeated tens or malformed repeated small-number sequences.
			if value >= 20 && value%10 == 0 {
				if group%100 >= 20 && group%10 == 0 {
					return 0, false
				}
				group += value
			} else if value < 20 {
				if lastWasHundred && value == 0 {
					return 0, false
				}
				group += value
			} else {
				group += value
			}
			used = true
			lastWasHundred = false
			continue
		}

		if token == "hundred" {
			if group == 0 || lastWasHundred {
				return 0, false
			}
			group *= 100
			lastWasHundred = true
			used = true
			continue
		}

		if scale, ok := wordScale[token]; ok {
			if group == 0 || scale <= 1 {
				return 0, false
			}
			if total > 0 && scale >= 1e3 && float64(total) >= scale {
				return 0, false
			}
			total += group * int64(scale)
			group = 0
			lastWasHundred = false
			used = true
			continue
		}

		return 0, false
	}

	if !used {
		return 0, false
	}
	return float64(total + group), true
}

func parseAtomicNumber(input string) (float64, error) {
	if numeric.MatchString(input) {
		return strconv.ParseFloat(input, 64)
	}
	if value, _, ok := parseWords(input); ok {
		return value, nil
	}
	return 0, fmt.Errorf("not a number")
}
