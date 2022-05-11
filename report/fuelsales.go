package report

import (
	"github.com/webbtech/gsales-xls-reports/model"
)

// FuelSales struct
type FuelSales struct {
	db      model.DBHandler
	dates   *model.RequestDates
	records []model.FuelSalesRecord
}

// ======================== Exported Methods =================================================== //

// GetRecords method
func (rep *FuelSales) GetRecords() ([]model.FuelSalesRecord, error) {

	var err error
	err = rep.setRecords()
	if err != nil {
		return nil, err
	}

	return rep.records, err
}

// ======================== Un-exported Methods ================================================ //

func (rep *FuelSales) setRecords() (err error) {

	sales, err := rep.db.GetFuelSales(rep.dates)
	if err != nil {
		return err
	}

	for _, s := range sales {
		fs, err := setFuels(s.Fuels)
		if err != nil {
			return err
		}
		record := model.FuelSalesRecord{
			StationName: s.StationName,
			Fuel1:       fs[1],
			Fuel2:       fs[2],
			Fuel3:       fs[3],
			Fuel4:       fs[4],
			Fuel5:       fs[5],
		}
		rep.records = append(rep.records, record)
	}

	return err
}

func setFuels(fuelsInput []model.Fuel) (fuels map[int]model.FuelType, err error) {

	fuels = make(map[int]model.FuelType)

	// set some default values
	fuels[1] = model.FuelType{}
	fuels[2] = model.FuelType{}
	fuels[3] = model.FuelType{}
	fuels[4] = model.FuelType{}
	fuels[5] = model.FuelType{}

	for _, f := range fuelsInput {
		fuels[f.Grade] = model.FuelType{
			Dollar: f.Dollar,
			Litre:  f.Litre,
		}
	}

	haveFuel := haveGrade2(fuels)
	if haveFuel {
		f2LtSplit := fuels[2].Litre / 2
		f1LtCost := fuels[1].Dollar / fuels[1].Litre
		f3LtCost := fuels[3].Dollar / fuels[3].Litre
		// here, for whatever reason (thoroughness perhaps?), we're calculating the % increase
		// in cost for the hi-grade (grade 3) fuel
		// fCostDiff := f3LtCost - f1LtCost
		// fCostPrct := fCostDiff / f1LtCost
		// f3DlrAdd := f2LtSplit * (1 + fCostPrct)

		// I've discovered that just using the calculated cost from f3LtCost was more accurate
		// helluva lot simpler too!
		f3DlrAdd := f2LtSplit * f3LtCost
		f1DlrAdd := f2LtSplit * f1LtCost

		// to test the variance uncomment below
		/* f2DlrChck := f1DlrAdd + f3DlrAdd
		DlrChk := f2DlrChck - fuels[2].Dollar
		fmt.Printf("DlrChk: %+v\n", DlrChk) */

		// zero out fuel2
		fuels[2] = model.FuelType{}
		fuel1Tmp := fuels[1]
		fuel3Tmp := fuels[3]
		fuels[1] = model.FuelType{
			Dollar: fuel1Tmp.Dollar + f1DlrAdd,
			Litre:  fuel1Tmp.Litre + f2LtSplit,
		}
		fuels[3] = model.FuelType{
			Dollar: fuel3Tmp.Dollar + f3DlrAdd,
			Litre:  fuel3Tmp.Litre + f2LtSplit,
		}
	}

	return fuels, err
}

func haveGrade2(fuels map[int]model.FuelType) bool {
	if fuels[2].Litre > 0 {
		return true
	}
	return false
}
