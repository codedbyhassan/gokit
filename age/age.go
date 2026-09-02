// Package age provides reusable age calculations from dates of birth.
package age

import (
	"fmt"
	"time"
)

// Result describes a person's age at a reference date.
type Result struct {
	Years      int
	Months     int
	Days       int
	BirthDate  time.Time
	AsOf       time.Time
}

// Calculate returns the calendar age at asOf. The birth date must not be in the future.
func Calculate(birthDate, asOf time.Time) (Result, error) {
	birthDate = dateOnly(birthDate)
	asOf = dateOnly(asOf)
	if birthDate.After(asOf) {
		return Result{}, fmt.Errorf("birth date %s is after reference date %s", birthDate.Format("2006-01-02"), asOf.Format("2006-01-02"))
	}

	years := asOf.Year() - birthDate.Year()
	anniversary := addYears(birthDate, years)
	if anniversary.After(asOf) {
		years--
		anniversary = addYears(birthDate, years)
	}

	months := 0
	cursor := anniversary
	for {
		next := cursor.AddDate(0, 1, 0)
		if next.After(asOf) {
			break
		}
		cursor = next
		months++
	}
	days := int(asOf.Sub(cursor).Hours() / 24)

	return Result{Years: years, Months: months, Days: days, BirthDate: birthDate, AsOf: asOf}, nil
}

// Years returns completed years between birthDate and asOf.
func Years(birthDate, asOf time.Time) (int, error) {
	result, err := Calculate(birthDate, asOf)
	if err != nil {
		return 0, err
	}
	return result.Years, nil
}

func dateOnly(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
}

func addYears(value time.Time, years int) time.Time {
	candidate := value.AddDate(years, 0, 0)
	if value.Month() == time.February && value.Day() == 29 && candidate.Month() == time.March {
		return time.Date(candidate.Year(), time.February, 28, 0, 0, 0, 0, value.Location())
	}
	return candidate
}
