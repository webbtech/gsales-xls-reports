package xlsx

import (
	"fmt"
	"time"

	"github.com/webbtech/gsales-xls-reports/model"
	"github.com/xuri/excelize/v2"
)

// BankCards method
func (x *XLSX) BankCards(records []*model.BankCardRecord) (err error) {

	f := x.file
	sheetNm := "Sheet1"

	index := f.NewSheet(sheetNm)
	f.SetActiveSheet(index)

	x.setHeader(sheetNm, BankCardJSON)
	x.setBankCardValues(sheetNm, records)
	x.setBankCardTotalsRow(sheetNm)
	f.SetSheetName(sheetNm, "Bank Cards")
	err = f.SetDocProps(&excelize.DocProperties{
		Title:   "Bank Cards Report",
		Created: time.Now().Format(time.RFC3339),
	})

	return err
}

func (x *XLSX) setBankCardValues(sheetNm string, records []*model.BankCardRecord) {

	col := 1
	row := 2

	firstRow = row
	lastRow = len(records) + 1

	for _, r := range records {
		x.displayCell(sheetNm, col, row, r.StationName)

		col++
		x.displayCell(sheetNm, col, row, r.RecordNumber)

		col++
		x.displayCell(sheetNm, col, row, r.BankAmex)

		col++
		x.displayCell(sheetNm, col, row, r.BankDiscover)

		col++
		x.displayCell(sheetNm, col, row, r.BankGales)

		col++
		x.displayCell(sheetNm, col, row, r.BankMC)

		col++
		x.displayCell(sheetNm, col, row, r.BankVisa)

		col++
		x.displayCell(sheetNm, col, row, r.CashDebit)

		col++
		x.displayCell(sheetNm, col, row, r.CashOther)

		col = 1
		row++
	}
}

func (x *XLSX) setBankCardTotalsRow(sheetNm string) {

	f := x.file
	totalsRow := lastRow + 1
	var cell, colNm, formula string
	const firstIteratorCol = 3
	const lastIteratorCol = 9
	var style int

	boldStyle, _ := f.NewStyle(`{"font":{"bold":true}}`)
	floatStyle, _ := f.NewStyle(`{"custom_number_format": "0.00; [red]0.00"}`)
	// numStyle, _ := f.NewStyle(`{}`)

	cell, _ = excelize.CoordinatesToCellName(1, totalsRow)
	f.SetCellValue(sheetNm, cell, "Totals")
	f.SetCellStyle(sheetNm, cell, cell, boldStyle)

	for c := firstIteratorCol; c <= lastIteratorCol; c++ {

		style = floatStyle

		colNm, _ = excelize.ColumnNumberToName(c)
		cell, _ = excelize.CoordinatesToCellName(c, totalsRow)
		formula = fmt.Sprintf("SUM(%s%d:%s%d)", colNm, firstRow, colNm, lastRow)
		f.SetCellFormula(sheetNm, cell, formula)
		f.SetCellStyle(sheetNm, cell, cell, style)
	}
}
