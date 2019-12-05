package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/pulpfree/gsales-xls-reports/model"
)

const (
	timeDayFormat   = "2006-01-02"
	timeMonthFormat = "2006-01"
)

// CreateDates function
func CreateDates(input *model.RequestInput) (dates *model.RequestDates, err error) {

	/* // Set start and end dates for monthly reports
	t, err := time.Parse(timeForm, dateStr)
	if err != nil {
		panic(err)
	}
	currentYear, currentMonth, _ := t.Date()
	currentLocation := t.Location()
	dte := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	suite.dates = &model.RequestDates{
		DateFrom: dte,
		DateTo:   dte.AddDate(0, 1, -1),
	} */

	fmt.Printf("input %+v\n", input)
	// If it's a date, then we assume a month date range
	if input.Date != "" {
		fmt.Printf("input.Date %s\n", input.Date)
		t, err := time.Parse(timeMonthFormat, input.Date)
		if err != nil {
			return nil, err
		}
		currentYear, currentMonth, _ := t.Date()
		currentLocation := t.Location()
		dte := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		dates = &model.RequestDates{
			DateFrom: dte,
			DateTo:   dte.AddDate(0, 1, -1),
		}
		fmt.Printf("t %+v\n", t)
	} else { // else we should have a start and end date
		if input.DateFrom == "" || input.DateTo == "" {
			return nil, errors.New("Invalid dates in util.CreateDates")
		}
	}
	return dates, err
}
