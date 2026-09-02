package unit_test

import (
	"fmt"

	"github.com/codedbyhassan/gokit/interpret/unit"
)

func ExampleParse() {
	result, err := unit.Parse("5kg to pounds")
	if err != nil { return }
	conversion := result.Value.(unit.Conversion)
	fmt.Printf("%.2f %s\n", conversion.Value, conversion.To.Symbol)
	// Output: 11.02 lb
}
