package util

import (
	"errors"
	"time"
)

type DateRange struct {
	From time.Time
	To   time.Time
}

func isValidTime(time time.Time) bool {
	return time.IsZero()
}

func ValidateTimeRange(from, to time.Time) error {
	if !isValidTime(from) {
		return errors.New("invalid [From]")
	}
	if !isValidTime(to) {
		return errors.New("invalid [To]")
	}
	if from.After(to) {
		return errors.New("invalid range [From] is after [To]")
	}
	return nil
}

func generateDailyRanges(from, to time.Time) []DateRange {
	var ranges []DateRange
	currentDate := from
	for !currentDate.After(to) {
		ranges = append(ranges, DateRange{From: currentDate, To: currentDate})
		currentDate = currentDate.AddDate(0, 0, 1)
	}
	return ranges
}

func generateWeeklyRanges(from, to time.Time) []DateRange {
	var ranges []DateRange
	currentStartOfWeek := from

	for !currentStartOfWeek.After(to) {
		currentEndOfWeek := currentStartOfWeek.AddDate(0, 0, 6)
		if currentEndOfWeek.After(to) {
			currentEndOfWeek = to
		}
		ranges = append(ranges, DateRange{From: currentStartOfWeek, To: currentEndOfWeek})
		currentStartOfWeek = currentEndOfWeek.AddDate(0, 0, 1)
	}
	return ranges
}

func GenerateDateRanges(from, to time.Time) []DateRange {
	diff := to.Sub(from).Hours() / 24
	if diff > 30 {
		return generateWeeklyRanges(from, to)
	}
	return generateDailyRanges(from, to)

}

func IsDateInRange(date time.Time, dateRange DateRange) bool {
	return (date.Equal(dateRange.From) || date.After(dateRange.From)) && (date.Equal(dateRange.To) || date.Before(dateRange.To))

}

func CalculatePreviousPeriod(from, to time.Time) (time.Time, time.Time) {
	duration := to.Sub(from)
	previousEndDate := from.Add(-time.Nanosecond)
	previousStartDate := previousEndDate.Add(-duration)

	return previousStartDate, previousEndDate
}
