package xlsx

import (
	"fmt"
	"time"

	excelize "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/pulpfree/gsales-xls-reports/model"
)

// FIXME: these vars are global to all reports... whew need to fix that
var firstRow int
var lastRow int

// MonthlySales method
func (x *XLSX) MonthlySales(sales []*model.MonthlySaleRecord) (err error) {

	f := x.file
	sheetNm := "Sheet1"

	index := f.NewSheet(sheetNm)
	f.SetActiveSheet(index)

	x.setHeader(sheetNm, MonthlySalesJSON)
	x.setMonthlySalesValues(sheetNm, sales)
	x.setMonthlySalesTotalsRow(sheetNm)
	f.SetSheetName(sheetNm, "Monthly Sales")
	err = f.SetDocProps(&excelize.DocProperties{
		Title:   "Monthly Sales Report",
		Created: time.Now().Format(time.RFC3339),
	})

	return err
}

func (x *XLSX) setMonthlySalesValues(sheetNm string, sales []*model.MonthlySaleRecord) {

	col := 1
	row := 2

	firstRow = row
	lastRow = len(sales) + 1

	for _, s := range sales {

		x.displayCell(sheetNm, col, row, s.StationName)

		col++
		x.displayCell(sheetNm, col, row, s.RecordNumber)

		col++
		x.displayCell(sheetNm, col, row, s.Employee)

		col++
		x.displayCell(sheetNm, col, row, s.ShiftOvershort)

		col++
		x.displayCell(sheetNm, col, row, s.FuelSales)

		col++
		x.displayCell(sheetNm, col, row, s.FuelSalesHST)

		col++
		x.displayCell(sheetNm, col, row, s.FuelSalesTotal)

		col++
		x.displayCell(sheetNm, col, row, s.FuelAdjustments) // <-- this one!!

		col++
		x.displayCell(sheetNm, col, row, s.FuelSalesOther)

		col++
		x.displayCell(sheetNm, col, row, s.PropaneSales)

		col++
		x.displayCell(sheetNm, col, row, s.PropaneQty)

		col++
		x.displayCell(sheetNm, col, row, s.MiscNonFuelSales)

		col++
		x.displayCell(sheetNm, col, row, s.MiscNonFuelQty)

		col++
		x.displayCell(sheetNm, col, row, s.GiftCertificates)

		col++
		x.displayCell(sheetNm, col, row, s.BobsSales)

		col++
		x.displayCell(sheetNm, col, row, s.BobsGiftCertificates)

		col++
		x.displayCell(sheetNm, col, row, s.NonFuelTotal)

		col++
		x.displayCell(sheetNm, col, row, s.CashBills)

		col++
		x.displayCell(sheetNm, col, row, s.CashDebit)

		col++
		x.displayCell(sheetNm, col, row, s.CashDieselDiscount)

		col++
		x.displayCell(sheetNm, col, row, s.CashDriveOffNSF)

		col++
		x.displayCell(sheetNm, col, row, s.CashGalesLoyaltyRedeem)

		col++
		x.displayCell(sheetNm, col, row, s.CashGiftCertRedeem)

		col++
		x.displayCell(sheetNm, col, row, s.CashLotteryPayout)

		col++
		x.displayCell(sheetNm, col, row, s.CashOther)

		col++
		x.displayCell(sheetNm, col, row, s.CashOSAdjusted)

		col++
		x.displayCell(sheetNm, col, row, s.CashPayout)

		col++
		x.displayCell(sheetNm, col, row, s.CashWriteOff)

		col++
		x.displayCell(sheetNm, col, row, s.BankAmex)

		col++
		x.displayCell(sheetNm, col, row, s.BankDiscover)

		col++
		x.displayCell(sheetNm, col, row, s.BankGales)

		col++
		x.displayCell(sheetNm, col, row, s.BankMC)

		col++
		x.displayCell(sheetNm, col, row, s.BankVisa)

		col++
		x.displayCell(sheetNm, col, row, s.ProductCigarettesSales)

		col++
		x.displayCell(sheetNm, col, row, s.ProductCigarettesQty)

		col++
		x.displayCell(sheetNm, col, row, s.ProductOilSales)

		col++
		x.displayCell(sheetNm, col, row, s.ProductOilQty)

		col++
		x.displayCell(sheetNm, col, row, s.CarWash)

		col++
		x.displayCell(sheetNm, col, row, s.GalesLoyalty)

		col = 1
		row++
	}
}

func (x *XLSX) setMonthlySalesTotalsRow(sheetNm string) {

	f := x.file
	totalsRow := lastRow + 1
	numericCols := []int{11, 13, 35, 37, 38, 39}
	var cell, colNm, formula string
	const lastIteratorCol = 39
	var style int

	boldStyle, _ := f.NewStyle(`{"font":{"bold":true}}`)
	// floatStyle, _ := f.NewStyle(`{"custom_number_format": "0.00; [red]0.00"}`)
	floatStyle, _ := f.NewStyle(`{"number_format": 2}`)
	numStyle, _ := f.NewStyle(`{}`)

	cell, _ = excelize.CoordinatesToCellName(1, totalsRow)
	f.SetCellValue(sheetNm, cell, "Totals")
	f.SetCellStyle(sheetNm, cell, cell, boldStyle)

	for c := 4; c <= lastIteratorCol; c++ {
		if findNumber(numericCols, c) == true {
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
