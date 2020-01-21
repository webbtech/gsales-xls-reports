package model

import (
	"errors"
)

// ReportType int
type ReportType int

// Constants
const (
	BankCardsReport ReportType = iota + 1
	EmployeeOSReport
	MonthlySalesReport
	PayPeriodReport
	ProductNumbersReport
)

// ReportStringToType function
func ReportStringToType(rType string) (ReportType, error) {
	var rt ReportType

	switch rType {
	case "bankcards":
		rt = BankCardsReport
	case "employeeos":
		rt = EmployeeOSReport
	case "monthlysales":
		rt = MonthlySalesReport
	case "payperiod":
		rt = PayPeriodReport
	case "productnumbers":
		rt = ProductNumbersReport
	default:
		rt = 0
	}
	if rt == 0 {
		return rt, errors.New("Invalid report type request")
	}
	return rt, nil
}
