package monthlysales

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
	ms *Sales
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
	suite.ms, err = Init(dates, cfg)
	suite.IsType(suite.ms.cfg, &model.Config{})
}

// ===================== Exported Functions =================================================== //

// TestGetRecords
func (suite *IntegSuite) TestGetRecords() {
	var err error
	records, err := suite.ms.GetRecords()
	suite.NoError(err)
	suite.True(len(records) > 10)
}

// ===================== Un-exported Functions ================================================ //

// TestsetRecords method
func (suite *IntegSuite) TestsetSales() {
	var err error
	err = suite.ms.setRecords()
	suite.NoError(err)
}

// TestsetMonthlySalesProducts method
func (suite *IntegSuite) TestsetMonthlySalesProducts() {
	var err error
	err = suite.ms.setRecords()
	suite.NoError(err)
	err = suite.ms.setMonthlySalesProducts()
	suite.NoError(err)
}

// TestIntegrationSuite function
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}
