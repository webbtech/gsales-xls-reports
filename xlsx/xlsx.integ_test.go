package xlsx

import (
	"os"
	"testing"

	"github.com/webbtech/gsales-xls-reports/config"
	"github.com/webbtech/gsales-xls-reports/model"

	"github.com/stretchr/testify/suite"
)

const (
	monthDate      = "2019-10"
	periodDateFrom = "2019-09-01"
	periodDateTo   = "2019-10-10"
	timeForm       = "2006-01-02"
)

const (
	monthlyReport    = "monthlysales"
	bankCardReport   = "bankcards"
	employeeOSReport = "employeeos"
	payPeriodReport  = "payperiod"
)

const (
	// dateStr    = "2019-08-01"
	// dateStr           = "2019-11-01"
	defaultsFP        = "../config/defaults.yml"
	filePathBankCards = "../tmp/bankCards.xlsx"
	filePathMonthly   = "../tmp/monthlySales.xlsx"
	filePathPay       = "../tmp/payPeriod.xlsx"
)

// IntegSuite struct
type IntegSuite struct {
	cfg  *config.Config
	db   model.DBHandler
	file *XLSX
	// bankCardReport *report.BankCard
	// employeeOS    *bankcards.Records
	// monthlysales *monthlysales.Sales
	// payPeriod    *payperiod.Records
	suite.Suite
}

// SetupTest method
func (s *IntegSuite) SetupTest() {
	// setup config
	os.Setenv("Stage", "test")
	cfg := &config.Config{IsDefaultsLocal: true}
	err := cfg.Init()
	s.NoError(err)

	s.file, err = NewFile()
	s.NoError(err)
	s.IsType(new(XLSX), s.file)

	/* bankCardReportInput := &model.RequestInput{
		DateFrom:   periodDateFrom,
		DateTo:     periodDateTo,
		ReportType: bankCardReport,
	}
	bankCardReportReq, _ := validate.SetRequest(bankCardReportInput)
	report, err := report.New(bankCardReportReq, cfg)
	s.NoError(err)
	fmt.Printf("report %+v\n", report) */
	/* bankCardReport := &report.BankCard{
		dates: report.dates,
		db:    report.db,
	} */

	/* suite.bankCards, err = bankcards.Init(dates, cfg)
	suite.NoError(err)

	suite.monthlysales, err = monthlysales.Init(dates, cfg)
	suite.NoError(err)

	suite.payPeriod, err = payperiod.Init(dates, cfg)
	suite.NoError(err) */
}

// TestBankCards method
/* func (s *IntegSuite) TestBankCards() {
	records, err := s.bankCards.GetRecords()
	s.NoError(err)

	err = s.file.BankCards(records)
	s.NoError(err)

	_, err = s.file.OutputToDisk(filePathBankCards)
	s.NoError(err)
} */

// TestEmployeeOS method
/* func (s *IntegSuite) TestEmployeeOS() {
	records, err := s.payPeriod.GetRecords()
	s.NoError(err)
	// fmt.Printf("records[0] %+v\n", records[0])

	err = s.file.PayPeriod(records)
	s.NoError(err)

	_, err = s.file.OutputToDisk(filePathPay)
	s.NoError(err)
} */

// TestMonthlySales method
/* func (s *IntegSuite) TestMonthlySales() {
	records, err := s.monthlysales.GetRecords()
	s.NoError(err)

	err = s.file.MonthlySales(records)
	s.NoError(err)

	_, err = s.file.OutputToDisk(filePathMonthly)
	s.NoError(err)

	// to open, use: open -a Numbers ./testfile.xlsx
} */

// TestPayPeriod method
/* func (s *IntegSuite) TestPayPeriod() {
	records, err := s.payPeriod.GetRecords()
	s.NoError(err)
	// fmt.Printf("records[0] %+v\n", records[0])

	err = s.file.PayPeriod(records)
	s.NoError(err)

	_, err = s.file.OutputToDisk(filePathPay)
	s.NoError(err)
} */

// TestIntegrationSuite function
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}
