package report

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/pulpfree/gsales-xls-reports/config"
	"github.com/pulpfree/gsales-xls-reports/model"
	"github.com/stretchr/testify/suite"
)

const (
	// dateStr = "2019-08-01"
	dateStr    = "2019-09-01"
	defaultsFP = "../config/defaults.yml"
	timeForm   = "2006-01-02"
)

// IntegSuite struct
type IntegSuite struct {
	report *Report
	suite.Suite
}

var (
	cfg   *config.Config
	dates *model.RequestDates
)

// SetupTest method
func (suite *IntegSuite) SetupTest() {
	// init config
	os.Setenv("Stage", "test")
	cfg = &config.Config{DefaultsFilePath: defaultsFP}
	err := cfg.Load()
	suite.NoError(err)

	// Set start and end dates for monthly reports
	t, err := time.Parse(timeForm, dateStr)
	if err != nil {
		panic(err)
	}
	currentYear, currentMonth, _ := t.Date()
	currentLocation := t.Location()

	dte := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	dates = &model.RequestDates{
		DateFrom: dte,
		DateTo:   dte.AddDate(0, 1, -1),
	}

}

// TestsetFileName method
func (suite *IntegSuite) TestsetFileName() {
	var err error
	suite.report, err = New(dates, cfg, "monthlysales")
	suite.NoError(err)
	suite.report.setFileName()
	fileNm := suite.report.getFileName()

	// stDteStr := dates.DateFrom.Format(timeFormat)
	startDate := dateStr[:len(dateStr)-3]
	expectedFileNm := fmt.Sprintf("MonthlySalesReport_%s.xlsx", startDate)
	suite.Equal(fileNm, expectedFileNm)

	// now test report with date range
	suite.report, err = New(dates, cfg, "payperiod")
	suite.NoError(err)
	suite.report.setFileName()
	fileNm = suite.report.getFileName()

	startDate = dates.DateFrom.Format(timeFormatLong)
	endDate := dates.DateTo.Format(timeFormatLong)
	expectedFileNm = fmt.Sprintf("PayPeriodReport_%s_thru_%s.xlsx", startDate, endDate)
	suite.Equal(fileNm, expectedFileNm)
}

// TestallowedTypes method
func (suite *IntegSuite) TestallowedTypes() {
	suite.report = &Report{}
	var err error
	validType := "monthlysales"
	invalidType := "monthlyreport"

	_, err = suite.report.testReportType(validType)
	suite.NoError(err)

	_, err = suite.report.testReportType(invalidType)
	suite.Error(err)
}

// TestCreate method
func (suite *IntegSuite) TestCreate() {
	var err error
	suite.report, err = New(dates, cfg, "monthlysales")
	err = suite.report.create()
	suite.NoError(err)

	suite.report, err = New(dates, cfg, "payperiod")
	err = suite.report.create()
	suite.NoError(err)
}

// TestSaveToDisk method
func (suite *IntegSuite) TestSaveToDisk() {
	var err error
	suite.report, err = New(dates, cfg, "payperiod")
	_, err = suite.report.SaveToDisk("../tmp")
	suite.NoError(err)
}

// TestCreateSignedURL method
func (suite *IntegSuite) TestCreateSignedURL() {
	var err error
	suite.report, err = New(dates, cfg, "payperiod")
	// suite.report, err = New(dates, cfg, "monthlysales")
	url, err := suite.report.CreateSignedURL()
	suite.NoError(err)
	suite.True(len(url) > 100)
}

// TestIntegrationSuite function
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}
