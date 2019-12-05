package payperiod

import (
	"github.com/pulpfree/gsales-xls-reports/config"
	"github.com/pulpfree/gsales-xls-reports/model"
	"github.com/pulpfree/gsales-xls-reports/model/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Records struct
type Records struct {
	cfg     *model.Config
	DB      model.DBHandler
	dates   *model.RequestDates
	Records []*model.PayPeriodRecord
}

// ======================== Exported Functions ================================================= //

// Init function
func Init(dates *model.RequestDates, config *config.Config) (*Records, error) {

	var err error
	db, err := db.NewDB(config.GetMongoConnectURL(), config.DBName)
	if err != nil {
		return nil, err
	}

	cfg, err := db.GetConfig()
	if err != nil {
		return nil, err
	}

	return &Records{
		cfg:   cfg,
		dates: dates,
		DB:    db,
	}, err
}

// ======================== Exported Methods =================================================== //

// GetRecords method
func (pp *Records) GetRecords() ([]*model.PayPeriodRecord, error) {

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

	return pp.Records, err
}

// ======================== Un-exported Methods ================================================ //

func (pp *Records) setRecords() (err error) {

	sales, err := pp.DB.GetPayPeriodSales(pp.dates)
	if err != nil {
		return err
	}

	stationMap, err := pp.DB.GetStationMap()
	if err != nil {
		return err
	}

	for _, s := range sales {

		employee, err := pp.DB.GetEmployee(s.Attendant.ID)
		if err != nil {
			return err
		}

		record := &model.PayPeriodRecord{
			AttendantAdjustment: s.Attendant.Adjustment,
			Employee:            employee,
			NonFuelSales:        model.SetFloat(s.Summary.TotalNonFuel),
			ProductSales:        model.SetFloat(s.Summary.Product),
			RecordNumber:        s.RecordNum,
			StationID:           s.StationID,
			StationName:         stationMap[s.StationID].Name,
			ShiftOvershort:      s.Overshort.Amount,
		}
		pp.Records = append(pp.Records, record)
	}

	return err
}

// setNonFuelCommission method
func (pp *Records) setNonFuelCommission() (err error) {
	for _, s := range pp.Records {
		commission, err := pp.DB.GetNonFuelCommission(s.RecordNumber, s.StationID)
		if err != nil {
			return err
		}
		s.Commission = commission
	}
	return err
}

// setCarWashes method
func (pp *Records) setCarWashes() (err error) {
	washes, err := pp.DB.GetCarWash(pp.dates)
	if err != nil {
		return err
	}
	for _, s := range pp.Records {
		if pp.testStation(s.StationID) {
			s.CarwashNumber = pp.searchCarWashSale(washes, s)
		}
	}
	return err
}

func (pp *Records) searchCarWashSale(carWashes []*model.NonFuelSale, rec *model.PayPeriodRecord) int {
	for _, cw := range carWashes {
		if cw.StationID == rec.StationID && cw.RecordNum == rec.RecordNumber {
			return cw.Qty.Sold
		}
	}
	return 0
}

func (pp *Records) testStation(stationID primitive.ObjectID) bool {
	cwStations := pp.getCarWashStations()
	for _, v := range cwStations {
		if v == stationID {
			return true
		}
	}
	return false
}

func (pp *Records) getCarWashStations() (stations []primitive.ObjectID) {
	// Currently there is only 1 station with car wash, creating a slice to accommodate more
	station1, _ := primitive.ObjectIDFromHex("56cf1815982d82b0f3000006")
	stations = make([]primitive.ObjectID, 1)
	stations = append(stations, station1)

	return stations
}
