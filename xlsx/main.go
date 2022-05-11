package xlsx

import (
	"bytes"
	"encoding/json"
	"math"

	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

// Defaults
const (
	abc             = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	floatFmt        = "#,#0"
	timeShortForm   = "20060102"
	timeMonthForm   = "200601"
	dateDayFormat   = "Jan _2"
	dateMonthFormat = "January 2006"
)

// XLSX struct
type XLSX struct {
	file *excelize.File
}

// HeaderVal struct
type HeaderVal struct {
	Column string  `json:"column"`
	Label  string  `json:"label"`
	Width  float64 `json:"width"`
}

// getHeaderJSON type
type getHeaderJSON = func() string

// NewFile function
func NewFile() (x *XLSX, err error) {

	x = new(XLSX)
	x.file = excelize.NewFile()
	if err != nil {
		log.Errorf("xlsx err %s: ", err)
	}
	return x, err
}

// OutputFile method
func (x *XLSX) OutputFile() (buf bytes.Buffer, err error) {
	err = x.file.Write(&buf)
	if err != nil {
		log.Errorf("xlsx err: %s", err)
	}
	return buf, err
}

// OutputToDisk method
func (x *XLSX) OutputToDisk(path string) (fp string, err error) {
	err = x.file.SaveAs(path)
	return path, err
}

func (x *XLSX) displayCell(sheetNm string, col int, row int, value interface{}) {

	f := x.file
	// Apparently, setting a style on a float value produces an error when opening in excel... whew!!
	// floatStyle, _ := f.NewStyle(`{"custom_number_format": "0.00; [Red]0.00"}`)
	// floatStyle, _ := f.NewStyle(`{"number_format": 3}`)
	defStyle, _ := f.NewStyle(`{}`)

	cell, _ := excelize.CoordinatesToCellName(col, row)

	switch value.(type) {
	case float64:
		f.SetCellFloat(sheetNm, cell, value.(float64), 2, 64)

	default:
		f.SetCellValue(sheetNm, cell, value)
		f.SetCellStyle(sheetNm, cell, cell, defStyle)
	}
}

// setHeader method
func (x *XLSX) setHeader(sheetNm string, jsonFunc getHeaderJSON) {

	var cell string
	var hdrs []HeaderVal

	f := x.file
	rowNo := 1
	style, _ := f.NewStyle(`{"font":{"bold":true}}`)

	json.Unmarshal([]byte(jsonFunc()), &hdrs)

	firstCell, _ := excelize.CoordinatesToCellName(1, rowNo)
	lastCell, _ := excelize.CoordinatesToCellName(len(hdrs), rowNo)

	for i, h := range hdrs {
		col := i + 1
		colNm, _ := excelize.ColumnNumberToName(col)
		cell, _ = excelize.CoordinatesToCellName(col, rowNo)
		f.SetCellValue(sheetNm, cell, h.Label)
		f.SetColWidth(sheetNm, colNm, colNm, h.Width)
	}

	f.SetCellStyle(sheetNm, firstCell, lastCell, style)
}

// ======================== Helper Methods ================================= //

// see: https://stackoverflow.com/questions/36803999/golang-alphabetic-representation-of-a-number
// for a way to map int to letters
func toChar(i int) string {
	return abc[i-1 : i]
}

// Found these function at: https://stackoverflow.com/questions/18390266/how-can-we-truncate-float64-type-to-a-particular-precision-in-golang
// Looks like a good way to deal with precision
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

// findNumber function
// takes a slice and looks for an element in it
// return bool
func findNumber(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

/* func displayFloat(num interface{}) float64 {
	var ret float64
	switch v := num.(type) {
	case *float64:
		// need to check for nil here to deal with null db values
		if v == nil {
			ret = 0.00
		} else {
			ret = *v
		}
	case float64:
		ret = v
	default:
		ret = 0.00
	}

	return ret
} */
