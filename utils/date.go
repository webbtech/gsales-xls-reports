package utils

import (
	"time"
)

const (
	timeDayFormat   = "2006-01-02"
	timeMonthFormat = "2006-01"
)

// DatesFromMonth creates a start and ending date for a specified month
func DatesFromMonth(monthDte string) (start, end time.Time, err error) {

	t, err := time.Parse(timeMonthFormat, monthDte)
	if err != nil {
		return start, end, err
	}
	currentYear, currentMonth, _ := t.Date()
	currentLocation := t.Location()

	start = time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	end = start.AddDate(0, 1, -1)

	return start, end, err
}

// DatesFromDays creates a start and ending date from specified start and end dates
func DatesFromDays(startStr, endStr string) (start, end time.Time, err error) {

	start, err = time.Parse(timeDayFormat, startStr)
	if err != nil {
		return start, end, err
	}
	end, err = time.Parse(timeDayFormat, endStr)
	if err != nil {
		return start, end, err
	}

	return start, end, err
}
