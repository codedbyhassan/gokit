// Package pipeline provides GoKit's unified human-input execution boundary.
package pipeline

import (
	"errors"
	"strings"
	"time"

	"github.com/codedbyhassan/gokit/interpret"
	"github.com/codedbyhassan/gokit/interpret/expression"
	"github.com/codedbyhassan/gokit/interpret/intent"
	"github.com/codedbyhassan/gokit/interpret/plan"
	"github.com/codedbyhassan/gokit/interpret/router"
)

type Source string

const (
	IntentSource Source = "intent"
	RouterSource Source = "router"
)

// Result is the unified output of interpretation and execution. Interpretation
// preserves the typed parser output; Value contains the executable result.
type Result struct {
	OriginalInput string
	Source        Source
	Kind          string
	Interpretation any
	Value         any
	Confidence    interpret.Confidence
	Assumptions   []string
	Plan          plan.Plan
}

func Parse(input string) (Result, error) { return ParseAt(input, time.Now()) }

// ParseAt is deterministic for callers that need a fixed reference time.
func ParseAt(input string, asOf time.Time) (Result, error) {
	if strings.TrimSpace(input) == "" {
		return Result{OriginalInput: input}, interpret.ErrEmptyInput
	}
	intentResult, err := intent.ParseAt(input, asOf)
	if err == nil { return fromIntent(intentResult) }
	if !errors.Is(err, interpret.ErrUnrecognizedInput) {
		return Result{OriginalInput: input}, err
	}
	routed, err := router.Parse(input)
	if err != nil { return Result{OriginalInput: input}, err }
	return fromRouter(routed)
}

func fromIntent(value intent.Result) (Result, error) {
	p := plan.Plan{OriginalInput: value.OriginalInput, Intent: string(value.Kind), Confidence: value.Confidence, Assumptions: append([]string(nil), value.Assumptions...)}
	p.Add(plan.Step{Name: "intent", Operation: string(value.Kind), Input: value.OriginalInput, Output: value.Value, Confidence: value.Confidence, Assumptions: value.Assumptions})
	result := Result{OriginalInput: value.OriginalInput, Source: IntentSource, Kind: string(value.Kind), Interpretation: value.Value, Value: value.Value, Confidence: value.Confidence, Assumptions: value.Assumptions, Plan: p}
	if operation, ok := value.Value.(expression.Operation); ok {
		executed, err := expression.Calculate(operation.OriginalInput)
		if err != nil { return Result{OriginalInput: value.OriginalInput}, err }
		result.Value = executed
		result.Plan.Add(plan.Step{Name: "execute", Operation: string(operation.Operator), Input: operation, Output: executed, Confidence: value.Confidence})
	}
	return result, nil
}

func fromRouter(value router.Result) (Result, error) {
	p := plan.Plan{OriginalInput: value.OriginalInput, Intent: string(value.Kind), Confidence: value.Confidence, Assumptions: append([]string(nil), value.Assumptions...)}
	p.Add(plan.Step{Name: "router", Operation: string(value.Kind), Input: value.OriginalInput, Output: value.Value, Confidence: value.Confidence, Assumptions: value.Assumptions})
	result := Result{OriginalInput: value.OriginalInput, Source: RouterSource, Kind: string(value.Kind), Interpretation: value.Value, Value: value.Value, Confidence: value.Confidence, Assumptions: value.Assumptions, Plan: p}
	if operation, ok := value.Value.(expression.Operation); ok {
		executed, err := expression.Calculate(operation.OriginalInput)
		if err != nil { return Result{OriginalInput: value.OriginalInput}, err }
		result.Value = executed
		result.Plan.Add(plan.Step{Name: "execute", Operation: string(operation.Operator), Input: operation, Output: executed, Confidence: value.Confidence})
	}
	return result, nil
}
