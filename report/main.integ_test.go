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
	monthDate      = "2019-10"
	periodDateFrom = "2019-11-01"
	periodDateTo   = "2019-11-30"
	defaultsFP     = "../config/defaults.yml"
	timeForm       = "2006-01-02"
)

const (
	monthlyReport        = "monthlysales"
	bankCardReport       = "bankcards"
	employeeOSReport     = "employeeos"
	payPeriodReport      = "payperiod"
	productNumbersReport = "productnumbers"
)

// IntegSuite struct
type IntegSuite struct {
	report                  *Report
	bankCardReportReq       *model.ReportRequest
	employeeOSReportReq     *model.ReportRequest
	monthReportReq          *model.ReportRequest
	payPeriodReportReq      *model.ReportRequest
	productNumbersReportReq *model.ReportRequest
	suite.Suite
}

var cfg *config.Config

// SetupTest method
func (s *IntegSuite) SetupTest() {
	// init config
	os.Setenv("Stage", "test")
	cfg = &config.Config{DefaultsFilePath: defaultsFP}
	err := cfg.Load()
	s.NoError(err)

	bankCardReportInput := &model.RequestInput{
		DateFrom:   periodDateFrom,
		DateTo:     periodDateTo,
		ReportType: bankCardReport,
	}
	s.bankCardReportReq, _ = validate.SetRequest(bankCardReportInput)

	employeeOSInput := &model.RequestInput{
		DateFrom:   periodDateFrom,
		DateTo:     periodDateTo,
		ReportType: employeeOSReport,
	}
	s.employeeOSReportReq, _ = validate.SetRequest(employeeOSInput)

	monthReportInput := &model.RequestInput{
		Date:       monthDate,
		ReportType: monthlyReport,
	}
	s.monthReportReq, _ = validate.SetRequest(monthReportInput)

	payPeriodReportInput := &model.RequestInput{
		DateFrom:   periodDateFrom,
		DateTo:     periodDateTo,
		ReportType: payPeriodReport,
	}
	s.payPeriodReportReq, _ = validate.SetRequest(payPeriodReportInput)

	productNumbersReportInput := &model.RequestInput{
		DateFrom:   periodDateFrom,
		DateTo:     periodDateTo,
		ReportType: productNumbersReport,
	}
	s.productNumbersReportReq, _ = validate.SetRequest(productNumbersReportInput)
}

// TestsetFileName method
func (s *IntegSuite) TestsetFileName() {
	var err error
	s.report, err = New(s.monthReportReq, cfg)
	s.NoError(err)
	s.report.setFileName()
	fileNm := s.report.getFileName()

	expectedFileNm := fmt.Sprintf("MonthlySalesReport_%s.xlsx", monthDate)
	s.Equal(fileNm, expectedFileNm)

	// now test report with date range
	s.report, err = New(s.payPeriodReportReq, cfg)
	s.NoError(err)
	s.report.setFileName()
	fileNm = s.report.getFileName()

	expectedFileNm = fmt.Sprintf("PayPeriodReport_%s_thru_%s.xlsx", periodDateFrom, periodDateTo)
	s.Equal(fileNm, expectedFileNm)
}

// TestGetRecords method
func (s *IntegSuite) TestGetRecords() {
	var err error

	s.report, err = New(s.bankCardReportReq, cfg)
	s.NoError(err)
	bc := &BankCard{
		dates: s.report.dates,
		db:    s.report.db,
	}
	bcRecs, err := bc.GetRecords()
	s.NoError(err)
	s.True(len(bcRecs) > 10)

	s.report, err = New(s.employeeOSReportReq, cfg)
	s.NoError(err)
	eo := &EmployeeOS{
		dates: s.report.dates,
		db:    s.report.db,
	}
	esRecs, err := eo.GetRecords()
	s.NoError(err)
	s.True(len(esRecs) > 10)

	s.report, err = New(s.monthReportReq, cfg)
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

	s.report, err = New(s.payPeriodReportReq, cfg)
	s.NoError(err)
	pp := &PayPeriod{
		dates: s.report.dates,
		db:    s.report.db,
	}
	ppRecs, err := pp.GetRecords()
	s.NoError(err)
	s.True(len(ppRecs) > 10)

	s.report, err = New(s.productNumbersReportReq, cfg)
	s.NoError(err)
	pn := &ProductNumbers{
		dates: s.report.dates,
		db:    s.report.db,
	}
	pnRecs, err := pn.GetRecords()
	s.NoError(err)
	s.True(len(pnRecs) > 10)
}

// Testcreate method
func (s *IntegSuite) Testcreate() {
	var err error
	s.report, err = New(s.bankCardReportReq, cfg)
	err = s.report.create()
	s.NoError(err)

	s.report, err = New(s.employeeOSReportReq, cfg)
	err = s.report.create()
	s.NoError(err)

	s.report, err = New(s.monthReportReq, cfg)
	err = s.report.create()
	s.NoError(err)

	s.report, err = New(s.payPeriodReportReq, cfg)
	err = s.report.create()
	s.NoError(err)

	s.report, err = New(s.productNumbersReportReq, cfg)
	err = s.report.create()
	s.NoError(err)
}

// TestSaveToDisk method
func (s *IntegSuite) TestSaveToDisk() {
	var err error
	/* s.report, err = New(s.bankCardReportReq, cfg)
	_, err = s.report.SaveToDisk("../tmp")
	s.NoError(err)

	s.report, err = New(s.employeeOSReportReq, cfg)
	s.NoError(err)
	_, err = s.report.SaveToDisk("../tmp")
	s.NoError(err)

	s.report, err = New(s.monthReportReq, cfg)
	_, err = s.report.SaveToDisk("../tmp")
	s.NoError(err) */

	s.report, err = New(s.payPeriodReportReq, cfg)
	_, err = s.report.SaveToDisk("../tmp")
	s.NoError(err)

	/* s.report, err = New(s.productNumbersReportReq, cfg)
	s.NoError(err)
	_, err = s.report.SaveToDisk("../tmp")
	s.NoError(err) */
}

// TestCreateSignedURL method
func (s *IntegSuite) TestCreateSignedURL() {
	var err error
	s.report, err = New(s.bankCardReportReq, cfg)
	url, err := s.report.CreateSignedURL()
	s.NoError(err)
	s.True(len(url) > 100)

	s.report, err = New(s.monthReportReq, cfg)
	url, err = s.report.CreateSignedURL()
	s.NoError(err)
	s.True(len(url) > 100)

	s.report, err = New(s.payPeriodReportReq, cfg)
	url, err = s.report.CreateSignedURL()
	s.NoError(err)
	s.True(len(url) > 100)
}

// TestIntegrationSuite function
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}
