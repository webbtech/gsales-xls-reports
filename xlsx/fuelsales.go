package xlsx

import (
	"fmt"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/pulpfree/gsales-xls-reports/model"
)

// stationTotal struct
type stationTotal struct {
	stationNm string
	model.FuelType
}

// FuelSales method
func (x *XLSX) FuelSales(records []model.FuelSalesRecord) (err error) {

	f := x.file
	sheetNm := "Sheet1"

	index := f.NewSheet(sheetNm)
	f.SetActiveSheet(index)

	x.setFSHeader(sheetNm)

	x.setFuelSalesValues(sheetNm, records)
	x.setFuelSalesTotals(sheetNm)
	x.setFuelSalesStationTotals(sheetNm, records)
	f.SetSheetName(sheetNm, "Fuel Sales")
	err = f.SetDocProps(&excelize.DocProperties{
		Title:   "Fuel Sales Report",
		Created: time.Now().Format(time.RFC3339),
	})

	return err
}

// setFSHeader method
// This is a custom header with 2 rows and cell merging, hence a new method
func (x *XLSX) setFSHeader(sheetNm string) {

	var cell string

	colWidth := 10.00
	f := x.file
	rowNo := 1
	style, _ := f.NewStyle(`{"font":{"bold":true}}`)

	firstCell, _ := excelize.CoordinatesToCellName(1, rowNo)
	lastCell, _ := excelize.CoordinatesToCellName(13, rowNo)
	f.SetCellValue(sheetNm, firstCell, "Station")
	colNm, _ := excelize.ColumnNumberToName(1)
	f.SetColWidth(sheetNm, colNm, colNm, colWidth)
	f.SetCellStyle(sheetNm, firstCell, lastCell, style)

	_ = f.MergeCell(sheetNm, "B1", "C1")
	f.SetCellValue(sheetNm, "B1", "NL")
	_ = f.MergeCell(sheetNm, "D1", "E1")
	f.SetCellValue(sheetNm, "D1", "SNL")
	_ = f.MergeCell(sheetNm, "F1", "G1")
	f.SetCellValue(sheetNm, "F1", "DSL")
	_ = f.MergeCell(sheetNm, "H1", "I1")
	f.SetCellValue(sheetNm, "H1", "CDSL")
	_ = f.MergeCell(sheetNm, "J1", "K1")
	f.SetCellValue(sheetNm, "J1", "Propane")
	_ = f.MergeCell(sheetNm, "L1", "M1")
	f.SetCellValue(sheetNm, "L1", "Station Totals")
	f.SetColWidth(sheetNm, "B", "M", 12)

	rowNo++
	col := 2
	for i := 1; i <= 6; i++ {
		cell, _ = excelize.CoordinatesToCellName(col, rowNo)
		f.SetCellValue(sheetNm, cell, "Litres")
		col++

		cell, _ = excelize.CoordinatesToCellName(col, rowNo)
		f.SetCellValue(sheetNm, cell, "Dollars")
		col++
	}
}

func (x *XLSX) setFuelSalesValues(sheetNm string, records []model.FuelSalesRecord) {

	col := 1
	row := 3

	firstRow = row
	lastRow = len(records) + 1

	for _, r := range records {
		fmt.Printf("r.StationName: %+v\n", r.StationName)
		fmt.Printf("row: %+v\n", row)
		x.displayCell(sheetNm, col, row, r.StationName)

		col++
		x.displayCell(sheetNm, col, row, r.Fuel1.Litre)

		col++
		x.displayCell(sheetNm, col, row, r.Fuel1.Dollar)

		col++
		x.displayCell(sheetNm, col, row, r.Fuel3.Litre)

		col++
		x.displayCell(sheetNm, col, row, r.Fuel3.Dollar)

		col++
		x.displayCell(sheetNm, col, row, r.Fuel4.Litre)

		col++
		x.displayCell(sheetNm, col, row, r.Fuel4.Dollar)

		col++
		x.displayCell(sheetNm, col, row, r.Fuel5.Litre)

		col++
		x.displayCell(sheetNm, col, row, r.Fuel5.Dollar)

		col++
		x.displayCell(sheetNm, col, row, r.Fuel6.Litre)

		col++
		x.displayCell(sheetNm, col, row, r.Fuel6.Dollar)

		col = 1
		row++
	}
}

func (x *XLSX) setFuelSalesTotals(sheetNm string) {

	f := x.file
	totalsRow := lastRow + 2
	var cell, colNm, formula string
	const firstIteratorCol = 2
	const lastIteratorCol = 13
	var style int

	boldStyle, _ := f.NewStyle(`{"font":{"bold":true}}`)
	floatStyle, _ := f.NewStyle(`{"custom_number_format": "0.00; [red]0.00"}`)

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

func (x *XLSX) setFuelSalesStationTotals(sheetNm string, records []model.FuelSalesRecord) {

	var colIter int
	var formula string
	var style int

	f := x.file
	startCol := 2
	startRow := 3
	curRow := startRow
	numRows := len(records) - 1 // Not clear why this is necessary, wee have 14 stations but len(records) shows 15???
	numFuels := 5
	totalStartCol := 12

	floatStyle, _ := f.NewStyle(`{"custom_number_format": "0.00; [red]0.00"}`)

	for c := 0; c < numRows; c++ {

		style = floatStyle
		colIter = totalStartCol
		rowIter := c + startRow

		// set Litres total formulas
		ltrFormula := createCellsFormula(startCol, numFuels, rowIter)
		ltrCell, _ := excelize.CoordinatesToCellName(colIter, curRow)
		formula = fmt.Sprintf("SUM(%s)", ltrFormula)
		f.SetCellFormula(sheetNm, ltrCell, formula)
		f.SetCellStyle(sheetNm, ltrCell, ltrCell, style)

		// set Dollars total formulas
		dlrFormula := createCellsFormula(startCol+1, numFuels, rowIter)
		dlrCell, _ := excelize.CoordinatesToCellName(colIter+1, curRow)
		formula = fmt.Sprintf("SUM(%s)", dlrFormula)
		f.SetCellFormula(sheetNm, dlrCell, formula)
		f.SetCellStyle(sheetNm, dlrCell, dlrCell, style)
		curRow++
	}
}

func createCellsFormula(startCol, numCols, row int) string {

	cells := []string{}
	curCol := startCol
	for i := 0; i < numCols; i++ {
		cell, _ := excelize.CoordinatesToCellName(curCol, row)
		cells = append(cells, cell)
		curCol += 2
	}
	return strings.Join(cells[:], "+")
}
