package xlsx

import (
	"fmt"
	"time"

	"github.com/webbtech/gsales-xls-reports/model"
	"github.com/xuri/excelize/v2"
)

// ProductNumbers method
func (x *XLSX) ProductNumbers(records []*model.ProductNumberRecord) (err error) {

	f := x.file
	sheetNm := "Sheet1"

	index := f.NewSheet(sheetNm)
	f.SetActiveSheet(index)

	x.setHeader(sheetNm, ProductNumbersJSON)
	x.setProductNumbersValues(sheetNm, records)
	x.setProductNumbersTotalsRow(sheetNm)
	f.SetSheetName(sheetNm, "Product Numbers")
	err = f.SetDocProps(&excelize.DocProperties{
		Title:   "Product Numbers Report",
		Created: time.Now().Format(time.RFC3339),
	})

	return err
}

func (x *XLSX) setProductNumbersValues(sheetNm string, records []*model.ProductNumberRecord) {

	col := 1
	row := 2

	firstRow = row
	lastRow = len(records) + 1

	for _, r := range records {

		x.displayCell(sheetNm, col, row, r.Product)

		col++
		x.displayCell(sheetNm, col, row, r.Qty)

		col = 1
		row++
	}
}

func (x *XLSX) setProductNumbersTotalsRow(sheetNm string) {

	f := x.file
	totalsRow := lastRow + 1
	var cell, colNm, formula string
	const firstIteratorCol = 2
	const lastIteratorCol = 2

	boldStyle, _ := f.NewStyle(`{"font":{"bold":true}}`)

	cell, _ = excelize.CoordinatesToCellName(1, totalsRow)
	f.SetCellValue(sheetNm, cell, "Total")
	f.SetCellStyle(sheetNm, cell, cell, boldStyle)

	c := 2
	colNm, _ = excelize.ColumnNumberToName(c)
	cell, _ = excelize.CoordinatesToCellName(c, totalsRow)
	formula = fmt.Sprintf("SUM(%s%d:%s%d)", colNm, firstRow, colNm, lastRow)
	f.SetCellFormula(sheetNm, cell, formula)
}
