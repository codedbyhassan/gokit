package number_test

import (
	"fmt"

	"github.com/codedbyhassan/gokit/interpret/number"
)

func ExampleParse_compound() {
	result, err := number.Parse("one thousand two hundred and thirty four")
	if err != nil { return }
	fmt.Println(result.Value)
	// Output: 1234
}
