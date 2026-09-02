package number_test

import (
	"fmt"

	"github.com/codedbyhassan/gokit/interpret/number"
)

func ExampleParse() {
	result, err := number.Parse("2.5 million")
	if err != nil {
		return
	}
	fmt.Println(result.Value)
	// Output: 2.5e+06
}
