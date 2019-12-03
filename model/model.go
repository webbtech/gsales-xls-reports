package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// DBHandler interface
type DBHandler interface {
	Close()
	GetConfig() (*Config, error)
	GetEmployee(primitive.ObjectID) (string, error)
	GetMonthlyProducts(*RequestDates) ([]*NonFuelProduct, error)
	GetMonthlySales(*RequestDates) ([]*Sales, error)
	GetPayPeriodSales(*RequestDates) ([]*Sales, error)
	GetStationMap() (map[primitive.ObjectID]*Station, error)
	GetNonFuelCommission(string, primitive.ObjectID) (*CommissionSale, error)
	GetCarWash(*RequestDates) ([]*NonFuelSale, error)
}

// Report interface
type Report interface {
	// Init(db DBHandler, dates *RequestDates)
	// GetRecords()
	GetRecords() (interface{}, error)
}

// ===================== Helper Functions ====================================================== //

// SetFloat function
func SetFloat(num interface{}) float64 {

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
}
