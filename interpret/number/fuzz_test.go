package number

import "testing"

func FuzzParseNeverPanics(f *testing.F) {
	for _, seed := range []string{"1,000", "two point five million", "negative one hundred", "million five", "one point", "05-06-2020"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		_, _ = Parse(input)
	})
}
