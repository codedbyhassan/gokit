package plan

import (
	"testing"

	"github.com/codedbyhassan/gokit/interpret"
)

func TestPlanAdd(t *testing.T) {
	var p Plan
	p.Add(Step{Name: "parse", Operation: "parse number", Confidence: interpret.HighConfidence})
	p.Add(Step{Name: "calculate", Operation: "add", Confidence: interpret.MediumConfidence, Assumptions: []string{"normalized units"}})

	if len(p.Steps) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(p.Steps))
	}
	if p.Confidence != interpret.MediumConfidence {
		t.Fatalf("expected weakest confidence, got %v", p.Confidence)
	}
	if len(p.Assumptions) != 1 || p.Assumptions[0] != "normalized units" {
		t.Fatalf("unexpected assumptions: %#v", p.Assumptions)
	}
}
