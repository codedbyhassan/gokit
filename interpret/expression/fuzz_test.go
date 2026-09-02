package expression

import "testing"

func FuzzCalculateNeverPanics(f *testing.F) {
	for _, seed := range []string{"1 + 2 * 3", "(10 - 4) / 2", "25% of 800", "20% increase on 500", "1 + (2 * 3)"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		_, _ = Calculate(input)
	})
}
