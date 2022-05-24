package report

import "github.com/webbtech/gsales-xls-reports/model"

// ProductNumbers struct
type ProductNumbers struct {
	db      model.DbHandler
	dates   *model.RequestDates
	records []*model.ProductNumberRecord
}

// ======================== Exported Methods =================================================== //

// GetRecords method
func (pn *ProductNumbers) GetRecords() ([]*model.ProductNumberRecord, error) {

	var err error
	records, err := pn.db.GetProductNumbers(pn.dates)
	if err != nil {
		return nil, err
	}

	return records, err
}
