package report

import (
	"github.com/pulpfree/gsales-xls-reports/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PayPeriod struct
type PayPeriod struct {
	db      model.DBHandler
	dates   *model.RequestDates
	records []*model.PayPeriodRecord
}

// ======================== Exported Methods =================================================== //

// GetRecords method
func (pp *PayPeriod) GetRecords() ([]*model.PayPeriodRecord, error) {

	var err error
	err = pp.setRecords()
	if err != nil {
		return nil, err
	}

	err = pp.setNonFuelCommission()
	if err != nil {
		return nil, err
	}

	err = pp.setCarWashes()
	if err != nil {
		return nil, err
	}

	err = pp.setGalesLoyalty()
	if err != nil {
		return nil, err
	}

	return pp.records, err
}

// ======================== Un-exported Methods ================================================ //

func (pp *PayPeriod) setRecords() (err error) {

	sales, err := pp.db.GetPayPeriodSales(pp.dates)
	if err != nil {
		return err
	}

	stationMap, err := pp.db.GetStationMap()
	if err != nil {
		return err
	}

	for _, s := range sales {

		employee, err := pp.db.GetEmployee(s.Attendant.ID)
		if err != nil {
			return err
		}

		record := &model.PayPeriodRecord{
			AttendantAdjustment: model.SetString(s.Attendant.Adjustment),
			Employee:            employee,
			NonFuelSales:        model.SetFloat(s.Summary.TotalNonFuel),
			ProductSales:        model.SetFloat(s.Summary.Product),
			RecordNumber:        s.RecordNum,
			StationID:           s.StationID,
			StationName:         stationMap[s.StationID].Name,
			ShiftOvershort:      s.Overshort.Amount,
		}
		pp.records = append(pp.records, record)
	}

	return err
}

// setNonFuelCommission method
func (pp *PayPeriod) setNonFuelCommission() (err error) {
	for _, s := range pp.records {
		commission, err := pp.db.GetNonFuelCommission(s.RecordNumber, s.StationID)
		if err != nil {
			return err
		}
		s.Commission = commission
	}
	return err
}

// setCarWashes method
func (pp *PayPeriod) setCarWashes() (err error) {
	washes, err := pp.db.GetCarWash(pp.dates)
	if err != nil {
		return err
	}
	for _, s := range pp.records {
		if pp.testStation(s.StationID) {
			s.CarwashNumber = pp.searchCarWashSale(washes, s)
		}
	}
	return err
}

func (pp *PayPeriod) searchCarWashSale(carWashes []*model.NonFuelSale, rec *model.PayPeriodRecord) int {
	for _, cw := range carWashes {
		if cw.StationID == rec.StationID && cw.RecordNum == rec.RecordNumber {
			return cw.Qty.Sold
		}
	}
	return 0
}

func (pp *PayPeriod) testStation(stationID primitive.ObjectID) bool {
	cwStations := pp.getCarWashStations()
	for _, v := range cwStations {
		if v == stationID {
			return true
		}
	}
	return false
}

func (pp *PayPeriod) getCarWashStations() (stations []primitive.ObjectID) {
	// Currently there is only 1 station with car wash, creating a slice to accommodate more
	station1, _ := primitive.ObjectIDFromHex("56cf1815982d82b0f3000006")
	stations = make([]primitive.ObjectID, 1)
	stations = append(stations, station1)

	return stations
}

// setGalesLoyalty method
func (pp *PayPeriod) setGalesLoyalty() (err error) {
	docs, err := pp.db.GetGalesLoyalty(pp.dates)
	for _, s := range pp.records {
		for _, p := range docs {
			if s.StationID == p.StationID && s.RecordNumber == p.RecordNum {
				// fmt.Printf("p.Qty %+v\n", p.Qty)
				s.GalesLoyaltyQty = p.Qty.Sold
			}
		}
	}
	return err
}
