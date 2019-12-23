package xlsx

import (
	"fmt"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/pulpfree/gsales-xls-reports/model"
)

// EmployeeOS method
func (x *XLSX) EmployeeOS(records []*model.EmployeeOSRecord) (err error) {

	f := x.file
	sheetNm := "Sheet1"

	index := f.NewSheet(sheetNm)
	f.SetActiveSheet(index)

	x.setHeader(sheetNm, EmployeeOSJSON)
	x.setEmployeeOSValues(sheetNm, records)
	x.setEmployeeOSTotalsRow(sheetNm)
	f.SetSheetName(sheetNm, "Employee Overshort")
	err = f.SetDocProps(&excelize.DocProperties{
		Title:   "Bank Cards Report",
		Created: time.Now().Format(time.RFC3339),
	})

	return err
}

func (x *XLSX) setEmployeeOSValues(sheetNm string, records []*model.EmployeeOSRecord) {

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
		x.displayCell(sheetNm, col, row, r.OvershortShift)

		col++
		x.displayCell(sheetNm, col, row, r.OvershortAttendant)

		col++
		x.displayCell(sheetNm, col, row, r.OvershortDiff)

		col++
		x.displayCell(sheetNm, col, row, r.DiscrepancyDescription)

		col = 1
		row++
	}
}

func (x *XLSX) setEmployeeOSTotalsRow(sheetNm string) {

	f := x.file
	totalsRow := lastRow + 1
	var cell, colNm, formula string
	const firstIteratorCol = 4
	const lastIteratorCol = 6
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
