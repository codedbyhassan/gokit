package pipeline_test

import (
	"fmt"
	"time"

	"github.com/codedbyhassan/gokit/interpret/pipeline"
)

func ExampleParse() {
	result, err := pipeline.Parse("what is 20% of 500")
	if err != nil {
		return
	}
	fmt.Println(result.Kind)
	fmt.Println(result.Value)
	// Output:
	// calculate
	// {500 20 * 100}
}

func ExampleParseAt() {
	asOf := time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC)
	result, err := pipeline.ParseAt("how old is someone born 11-11-2011", asOf)
	if err != nil {
		return
	}
	fmt.Println(result.Kind)
	fmt.Println(result.Plan.Intent)
	// Output:
	// age
	// age
}
