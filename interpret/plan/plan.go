// Package plan defines an explainable execution plan produced from interpreted input.
package plan

import "github.com/codedbyhassan/gokit/interpret"

// Step is one deterministic operation in an execution plan.
type Step struct {
	Name        string
	Operation   string
	Input       any
	Output      any
	Confidence  interpret.Confidence
	Assumptions []string
}

// Plan represents the interpreted intent before or alongside execution.
type Plan struct {
	OriginalInput string
	Intent        string
	Steps         []Step
	Confidence    interpret.Confidence
	Assumptions   []string
}

// Add appends a step and updates plan confidence when the new step is weaker.
func (p *Plan) Add(step Step) {
	p.Steps = append(p.Steps, step)
	if p.Confidence == 0 || step.Confidence < p.Confidence {
		p.Confidence = step.Confidence
	}
	p.Assumptions = append(p.Assumptions, step.Assumptions...)
}
