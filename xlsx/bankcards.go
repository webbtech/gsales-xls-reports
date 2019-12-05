package xlsx

import "github.com/pulpfree/gsales-xls-reports/model"

// BankCards method
func (x *XLSX) BankCards(records []*model.BankCardRecord) (err error) {

	f := x.file
	sheetNm := "Sheet1"

	index := f.NewSheet(sheetNm)
	f.SetActiveSheet(index)

	x.setHeader(sheetNm, BankCardJSON)
	// x.setPayPeriodValues(sheetNm, records)
	// x.setPayPeriodTotalsRow(sheetNm)
	f.SetSheetName(sheetNm, "Bank Cards")

	return err
}
