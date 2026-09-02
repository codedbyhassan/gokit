package age

import (
	"testing"
	"time"
)

func TestCalculate(t *testing.T) {
	location := time.UTC
	birth := time.Date(2000, 5, 10, 0, 0, 0, 0, location)
	asOf := time.Date(2026, 9, 2, 0, 0, 0, 0, location)

	result, err := Calculate(birth, asOf)
	if err != nil {
		t.Fatal(err)
	}
	if result.Years != 26 || result.Months != 3 || result.Days != 23 {
		t.Fatalf("unexpected age: %+v", result)
	}
}

func TestCalculateRejectsFutureBirthDate(t *testing.T) {
	birth := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	asOf := time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC)
	if _, err := Calculate(birth, asOf); err == nil {
		t.Fatal("expected future birth date to fail")
	}
}
