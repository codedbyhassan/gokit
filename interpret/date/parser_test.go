package date

import (
	"errors"
	"testing"
	"time"

	"github.com/codedbyhassan/gokit/interpret"
)

func TestParseFlexibleDates(t *testing.T) {
	tests := []struct {
		name string
		input string
		want time.Time
	}{
		{"day month year", "11-11-2011", time.Date(2011, 11, 11, 0, 0, 0, 0, time.UTC)},
		{"year month day", "2011-11-11", time.Date(2011, 11, 11, 0, 0, 0, 0, time.UTC)},
		{"year in middle", "11-2011-11", time.Date(2011, 11, 11, 0, 0, 0, 0, time.UTC)},
		{"named month", "11 November 2011", time.Date(2011, 11, 11, 0, 0, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if !got.Value.Equal(tt.want) {
				t.Fatalf("Parse() = %v, want %v", got.Value, tt.want)
			}
			if got.Confidence <= 0 {
				t.Fatal("expected a positive confidence score")
			}
		})
	}
}

func TestParseAmbiguousDate(t *testing.T) {
	_, err := Parse("05-06-2020")
	if !errors.Is(err, interpret.ErrAmbiguousInput) {
		t.Fatalf("expected ErrAmbiguousInput, got %v", err)
	}
}

func TestParseInvalidDate(t *testing.T) {
	_, err := Parse("31-02-2020")
	if !errors.Is(err, interpret.ErrUnrecognizedInput) {
		t.Fatalf("expected ErrUnrecognizedInput, got %v", err)
	}
}
