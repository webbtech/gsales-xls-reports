package xlsx

import (
	"fmt"

	"github.com/webbtech/gsales-xls-reports/model"
	"github.com/xuri/excelize/v2"
)

// PayPeriod method
func (x *XLSX) PayPeriod(records []*model.PayPeriodRecord) (err error) {

	f := x.file
	sheetNm := "Sheet1"

	index := f.NewSheet(sheetNm)
	f.SetActiveSheet(index)

	x.setHeader(sheetNm, PayPeriodJSON)
	x.setPayPeriodValues(sheetNm, records)
	x.setPayPeriodTotalsRow(sheetNm)
	f.SetSheetName(sheetNm, "Pay Period")

	return err
}

func (x *XLSX) setPayPeriodValues(sheetNm string, records []*model.PayPeriodRecord) {

	col := 1
	row := 2

	firstRow = row
	lastRow = len(records) + 1

	for _, r := range records {

		x.displayCell(sheetNm, col, row, r.Employee)

		col++
		x.displayCell(sheetNm, col, row, r.RecordNumber)

		col++
		x.displayCell(sheetNm, col, row, r.StationName)

		col++
		x.displayCell(sheetNm, col, row, r.ShiftOvershort)

		col++
		x.displayCell(sheetNm, col, row, r.NonFuelSales)

		col++
		x.displayCell(sheetNm, col, row, r.ProductSales)

		col++
		x.displayCell(sheetNm, col, row, r.Commission.Qty)

		col++
		x.displayCell(sheetNm, col, row, r.Commission.Sales)

		col++
		x.displayCell(sheetNm, col, row, r.Commission.Commission)

		col++
		x.displayCell(sheetNm, col, row, r.CarwashNumber)

		col++
		x.displayCell(sheetNm, col, row, r.GalesLoyaltyQty)

		col++
		x.displayCell(sheetNm, col, row, r.AttendantAdjustment)

		col = 1
		row++
	}
}

func (x *XLSX) setPayPeriodTotalsRow(sheetNm string) {

	f := x.file
	totalsRow := lastRow + 1
	var cell, colNm, formula string
	const firstIteratorCol = 4
	const lastIteratorCol = 11
	var style int

	boldStyle, _ := f.NewStyle(`{"font":{"bold":true}}`)
	floatStyle, _ := f.NewStyle(`{"custom_number_format": "0.00; [red]0.00"}`)
	numStyle, _ := f.NewStyle(`{}`)

	cell, _ = excelize.CoordinatesToCellName(1, totalsRow)
	f.SetCellValue(sheetNm, cell, "Totals")
	f.SetCellStyle(sheetNm, cell, cell, boldStyle)

	for c := firstIteratorCol; c <= lastIteratorCol; c++ {
		if c == 7 || c == 10 || c == 11 {
			style = numStyle
		} else {
			style = floatStyle
		}

		colNm, _ = excelize.ColumnNumberToName(c)
		cell, _ = excelize.CoordinatesToCellName(c, totalsRow)
		formula = fmt.Sprintf("SUM(%s%d:%s%d)", colNm, firstRow, colNm, lastRow)
		f.SetCellFormula(sheetNm, cell, formula)
		f.SetCellStyle(sheetNm, cell, cell, style)
	}
}
