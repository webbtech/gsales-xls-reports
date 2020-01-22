package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// DBHandler interface
type DBHandler interface {
	Close()
	GetBankCards(*RequestDates) ([]*Sales, error)
	GetCarWash(*RequestDates) ([]*NonFuelSale, error)
	GetConfig() (*Config, error)
	GetEmployee(primitive.ObjectID) (string, error)
	GetEmployeeOS(*RequestDates) ([]*Sales, error)
	GetMonthlyGalesLoyalty(*RequestDates) ([]*NonFuelSale, error)
	GetMonthlyProducts(*RequestDates) ([]*NonFuelProduct, error)
	GetMonthlySales(*RequestDates) ([]*Sales, error)
	GetNonFuelCommission(string, primitive.ObjectID) (*CommissionSale, error)
	GetPayPeriodSales(*RequestDates) ([]*Sales, error)
	GetProductNumbers(*RequestDates) ([]*ProductNumberRecord, error)
	GetStationMap() (map[primitive.ObjectID]*Station, error)
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

// SetString function
func SetString(s interface{}) string {
	var ret string
	switch v := s.(type) {
	case *string:
		// need to check for nil here to deal with null db values
		if v == nil {
			ret = ""
		} else {
			ret = *v
		}
	case string:
		ret = v
	default:
		ret = ""
	}

	return ret
}
