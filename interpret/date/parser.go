// Package date interprets common human date formats into time.Time values.
package date

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/codedbyhassan/gokit/interpret"
)

var numericDate = regexp.MustCompile(`^(\d{1,4})[-/.](\d{1,2})[-/.](\d{1,4})$`)

// Parse interprets a date without requiring one exact input format.
func Parse(input string) (interpret.Result[time.Time], error) {
	original := input
	input = interpret.Normalize(input)
	if input == "" {
		return interpret.Result[time.Time]{OriginalInput: original}, interpret.ErrEmptyInput
	}
	if value, layout, ok := parseNamed(input); ok {
		return interpret.Result[time.Time]{Value: value, Format: layout, Confidence: interpret.HighConfidence, OriginalInput: original}, nil
	}
	parts := numericDate.FindStringSubmatch(input)
	if parts == nil {
		return interpret.Result[time.Time]{OriginalInput: original}, fmt.Errorf("%w: %q", interpret.ErrUnrecognizedInput, input)
	}
	numbers := make([]int, 3)
	for i := 0; i < 3; i++ {
		n, err := strconv.Atoi(parts[i+1])
		if err != nil {
			return interpret.Result[time.Time]{OriginalInput: original}, fmt.Errorf("%w: %q", interpret.ErrUnrecognizedInput, input)
		}
		numbers[i] = n
	}
	candidates := buildCandidates(numbers)
	if len(candidates) == 0 {
		return interpret.Result[time.Time]{OriginalInput: original}, fmt.Errorf("%w: %q", interpret.ErrUnrecognizedInput, input)
	}
	if len(candidates) > 1 && candidates[0].score == candidates[1].score {
		result := interpret.Result[time.Time]{Format: "ambiguous", Confidence: candidates[0].score, Ambiguous: true, OriginalInput: original}
		for _, candidate := range candidates {
			result.Candidates = append(result.Candidates, interpret.Candidate{Format: candidate.format, Confidence: candidate.score, Reason: candidate.reason})
		}
		return result, interpret.ErrAmbiguousInput
	}
	best := candidates[0]
	result := interpret.Result[time.Time]{Value: best.value, Format: best.format, Confidence: best.score, OriginalInput: original}
	if best.assumption != "" {
		result.Assumptions = []string{best.assumption}
	}
	for _, candidate := range candidates {
		result.Candidates = append(result.Candidates, interpret.Candidate{Format: candidate.format, Confidence: candidate.score, Reason: candidate.reason})
	}
	return result, nil
}

type candidate struct {
	value time.Time
	format string
	score interpret.Confidence
	reason string
	assumption string
}

func buildCandidates(n []int) []candidate {
	permutations := []struct { order [3]int; format, reason string }{
		{[3]int{0, 1, 2}, "DD-MM-YYYY", "day-month-year"},
		{[3]int{1, 0, 2}, "MM-DD-YYYY", "month-day-year"},
		{[3]int{2, 1, 0}, "YYYY-MM-DD", "year-month-day"},
		{[3]int{0, 2, 1}, "DD-YYYY-MM", "day-year-month"},
		{[3]int{1, 2, 0}, "MM-YYYY-DD", "month-year-day"},
		{[3]int{2, 0, 1}, "YYYY-DD-MM", "year-day-month"},
	}
	var out []candidate
	for _, p := range permutations {
		year := n[p.order[2]]
		month := n[p.order[1]]
		day := n[p.order[0]]
		if year < 1000 || year > 9999 || month < 1 || month > 12 || day < 1 || day > 31 {
			continue
		}
		value := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		if value.Year() != year || int(value.Month()) != month || value.Day() != day {
			continue
		}
		score := interpret.Confidence(0.70)
		assumption := ""
		if year >= 1000 {
			score = 0.95
			assumption = fmt.Sprintf("%d identified as the year", year)
		}
		if strings.Contains(p.format, "YYYY") && (n[0] >= 1000 || n[1] >= 1000 || n[2] >= 1000) {
			score = 0.98
		}
		out = append(out, candidate{value: value, format: p.format, score: score, reason: p.reason, assumption: assumption})
	}

	// Equivalent dates are not ambiguous even if their format labels differ.
	unique := make(map[string]candidate, len(out))
	for _, item := range out {
		key := item.value.Format("2006-01-02")
		if existing, ok := unique[key]; !ok || item.score > existing.score {
			unique[key] = item
		}
	}
	out = out[:0]
	for _, item := range unique {
		out = append(out, item)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].score > out[j].score })
	return out
}

func parseNamed(input string) (time.Time, string, bool) {
	clean := strings.ToLower(strings.TrimSpace(input))
	layouts := []struct { layout, name string }{
		{"2 January 2006", "D Month YYYY"},
		{"2 Jan 2006", "D Mon YYYY"},
		{"January 2 2006", "Month D YYYY"},
		{"Jan 2 2006", "Mon D YYYY"},
		{"2006 January 2", "YYYY Month D"},
		{"2006-01-02", "YYYY-MM-DD"},
		{"02-01-2006", "DD-MM-YYYY"},
		{"01/02/2006", "DD/MM/YYYY"},
	}
	for _, item := range layouts {
		if value, err := time.Parse(item.layout, clean); err == nil {
			return value, item.name, true
		}
	}
	return time.Time{}, "", false
}
