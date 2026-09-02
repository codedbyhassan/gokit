package pipeline

import "testing"

func BenchmarkParseExpression(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := Parse("what is 20% of 500"); err != nil { b.Fatal(err) }
	}
}

func BenchmarkParseNumber(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := Parse("two point five million"); err != nil { b.Fatal(err) }
	}
}
