package util

import (
	"errors"
	"time"

	"github.com/pulpfree/gsales-xls-reports/model"
	"github.com/pulpfree/gsales-xls-reports/pkgerrors"
)

const (
	timeDayFormat   = "2006-01-02"
	timeMonthFormat = "2006-01"
)

// CreateDates function
func CreateDates(input *model.RequestInput) (*model.RequestDates, error) {

	var dates *model.RequestDates
	// If it's a date, then we assume a month date range
	if input.Date != "" {
		t, err := time.Parse(timeMonthFormat, input.Date)
		if err != nil {
			return nil, &pkgerrors.StdError{Err: err.Error(), Caller: "util.CreateDates", Msg: "Error parsing time input.Date in CreateDates"}
		}
		currentYear, currentMonth, _ := t.Date()
		currentLocation := t.Location()
		dte := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		dates = &model.RequestDates{
			DateFrom: dte,
			DateTo:   dte.AddDate(0, 1, -1),
		}

		// else we should have a start and end date
	} else {
		if input.DateFrom == "" || input.DateTo == "" {
			return nil, errors.New("Invalid dates in util.CreateDates")
		}
		tStart, err := time.Parse(timeDayFormat, input.DateFrom)
		if err != nil {
			return nil, err
		}
		tEnd, err := time.Parse(timeDayFormat, input.DateTo)
		if err != nil {
			return nil, err
		}
		dates = &model.RequestDates{
			DateFrom: tStart,
			DateTo:   tEnd,
		}
	}
	return dates, nil
}
