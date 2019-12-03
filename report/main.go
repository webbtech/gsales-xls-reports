package report

import (
	"errors"
	"fmt"
	"path"

	"github.com/pulpfree/gsales-xls-reports/awsservices"
	"github.com/pulpfree/gsales-xls-reports/config"
	"github.com/pulpfree/gsales-xls-reports/model"
	"github.com/pulpfree/gsales-xls-reports/model/monthlysales"
	"github.com/pulpfree/gsales-xls-reports/xlsx"
)

// Report struct
type Report struct {
	allowedTypes map[string]bool
	cfg          *config.Config
	dates        *model.RequestDates
	file         *xlsx.XLSX
	filename     string
	haveFile     bool
	reportType   string
}

// Constants
const (
	reportFileName = "MonthlySalesReport"
	timeFormat     = "2006-01"
)

// New function
func New(dates *model.RequestDates, cfg *config.Config, reportType string) (r *Report, err error) {
	r = &Report{
		cfg:        cfg,
		dates:      dates,
		haveFile:   false,
		reportType: reportType,
	}

	// test report type
	r.setAllowedTypes()
	_, ok := r.allowedTypes[reportType]
	if ok != true {
		err = errors.New("Invalid report type")
	}

	return r, err
}

// Create method
func (r *Report) Create() (err error) {

	r.setFileName()

	switch r.reportType {
	case "monthlysales":
		return r.createMonthlySales()
	}

	return err
}

// SaveToDisk method
func (r *Report) SaveToDisk(dir string) (fp string, err error) {

	if r.haveFile == false {
		err = r.Create()
		if err != nil {
			return fp, err
		}
	}

	filePath := path.Join(dir, r.getFileName())
	fp, err = r.file.OutputToDisk(filePath)

	return fp, err
}

// CreateSignedURL method
func (r *Report) CreateSignedURL() (url string, err error) {

	if r.haveFile == false {
		err = r.Create()
		if err != nil {
			return url, err
		}
	}

	fileOutput, err := r.file.OutputFile()
	if err != nil {
		return url, err
	}

	s3Service, err := awsservices.NewS3(r.cfg)
	filePrefix := r.getFileName()

	return s3Service.GetSignedURL(filePrefix, &fileOutput)
}

func (r *Report) createMonthlySales() (err error) {

	sales, err := monthlysales.Init(r.dates, r.cfg)
	records, err := sales.GetRecords()
	defer sales.DB.Close()

	r.file, err = xlsx.NewFile()
	err = r.file.MonthlySales(records)
	r.haveFile = true

	return err
}

// ======================== Helper Methods =============================== //

func (r *Report) setFileName() {
	r.filename = fmt.Sprintf("%s_%s.xlsx", reportFileName, r.dates.DateFrom.Format(timeFormat))
}

func (r *Report) getFileName() string {
	return r.filename
}

func (r *Report) setAllowedTypes() {
	r.allowedTypes = make(map[string]bool)
	r.allowedTypes["monthlysales"] = true
	r.allowedTypes["payperiod"] = true
}
