package report

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/webbtech/gsales-xls-reports/config"
	"github.com/webbtech/gsales-xls-reports/model"
	"github.com/webbtech/gsales-xls-reports/mongodb"
	"github.com/webbtech/gsales-xls-reports/utils"
)

const (
	monthDate      = "2022-03"
	periodDateFrom = "2022-02-19"
	periodDateTo   = "2022-03-05"
	timeForm       = "2006-01-02"
)

const (
	monthlyReport        = "monthlysales"
	bankCardReport       = "bankcards"
	employeeOSReport     = "employeeos"
	fuelSalesReport      = "fuelsales"
	payPeriodReport      = "payperiod"
	productNumbersReport = "productnumbers"
)

// IntegSuite struct
type IntegSuite struct {
	report                  *Report
	bankCardReportReq       *model.ReportRequest
	employeeOSReportReq     *model.ReportRequest
	fuelSalesReportReq      *model.ReportRequest
	monthReportReq          *model.ReportRequest
	payPeriodReportReq      *model.ReportRequest
	productNumbersReportReq *model.ReportRequest
	suite.Suite
}

var (
	cfg              *config.Config
	db               *mongodb.MDB
	startDte, endDte time.Time
	tp               model.ReportType
)

// SetupTest method
func (s *IntegSuite) SetupTest() {
	// init config
	os.Setenv("Stage", "test")
	cfg = &config.Config{}
	err := cfg.Init()
	s.NoError(err)

	// Set db
	db, err = mongodb.NewDB(cfg.GetMongoConnectURL(), cfg.DbName)
	if err != nil {
		log.Fatalf("Failed to initial db with error: %s", err)
	}

	// create shared dates struct
	startDte, endDte, _ = utils.DatesFromDays(periodDateFrom, periodDateTo)
	dayDates := &model.RequestDates{
		DateFrom: startDte,
		DateTo:   endDte,
	}

	startDte, endDte, _ = utils.DatesFromMonth(monthDate)
	monthDates := &model.RequestDates{
		DateFrom: startDte,
		DateTo:   endDte,
	}

	// BankCard report request
	tp, _ = model.ReportStringToType(bankCardReport)
	s.bankCardReportReq = &model.ReportRequest{
		Dates:      dayDates,
		ReportType: tp,
	}
	// EmployeeOS report request
	tp, _ = model.ReportStringToType(employeeOSReport)
	s.employeeOSReportReq = &model.ReportRequest{
		Dates:      dayDates,
		ReportType: tp,
	}

	// FuelSales report request
	tp, _ = model.ReportStringToType(fuelSalesReport)
	s.fuelSalesReportReq = &model.ReportRequest{
		Dates:      monthDates,
		ReportType: tp,
	}

	// Monthly report request
	tp, _ = model.ReportStringToType(monthlyReport)
	s.monthReportReq = &model.ReportRequest{
		Dates:      monthDates,
		ReportType: tp,
	}

	// PayPeriod report request
	tp, _ = model.ReportStringToType(payPeriodReport)
	s.payPeriodReportReq = &model.ReportRequest{
		Dates:      dayDates,
		ReportType: tp,
	}

	// ProducctNumbers report request
	tp, _ = model.ReportStringToType(productNumbersReport)
	s.productNumbersReportReq = &model.ReportRequest{
		Dates:      dayDates,
		ReportType: tp,
	}
}

// TestSetFileName method
func (s *IntegSuite) TestSetFileName() {
	var err error
	s.report, err = New(s.monthReportReq, cfg, db)
	s.NoError(err)
	s.report.setFileName()
	fileNm := s.report.getFileName()

	expectedFileNm := fmt.Sprintf("MonthlySalesReport_%s.xlsx", monthDate)
	s.Equal(expectedFileNm, fileNm)

	// now test report with date range
	s.report, err = New(s.payPeriodReportReq, cfg, db)
	s.NoError(err)
	s.report.setFileName()
	fileNm = s.report.getFileName()

	expectedFileNm = fmt.Sprintf("PayPeriodReport_%s_thru_%s.xlsx", periodDateFrom, periodDateTo)
	s.Equal(expectedFileNm, fileNm)
}

// TestFuelSalesGetRecords method
func (s *IntegSuite) TestFuelSalesGetRecords() {
	var err error

	s.report, err = New(s.fuelSalesReportReq, cfg, db)
	s.NoError(err)
	rep := &FuelSales{
		dates: s.report.dates,
		db:    s.report.db,
	}
	fsRecs, err := rep.GetRecords()
	s.NoError(err)
	s.True(len(fsRecs) > 10)
}

// TestShiftTypeGetRecords method
func (s *IntegSuite) TestShiftTypeGetRecords() {
	var err error

	s.report, err = New(s.bankCardReportReq, cfg, db)
	s.NoError(err)
	bc := &BankCard{
		dates: s.report.dates,
		db:    s.report.db,
	}
	bcRecs, err := bc.GetRecords()
	s.NoError(err)
	s.True(len(bcRecs) > 10)

	s.report, err = New(s.employeeOSReportReq, cfg, db)
	s.NoError(err)
	eo := &EmployeeOS{
		dates: s.report.dates,
		db:    s.report.db,
	}
	esRecs, err := eo.GetRecords()
	s.NoError(err)
	s.True(len(esRecs) > 10)

	s.report, err = New(s.payPeriodReportReq, cfg, db)
	s.NoError(err)
	pp := &PayPeriod{
		dates: s.report.dates,
		db:    s.report.db,
	}
	ppRecs, err := pp.GetRecords()
	s.NoError(err)
	s.True(len(ppRecs) > 10)

	s.report, err = New(s.productNumbersReportReq, cfg, db)
	s.NoError(err)
	pn := &ProductNumbers{
		dates: s.report.dates,
		db:    s.report.db,
	}
	pnRecs, err := pn.GetRecords()
	s.NoError(err)
	s.True(len(pnRecs) > 10)
}

// TestMonthlyGetRecords method
func (s *IntegSuite) TestMonthlyGetRecords() {
	var err error

	s.report, err = New(s.monthReportReq, cfg, db)
	s.NoError(err)

	conf, _ := s.report.db.GetConfig()
	ms := &MonthlySales{
		cfg:   conf,
		dates: s.report.dates,
		db:    s.report.db,
	}

	msRecs, err := ms.GetRecords()
	s.NoError(err)
	s.True(len(msRecs) > 10)
}

// TestCreate method
func (s *IntegSuite) TestCreate() {
	var err error

	s.report, err = New(s.bankCardReportReq, cfg, db)
	err = s.report.create()
	s.NoError(err)

	s.report, err = New(s.employeeOSReportReq, cfg, db)
	err = s.report.create()
	s.NoError(err)

	s.report, err = New(s.fuelSalesReportReq, cfg, db)
	err = s.report.create()
	s.NoError(err)

	s.report, err = New(s.monthReportReq, cfg, db)
	err = s.report.create()
	s.NoError(err)

	s.report, err = New(s.payPeriodReportReq, cfg, db)
	err = s.report.create()
	s.NoError(err)

	s.report, err = New(s.productNumbersReportReq, cfg, db)
	err = s.report.create()
	s.NoError(err)
}

// TestSaveFuelSalesToDisk method
func (s *IntegSuite) TestSaveFuelSalesToDisk() {
	var err error

	s.report, err = New(s.fuelSalesReportReq, cfg, db)
	s.NoError(err)
	_, err = s.report.SaveToDisk("../tmp")
	s.NoError(err)
}

// TestSaveProductNumbersToDisk method
func (s *IntegSuite) TestSaveProductNumbersToDisk() {
	var err error

	s.report, err = New(s.productNumbersReportReq, cfg, db)
	s.NoError(err)
	_, err = s.report.SaveToDisk("../tmp")
	s.NoError(err)
}

// TestSaveEmployeeOSToDisk method
func (s *IntegSuite) TestSaveEmployeeOSToDisk() {
	var err error

	s.report, err = New(s.employeeOSReportReq, cfg, db)
	s.NoError(err)
	_, err = s.report.SaveToDisk("../tmp")
	s.NoError(err)
}

// TestSaveBankCardToDisk method
func (s *IntegSuite) TestSaveBankCardToDisk() {
	var err error

	s.report, err = New(s.bankCardReportReq, cfg, db)
	_, err = s.report.SaveToDisk("../tmp")
	s.NoError(err)
}

// TestSavePayPeriodToDisk method
func (s *IntegSuite) TestSavePayPeriodToDisk() {
	var err error

	s.report, err = New(s.payPeriodReportReq, cfg, db)
	_, err = s.report.SaveToDisk("../tmp")
	s.NoError(err)
}

// TestSaveMonthlyToDisk method
func (s *IntegSuite) TestSaveMonthlyToDisk() {
	var err error

	s.report, err = New(s.monthReportReq, cfg, db)
	_, err = s.report.SaveToDisk("../tmp")
	s.NoError(err)
}

// TestCreateSignedURL method
func (s *IntegSuite) TestCreateSignedURL() {
	var err error
	s.report, err = New(s.bankCardReportReq, cfg, db)
	url, err := s.report.CreateSignedURL()
	s.NoError(err)
	s.True(len(url) > 100)

	s.report, err = New(s.monthReportReq, cfg, db)
	url, err = s.report.CreateSignedURL()
	s.NoError(err)
	s.True(len(url) > 100)

	s.report, err = New(s.payPeriodReportReq, cfg, db)
	url, err = s.report.CreateSignedURL()
	s.NoError(err)
	s.True(len(url) > 100)
}

// TestIntegrationSuite function
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}
