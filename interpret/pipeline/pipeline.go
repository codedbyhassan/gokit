// Package pipeline provides GoKit's unified human-input execution boundary.
package pipeline

import (
	"errors"
	"strings"
	"time"

	"github.com/codedbyhassan/gokit/interpret"
	"github.com/codedbyhassan/gokit/interpret/intent"
	"github.com/codedbyhassan/gokit/interpret/plan"
	"github.com/codedbyhassan/gokit/interpret/router"
)

// Source identifies the layer that successfully interpreted the input.
type Source string

const (
	IntentSource Source = "intent"
	RouterSource Source = "router"
)

// Result is the unified output of interpretation and execution.
type Result struct {
	OriginalInput string
	Source        Source
	Kind          string
	Value         any
	Confidence    interpret.Confidence
	Assumptions   []string
	Plan          plan.Plan
}

// Parse interprets and executes input using the current time for contextual intents.
func Parse(input string) (Result, error) { return ParseAt(input, time.Now()) }

// ParseAt is deterministic for callers that need a fixed reference time.
func ParseAt(input string, asOf time.Time) (Result, error) {
	if strings.TrimSpace(input) == "" {
		return Result{OriginalInput: input}, interpret.ErrEmptyInput
	}

	intentResult, err := intent.ParseAt(input, asOf)
	if err == nil {
		return fromIntent(intentResult), nil
	}
	if !errors.Is(err, interpret.ErrUnrecognizedInput) {
		return Result{OriginalInput: input}, err
	}

	routed, err := router.Parse(input)
	if err != nil {
		return Result{OriginalInput: input}, err
	}
	return fromRouter(routed), nil
}

func fromIntent(value intent.Result) Result {
	p := plan.Plan{OriginalInput: value.OriginalInput, Intent: string(value.Kind), Confidence: value.Confidence, Assumptions: append([]string(nil), value.Assumptions...)}
	p.Add(plan.Step{Name: "intent", Operation: string(value.Kind), Input: value.OriginalInput, Output: value.Value, Confidence: value.Confidence, Assumptions: value.Assumptions})
	return Result{OriginalInput: value.OriginalInput, Source: IntentSource, Kind: string(value.Kind), Value: value.Value, Confidence: value.Confidence, Assumptions: value.Assumptions, Plan: p}
}

func fromRouter(value router.Result) Result {
	p := plan.Plan{OriginalInput: value.OriginalInput, Intent: string(value.Kind), Confidence: value.Confidence, Assumptions: append([]string(nil), value.Assumptions...)}
	p.Add(plan.Step{Name: "router", Operation: string(value.Kind), Input: value.OriginalInput, Output: value.Value, Confidence: value.Confidence, Assumptions: value.Assumptions})
	return Result{OriginalInput: value.OriginalInput, Source: RouterSource, Kind: string(value.Kind), Value: value.Value, Confidence: value.Confidence, Assumptions: value.Assumptions, Plan: p}
}
