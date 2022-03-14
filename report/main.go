package report

import (
	"errors"
	"fmt"
	"path"

	"github.com/pulpfree/gsales-xls-reports/awsservices"
	"github.com/pulpfree/gsales-xls-reports/config"
	"github.com/pulpfree/gsales-xls-reports/model"
	"github.com/pulpfree/gsales-xls-reports/model/db"
	"github.com/pulpfree/gsales-xls-reports/xlsx"
)

// Report struct
type Report struct {
	cfg        *config.Config
	dates      *model.RequestDates
	db         model.DBHandler
	file       *xlsx.XLSX
	filename   string
	reportType *model.ReportType
}

// Constants
const (
	timeFormatShort = "2006-01"
	timeFormatLong  = "2006-01-02"
)

// New function
func New(req *model.ReportRequest, cfg *config.Config) (report *Report, err error) {
	db, err := db.NewDB(cfg.GetMongoConnectURL(), cfg.DBName)
	if err != nil {
		return nil, err
	}

	return &Report{
		cfg:        cfg,
		dates:      req.Dates,
		db:         db,
		reportType: req.ReportType,
	}, err
}

// ===================== Exported Methods ====================================================== //

// CreateSignedURL method
func (r *Report) CreateSignedURL() (url string, err error) {

	err = r.create()
	if err != nil {
		return url, err
	}

	fileOutput, err := r.file.OutputFile()
	if err != nil {
		return url, err
	}

	s3Service, err := awsservices.NewS3(r.cfg)
	filePrefix := r.getFileName()

	return s3Service.GetSignedURL(filePrefix, &fileOutput)
}

// SaveToDisk method
func (r *Report) SaveToDisk(dir string) (fp string, err error) {

	err = r.create()
	if err != nil {
		return fp, err
	}

	filePath := path.Join(dir, r.getFileName())
	fp, err = r.file.OutputToDisk(filePath)

	return fp, err
}

// ===================== Un-exported Methods =================================================== //

// create method
func (r *Report) create() (err error) {

	err = r.setFileName()
	if err != nil {
		return err
	}

	rt := *r.reportType
	switch rt {
	case model.BankCardsReport:
		return r.createBankCardsReport()
	case model.EmployeeOSReport:
		return r.createEmployeeOSReport()
	case model.FuelSalesReport:
		return r.createFuelSalesReport()
	case model.MonthlySalesReport:
		return r.createMonthlySales()
	case model.PayPeriodReport:
		return r.createPayPeriod()
	case model.ProductNumbersReport:
		return r.createProductNumbers()
	default:
		errStr := fmt.Sprintf("Invalid report type: %v", rt)
		err = errors.New(errStr)
	}

	return err
}

func (r *Report) createBankCardsReport() (err error) {

	bc := &BankCard{
		dates: r.dates,
		db:    r.db,
	}

	records, err := bc.GetRecords()
	defer r.db.Close()

	r.file, err = xlsx.NewFile()
	err = r.file.BankCards(records)
	return err
}

func (r *Report) createEmployeeOSReport() (err error) {

	eo := &EmployeeOS{
		dates: r.dates,
		db:    r.db,
	}

	records, err := eo.GetRecords()
	defer r.db.Close()

	r.file, err = xlsx.NewFile()
	err = r.file.EmployeeOS(records)

	return err
}

func (r *Report) createFuelSalesReport() (err error) {

	rep := &FuelSales{
		dates: r.dates,
		db:    r.db,
	}

	records, err := rep.GetRecords()
	defer r.db.Close()

	r.file, err = xlsx.NewFile()
	err = r.file.FuelSales(records)

	return err
}

func (r *Report) createMonthlySales() (err error) {

	cfg, err := r.db.GetConfig()
	if err != nil {
		return err
	}

	ms := &MonthlySales{
		cfg:   cfg,
		dates: r.dates,
		db:    r.db,
	}

	records, err := ms.GetRecords()
	defer r.db.Close()

	r.file, err = xlsx.NewFile()
	err = r.file.MonthlySales(records)

	return err
}

func (r *Report) createPayPeriod() (err error) {

	pp := &PayPeriod{
		dates: r.dates,
		db:    r.db,
	}
	records, err := pp.GetRecords()
	defer pp.db.Close()
	if err != nil {
		return err
	}

	r.file, err = xlsx.NewFile()
	err = r.file.PayPeriod(records)

	return err
}

func (r *Report) createProductNumbers() (err error) {

	pn := &ProductNumbers{
		dates: r.dates,
		db:    r.db,
	}
	records, err := pn.GetRecords()
	defer pn.db.Close()

	r.file, err = xlsx.NewFile()
	err = r.file.ProductNumbers(records)

	return err
}

// ===================== Helper Methods ======================================================== //

func (r *Report) setFileName() (err error) {
	rt := *r.reportType
	switch rt {
	case model.BankCardsReport:
		r.filename = fmt.Sprintf("BankCardsReport_%s_thru_%s.xlsx", r.dates.DateFrom.Format(timeFormatLong), r.dates.DateTo.Format(timeFormatLong))
	case model.EmployeeOSReport:
		r.filename = fmt.Sprintf("EmployeeOSReport_%s_thru_%s.xlsx", r.dates.DateFrom.Format(timeFormatLong), r.dates.DateTo.Format(timeFormatLong))
	case model.FuelSalesReport:
		r.filename = fmt.Sprintf("FuelSalesReport_%s.xlsx", r.dates.DateFrom.Format(timeFormatShort))
	case model.MonthlySalesReport:
		r.filename = fmt.Sprintf("MonthlySalesReport_%s.xlsx", r.dates.DateFrom.Format(timeFormatShort))
	case model.PayPeriodReport:
		r.filename = fmt.Sprintf("PayPeriodReport_%s_thru_%s.xlsx", r.dates.DateFrom.Format(timeFormatLong), r.dates.DateTo.Format(timeFormatLong))
	case model.ProductNumbersReport:
		r.filename = fmt.Sprintf("ProductNumbersReport_%s_thru_%s.xlsx", r.dates.DateFrom.Format(timeFormatLong), r.dates.DateTo.Format(timeFormatLong))
	default:
		errStr := fmt.Sprintf("Missing or invalid report type %v in setFileName method", rt)
		err = errors.New(errStr)
	}
	return err
}

func (r *Report) getFileName() string {
	return r.filename
}
