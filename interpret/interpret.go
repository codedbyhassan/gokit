// Package interpret provides flexible, explainable input interpretation primitives.
package interpret

// Confidence describes how strongly an input matches an interpretation.
type Confidence float64

const (
	LowConfidence    Confidence = 0.40
	MediumConfidence Confidence = 0.70
	HighConfidence   Confidence = 0.90
)

// Result is returned by interpreters that turn human input into typed values.
// Assumptions explains any non-obvious decisions made by the interpreter.
type Result[T any] struct {
	Value        T
	Format       string
	Confidence   Confidence
	Ambiguous    bool
	Assumptions  []string
	Candidates   []Candidate
	OriginalInput string
}

// Candidate describes one possible interpretation of an input.
type Candidate struct {
	Format     string
	Confidence Confidence
	Reason     string
}
