package payperiod

import (
	"os"
	"testing"
	"time"

	"github.com/pulpfree/gsales-xls-reports/config"
	"github.com/pulpfree/gsales-xls-reports/model"
	"github.com/stretchr/testify/suite"
)

const (
	cfgHST = int32(13)
	// dateStr = "2019-08-01"
	dateStr    = "2019-09-01"
	defaultsFP = "../../config/defaults.yml"
	timeForm   = "2006-01-02"
)

// IntegSuite struct
type IntegSuite struct {
	pp *Records
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
	suite.pp, err = Init(dates, cfg)
	suite.IsType(suite.pp.cfg, &model.Config{})
}

// ===================== Exported Functions ============================================ //

/* func (suite *IntegSuite) TestPPInit() {
pp := Init(db *MDB, dates *model.RequestDates) (*PayPeriod, error) {
} */

// ===================== Un-exported Functions ================================================ //

// TestsetRecords method
func (suite *IntegSuite) TestsetRecords() {
	var err error
	err = suite.pp.setRecords()
	suite.NoError(err)
	suite.True(len(suite.pp.Records) > 0)
}

// TestsetNonFuelCommission
func (suite *IntegSuite) TestsetNonFuelCommission() {
	var err error
	err = suite.pp.setRecords()
	err = suite.pp.setNonFuelCommission()
	suite.NoError(err)
	suite.True(suite.pp.Records[0].Commission != nil)
}

// TestsetCarWashes method
func (suite *IntegSuite) TestsetCarWashes() {
	var err error
	err = suite.pp.setRecords()
	err = suite.pp.setCarWashes()
	suite.NoError(err)
}

// TestGetRecords method
func (suite *IntegSuite) TestGetRecords() {
	var err error
	records, err := suite.pp.GetRecords()
	suite.NoError(err)
	suite.True(len(records) > 0)
	// fmt.Printf("records[0] %+v\n", records[0])
}

// TestIntegrationSuite function
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}
