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

// SetupTest method
func (suite *IntegSuite) SetupTest() {
	// init config
	os.Setenv("Stage", "test")
	cfg := &config.Config{DefaultsFilePath: defaultsFP}
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
	dates := &model.RequestDates{
		DateFrom: dte,
		DateTo:   dte.AddDate(0, 1, -1),
	}

	suite.report, err = New(dates, cfg, "monthlysales")
}

// TestsetFileName method
func (suite *IntegSuite) TestsetFileName() {
	suite.report.setFileName()
	fileNm := suite.report.getFileName()

	date := dateStr[:len(dateStr)-3]
	expectedFileNm := fmt.Sprintf("%s_%s.xlsx", reportFileName, date)
	suite.Equal(fileNm, expectedFileNm)
}

// TestCreate method
func (suite *IntegSuite) TestCreate() {
	var err error
	err = suite.report.Create()
	suite.NoError(err)
}

// TestSaveToDisk method
func (suite *IntegSuite) TestSaveToDisk() {
	suite.report.haveFile = false
	var err error
	_, err = suite.report.SaveToDisk("../tmp")
	suite.NoError(err)
}

// TestCreateSignedURL method
func (suite *IntegSuite) TestCreateSignedURL() {
	suite.report.haveFile = false
	var err error
	url, err := suite.report.CreateSignedURL()
	suite.NoError(err)
	suite.True(len(url) > 100)
}

// TestIntegrationSuite function
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}
