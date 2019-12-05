package util

import (
	"fmt"
	"testing"

	"github.com/pulpfree/gsales-xls-reports/model"
	"gotest.tools/assert"
)

// TestCreateDates function
func TestCreateDates(t *testing.T) {
	var err error
	monthDate := "2019-08"
	// expectedMonthEndDate := "2019-08-31"
	input := &model.RequestInput{
		Date:     monthDate,
		DateFrom: "",
		DateTo:   "",
	}
	dates, err := CreateDates(input)
	assert.NilError(t, err)
	assert.Equal(t, dates.DateFrom.Format(timeMonthFormat), monthDate)
	// assert.Equal(t, dates.DateTo.Format(timeMonthFormat), expectedMonthEndDate)
	fmt.Printf("dates %+v\n", dates)
}
