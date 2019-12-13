package report

import (
	"fmt"
	"os"
	"testing"

	"github.com/pulpfree/gsales-xls-reports/config"
	"github.com/pulpfree/gsales-xls-reports/model"
	"github.com/pulpfree/gsales-xls-reports/validate"
	"github.com/stretchr/testify/suite"
)

const (
	monthDate       = "2019-10"
	monthlyReport   = "monthlysales"
	payPeriodReport = "payperiod"
	periodDateFrom  = "2018-09-01"
	periodDateTo    = "2018-10-10"

	// dateMonth    = "2019-08"
	// dateDayStart = "2019-08-01"
	// dateDayEnd   = "2019-08-16"
	defaultsFP = "../config/defaults.yml"
	timeForm   = "2006-01-02"
)

// IntegSuite struct
type IntegSuite struct {
	report          *Report
	periodReportReq *model.ReportRequest
	monthReportReq  *model.ReportRequest
	suite.Suite
}

var (
	cfg *config.Config
	// dates         *model.RequestDates
	// reportRequest *model.ReportRequest
)

// SetupTest method
func (s *IntegSuite) SetupTest() {
	// init config
	os.Setenv("Stage", "test")
	cfg = &config.Config{DefaultsFilePath: defaultsFP}
	err := cfg.Load()
	s.NoError(err)

	monthReportInput := &model.RequestInput{
		Date:       monthDate,
		ReportType: monthlyReport,
	}
	s.monthReportReq, _ = validate.SetRequest(monthReportInput)

	periodReportInput := &model.RequestInput{
		DateFrom:   periodDateFrom,
		DateTo:     periodDateTo,
		ReportType: payPeriodReport,
	}
	s.periodReportReq, _ = validate.SetRequest(periodReportInput)
}

// TestsetFileName method
func (s *IntegSuite) TestsetFileName() {
	fmt.Printf("cfg in TestsetFileNam %+v\n", cfg)
	var err error
	s.report = New(s.monthReportReq, cfg)
	s.NoError(err)
	s.report.setFileName()
	fileNm := s.report.getFileName()

	expectedFileNm := fmt.Sprintf("MonthlySalesReport_%s.xlsx", monthDate)
	s.Equal(fileNm, expectedFileNm)

	// now test report with date range
	s.report = New(s.periodReportReq, cfg)
	s.NoError(err)
	s.report.setFileName()
	fileNm = s.report.getFileName()

	expectedFileNm = fmt.Sprintf("PayPeriodReport_%s_thru_%s.xlsx", periodDateFrom, periodDateTo)
	s.Equal(fileNm, expectedFileNm)
}

// Testcreate method
func (s *IntegSuite) Testcreate() {
	var err error
	s.report = New(s.monthReportReq, cfg)
	err = s.report.create()
	s.NoError(err)

	s.report = New(s.periodReportReq, cfg)
	err = s.report.create()
	s.NoError(err)
}

// TestSaveToDisk method
func (s *IntegSuite) TestSaveToDisk() {
	var err error
	s.report = New(s.periodReportReq, cfg)
	_, err = s.report.SaveToDisk("../tmp")
	s.NoError(err)
}

// TestCreateSignedURL method
func (s *IntegSuite) TestCreateSignedURL() {
	var err error
	s.report = New(s.periodReportReq, cfg)
	url, err := s.report.CreateSignedURL()
	s.NoError(err)
	s.True(len(url) > 100)
}

// TestIntegrationSuite function
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}
