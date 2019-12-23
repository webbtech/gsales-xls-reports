package report

import (
	"github.com/pulpfree/gsales-xls-reports/model"
)

// EmployeeOS struct
type EmployeeOS struct {
	db      model.DBHandler
	dates   *model.RequestDates
	records []*model.EmployeeOSRecord
}

// ======================== Exported Methods =================================================== //

// GetRecords method
func (eo *EmployeeOS) GetRecords() ([]*model.EmployeeOSRecord, error) {

	var err error
	err = eo.setRecords()
	if err != nil {
		return nil, err
	}

	return eo.records, err
}

// ======================== Un-exported Methods ================================================ //

func (eo *EmployeeOS) setRecords() (err error) {

	sales, err := eo.db.GetEmployeeOS(eo.dates)
	if err != nil {
		return err
	}

	stationMap, err := eo.db.GetStationMap()
	if err != nil {
		return err
	}

	for _, s := range sales {
		employee, err := eo.db.GetEmployee(s.Attendant.ID)
		if err != nil {
			return err
		}

		record := &model.EmployeeOSRecord{
			DiscrepancyDescription: s.Overshort.Descrip,
			Employee:               employee,
			RecordNumber:           s.RecordNum,
			OvershortShift:         model.SetFloat(s.Overshort.Amount),
			OvershortAttendant:     model.SetFloat(s.Attendant.OvershortValue),
			OvershortDiff:          (model.SetFloat(s.Overshort.Amount) - model.SetFloat(s.Attendant.OvershortValue)),
			StationID:              s.StationID,
			StationName:            stationMap[s.StationID].Name,
		}
		eo.records = append(eo.records, record)
	}

	return err
}
