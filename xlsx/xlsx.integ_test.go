package xlsx

import (
	"os"
	"testing"
	"time"

	"github.com/pulpfree/gsales-xls-reports/config"
	"github.com/pulpfree/gsales-xls-reports/model"
	"github.com/pulpfree/gsales-xls-reports/model/bankcards"
	"github.com/pulpfree/gsales-xls-reports/model/monthlysales"
	"github.com/pulpfree/gsales-xls-reports/model/payperiod"
	"github.com/stretchr/testify/suite"
)

const (
	// dateStr    = "2019-08-01"
	dateStr           = "2019-09-01"
	defaultsFP        = "../config/defaults.yml"
	filePathBankCards = "../tmp/bankCards.xlsx"
	filePathMonthly   = "../tmp/monthlySales.xlsx"
	filePathPay       = "../tmp/payPeriod.xlsx"
	timeForm          = "2006-01-02"
)

// IntegSuite struct
type IntegSuite struct {
	cfg          *config.Config
	db           model.DBHandler
	file         *XLSX
	bankCards    *bankcards.Records
	monthlysales *monthlysales.Sales
	payPeriod    *payperiod.Records
	suite.Suite
}

// SetupTest method
func (suite *IntegSuite) SetupTest() {
	// setup config
	os.Setenv("Stage", "test")
	cfg := &config.Config{DefaultsFilePath: defaultsFP}
	err := cfg.Load()
	suite.NoError(err)

	suite.file, err = NewFile()
	suite.NoError(err)
	suite.IsType(new(XLSX), suite.file)

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

	suite.bankCards, err = bankcards.Init(dates, cfg)
	suite.NoError(err)

	suite.monthlysales, err = monthlysales.Init(dates, cfg)
	suite.NoError(err)

	suite.payPeriod, err = payperiod.Init(dates, cfg)
	suite.NoError(err)
}

// TestBankCards method
func (suite *IntegSuite) TestBankCards() {
	records, err := suite.bankCards.GetRecords()
	suite.NoError(err)

	err = suite.file.BankCards(records)
	suite.NoError(err)

	_, err = suite.file.OutputToDisk(filePathBankCards)
	suite.NoError(err)
}

// TestMonthlySales method
func (suite *IntegSuite) TestMonthlySales() {
	records, err := suite.monthlysales.GetRecords()
	suite.NoError(err)

	err = suite.file.MonthlySales(records)
	suite.NoError(err)

	_, err = suite.file.OutputToDisk(filePathMonthly)
	suite.NoError(err)

	// to open, use: open -a Numbers ./testfile.xlsx
}

// TestPayPeriod method
func (suite *IntegSuite) TestPayPeriod() {
	records, err := suite.payPeriod.GetRecords()
	suite.NoError(err)
	// fmt.Printf("records[0] %+v\n", records[0])

	err = suite.file.PayPeriod(records)
	suite.NoError(err)

	_, err = suite.file.OutputToDisk(filePathPay)
	suite.NoError(err)
}

// TestIntegrationSuite function
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}
