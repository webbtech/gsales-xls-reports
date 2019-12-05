package report

import (
	"errors"
	"fmt"
	"path"

	"github.com/pulpfree/gsales-xls-reports/awsservices"
	"github.com/pulpfree/gsales-xls-reports/config"
	"github.com/pulpfree/gsales-xls-reports/model"
	"github.com/pulpfree/gsales-xls-reports/model/monthlysales"
	"github.com/pulpfree/gsales-xls-reports/model/payperiod"
	"github.com/pulpfree/gsales-xls-reports/xlsx"
)

// Report struct
type Report struct {
	allowedTypes []reportType
	cfg          *config.Config
	dates        *model.RequestDates
	file         *xlsx.XLSX
	filename     string
	reportType   reportType
}

// reportType string
type reportType string

// Constants
const (
	timeFormatShort               = "2006-01"
	timeFormatLong                = "2006-01-02"
	reportMonthlySales reportType = "monthlysales"
	reportPayPeriod    reportType = "payperiod"
)

var validType bool

// New function
func New(dates *model.RequestDates, cfg *config.Config, rType string) (r *Report, err error) {
	r = &Report{
		cfg:   cfg,
		dates: dates,
	}

	r.reportType, err = r.testReportType(rType)
	return r, err
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

	r.setFileName()

	switch r.reportType {
	case reportMonthlySales:
		return r.createMonthlySales()

	case reportPayPeriod:
		return r.createPayPeriod()
	}

	return err
}

func (r *Report) createMonthlySales() (err error) {

	sales, err := monthlysales.Init(r.dates, r.cfg)
	records, err := sales.GetRecords()
	defer sales.DB.Close()

	r.file, err = xlsx.NewFile()
	err = r.file.MonthlySales(records)

	return err
}

func (r *Report) createPayPeriod() (err error) {

	pp, err := payperiod.Init(r.dates, r.cfg)
	records, err := pp.GetRecords()
	defer pp.DB.Close()

	r.file, err = xlsx.NewFile()
	err = r.file.PayPeriod(records)

	return err
}

// ===================== Helper Methods ======================================================== //

func (r *Report) setFileName() {
	switch r.reportType {
	case reportMonthlySales:
		r.filename = fmt.Sprintf("MonthlySalesReport_%s.xlsx", r.dates.DateFrom.Format(timeFormatShort))
	case reportPayPeriod:
		r.filename = fmt.Sprintf("PayPeriodReport_%s_thru_%s.xlsx", r.dates.DateFrom.Format(timeFormatLong), r.dates.DateTo.Format(timeFormatLong))
	}
}

func (r *Report) getFileName() string {
	return r.filename
}

func (r *Report) setAllowedTypes() {
	r.allowedTypes = make([]reportType, 2)
	r.allowedTypes[0] = reportMonthlySales
	r.allowedTypes[1] = reportPayPeriod
}

func (r *Report) testReportType(rType string) (rt reportType, err error) {

	if len(r.allowedTypes) == 0 {
		r.setAllowedTypes()
	}
	for _, at := range r.allowedTypes {
		strType := string(at)
		if rType == strType {
			rt = at
			break
		}
	}

	if rt == "" {
		errStr := fmt.Sprintf("Invalid report type: %s", rType)
		err = errors.New(errStr)
	}

	return rt, err
}
