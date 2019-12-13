package validate

import (
	"testing"
	"time"

	"github.com/pulpfree/gsales-xls-reports/model"
	"github.com/stretchr/testify/suite"
)

const (
	monthDate       = "2019-10"
	monthlyReport   = "monthlysales"
	payPeriodReport = "payperiod"
	periodDateFrom  = "2018-09-01"
	periodDateTo    = "2018-10-10"
)

// UnitSuite struct
type UnitSuite struct {
	suite.Suite
	requestPeriodReport *model.RequestInput
	requestMonthReport  *model.RequestInput
	requestVars         *model.ReportRequest
}

// SetupTest method
func (s *UnitSuite) SetupTest() {
	s.requestMonthReport = &model.RequestInput{
		Date:       monthDate,
		ReportType: monthlyReport,
	}

	s.requestPeriodReport = &model.RequestInput{
		DateFrom:   periodDateFrom,
		DateTo:     periodDateTo,
		ReportType: payPeriodReport,
	}
}

// TestSetMonthlyRequest method
func (s *UnitSuite) TestSetMonthlyRequest() {

	var rt *model.ReportType
	var t time.Time

	req, err := SetRequest(s.requestMonthReport)
	s.NoError(err)
	s.IsType(&model.RequestDates{}, req.Dates)
	s.IsType(&model.ReportRequest{}, req)
	s.Equal(int(model.MonthlySalesReport), int(*req.ReportType))
	s.IsType(rt, req.ReportType)
	s.IsType(t, req.Dates.DateFrom)
	s.IsType(t, req.Dates.DateTo)
}

// TestFailTypeSetRequest method
func (s *UnitSuite) TestFailTypeSetRequest() {
	requestFail := &model.RequestInput{
		DateFrom:   periodDateFrom,
		DateTo:     periodDateTo,
		ReportType: "falsetype",
	}
	_, err := SetRequest(requestFail)
	s.Error(err)
}

// TestUnitSuite function
func TestUnitSuite(t *testing.T) {
	suite.Run(t, new(UnitSuite))
}
