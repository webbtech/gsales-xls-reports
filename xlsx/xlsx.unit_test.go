package xlsx

import (
	"testing"
)

// TestcreateFormula_Success
func TestCreateCellsFormula_Success(t *testing.T) {

	expectedStr := "B2+D2+F2+H2"
	startCol := 2
	numCols := 4
	row := 2
	frm := createCellsFormula(startCol, numCols, row)

	if frm != expectedStr {
		t.Errorf("Expect formula to be: %s got %s", expectedStr, frm)
	}
}
