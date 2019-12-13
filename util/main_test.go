package util

import (
	"errors"
	"testing"

	"github.com/pulpfree/gsales-xls-reports/model"
	"github.com/pulpfree/gsales-xls-reports/pkgerrors"
	"gotest.tools/assert"
)

// TestCreateDatesMonth function
func TestCreateDatesMonth(t *testing.T) {
	var err error
	monthDate := "2019-08"
	expectedMonthStartDate := "2019-08-01"
	expectedMonthEndDate := "2019-08-31"
	input := &model.RequestInput{
		Date:     monthDate,
		DateFrom: "",
		DateTo:   "",
	}
	dates, err := CreateDates(input)
	assert.NilError(t, err)
	assert.Equal(t, dates.DateFrom.Format(timeDayFormat), expectedMonthStartDate)
	assert.Equal(t, dates.DateTo.Format(timeDayFormat), expectedMonthEndDate)
}

// TestCreateDatesStartEnd function
func TestCreateDatesStartEnd(t *testing.T) {
	var err error
	startDate := "2019-08-01"
	endDate := "2019-08-16"
	input := &model.RequestInput{
		Date:     "",
		DateFrom: startDate,
		DateTo:   endDate,
	}
	dates, err := CreateDates(input)
	assert.NilError(t, err)
	assert.Equal(t, dates.DateFrom.Format(timeDayFormat), startDate)
	assert.Equal(t, dates.DateTo.Format(timeDayFormat), endDate)
}

// TestCreateDatesErrors function
func TestCreateDatesErrors(t *testing.T) {
	var err error
	expectedErrorStr := "Invalid dates in util.CreateDates"
	input := &model.RequestInput{}
	_, err = CreateDates(input)
	assert.Error(t, err, expectedErrorStr)

	expectedErrorStr = "Error parsing time input.Date in CreateDates"
	input = &model.RequestInput{
		Date: "20",
	}
	_, err = CreateDates(input)

	var err2 *pkgerrors.StdError
	if ok := errors.As(err, &err2); ok {
		// 	// handle gracefully
		// fmt.Printf("err2.Info %+v\n", err2.Msg)
		// fmt.Printf("err2.Err %+v\n", err2.Err)
		assert.Equal(t, err2.Msg, expectedErrorStr)
	}

	if errors.Is(err, err2) {
		// fmt.Printf("err in Is %+v\n", err)
		assert.Equal(t, err2.Msg, expectedErrorStr)
	}
}
