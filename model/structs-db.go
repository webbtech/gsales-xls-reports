package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ===================== Data Structs ========================================================== //

// Attendant struct
type Attendant struct {
	ID                primitive.ObjectID `bson:"ID" json:"ID"`
	Adjustment        *string            `bson:"adjustment" json:"adjustment"`
	OvershortComplete bool               `bson:"overshortComplete" json:"overshortComplete"`
	OvershortValue    *float64           `bson:"overshortValue" json:"overshortValue"`
	SheetComplete     bool               `bson:"sheetComplete" json:"sheetComplete"`
	Name              string             `bson:"name" json:"name"`
}

// Cash struct
type Cash struct {
	Bills              *float64 `bson:"bills" json:"bills"`
	Debit              *float64 `bson:"debit" json:"debit"`
	DieselDiscount     *float64 `bson:"dieselDiscount" json:"dieselDiscount"`
	DriveOffNSF        *float64 `bson:"driveOffNSF" json:"driveOffNSF"`
	GalesLoyaltyRedeem *float64 `bson:"galesLoyaltyRedeem" json:"galesLoyaltyRedeem"`
	GiftCertRedeem     *float64 `bson:"giftCertRedeem" json:"giftCertRedeem"`
	LotteryPayout      *float64 `bson:"lotteryPayout" json:"lotteryPayout"`
	Other              *float64 `bson:"other" json:"other"`
	OSAdjusted         *float64 `bson:"osAdjusted" json:"osAdjusted"`
	Payout             *float64 `bson:"payout" json:"payout"`
	WriteOff           *float64 `bson:"writeOff" json:"writeOff"`
}

// Config struct
type Config struct {
	ColouredDieselDiscount float64 `bson:"colouredDieselDsc" json:"colouredDieselDiscount"`
	Commission             int32   `bson:"commission" json:"commission"`
	DiscrepancyFlag        int32   `bson:"discrepancyFlag" json:"discrepancyFlag"`
	HST                    int32   `bson:"HST" json:"HST"`
	HiGradePremium         int32   `bson:"hiGradePremium" json:"hiGradePremium"`
}

// CreditCard struct
type CreditCard struct {
	Amex     *float64 `bson:"amex" json:"amex"`
	Discover *float64 `bson:"discover" json:"discover"`
	Gales    *float64 `bson:"gales" json:"gales"`
	MC       *float64 `bson:"mc" json:"mc"`
	Visa     *float64 `bson:"visa" json:"visa"`
}

// CommissionSale struct
type CommissionSale struct {
	Qty        int32
	Commission float64
	Sales      float64
}

// Employee struct
type Employee struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Active    bool               `bson:"active" json:"active"`
	NameFirst string             `bson:"nameFirst" json:"nameFirst"`
	NameLast  string             `bson:"nameLast" json:"nameLast"`
}

// FuelSales struct
type FuelSales struct {
	StationName  string             `bson:"stationName" json:"stationName"`
	StationSetID primitive.ObjectID `bson:"stationSetID" json:"stationSetID"`
	Fuels        []Fuel
}

// Fuel struct
type Fuel struct {
	Dollar float64 `bson:"dollars" json:"dollars"`
	Grade  int     `bson:"grade" json:"grade"`
	Litre  float64 `bson:"litres" json:"litres"`
}

// FuelType struct
type FuelType struct {
	Dollar float64 `bson:"dollar" json:"dollar"`
	Litre  float64 `bson:"litre" json:"litre"`
}

// NonFuelProduct struct
type NonFuelProduct struct {
	Qty   int     `bson:"qty"`
	Sales float64 `bson:"sales"`
	ID    struct {
		Station         primitive.ObjectID `bson:"station"`
		RecordNum       string             `bson:"recordNum"`
		ProductCategory string             `bson:"productCategory"`
	} `bson:"_id"`
}

// NonFuelSale struct
type NonFuelSale struct {
	ID        primitive.ObjectID `bson:"_id"`
	ProductID primitive.ObjectID `bson:"productID"`
	Qty       struct {
		Restock int `bson:"restock"`
		Open    int `bson:"open"`
		Sold    int `bson:"sold"`
		Close   int `bson:"close"`
	} `bson:"qty"`
	RecordNum string             `bson:"recordNum"`
	Sales     *float64           `bson:"sales"`
	StationID primitive.ObjectID `bson:"stationID"`
}

// OtherNonFuel struct
type OtherNonFuel struct {
	Bobs      *float64 `bson:"bobs" json:"bobs"`
	GiftCerts *float64 `bson:"giftCerts" json:"giftCerts"`
}

// OtherNonFuelBobs struct
type OtherNonFuelBobs struct {
	BobsGiftCerts *float64 `bson:"bobsGiftCerts" json:"bobsGiftCerts"`
}

// Overshort struct
type Overshort struct {
	Amount  float64 `bson:"amount" json:"amount"`
	Descrip string  `bson:"descrip" json:"descrip"`
}

// Product struct
type Product struct {
	Name  string
	Qty   int
	Sales float64
}

// Sales struct
type Sales struct {
	Attendant        *Attendant
	Cash             *Cash
	CreditCard       *CreditCard       `bson:"creditCard"`
	NonFuelAdjustOS  float64           `bson:"nonFuelAdjustOS"`
	OtherNonFuel     *OtherNonFuel     `bson:"otherNonFuel"`
	OtherNonFuelBobs *OtherNonFuelBobs `bson:"otherNonFuelBobs"`
	Overshort        *Overshort
	Products         []*Product
	RecordNum        string             `bson:"recordNum"`
	StationID        primitive.ObjectID `bson:"stationID"`
	Summary          *SalesSummary      `bson:"salesSummary"`
}

// SalesSummary struct
type SalesSummary struct {
	Fuel struct {
		Fuel1 *FuelType `bson:"fuel_1" json:"fuel_1"`
		Fuel2 *FuelType `bson:"fuel_2" json:"fuel_2"`
		Fuel3 *FuelType `bson:"fuel_3" json:"fuel_3"`
		Fuel4 *FuelType `bson:"fuel_4" json:"fuel_4"`
		Fuel5 *FuelType `bson:"fuel_5" json:"fuel_5"`
		Fuel6 *FuelType `bson:"fuel_6" json:"fuel_6"`
	}
	FuelAdjust             *float64 `bson:"fuelAdjust" json:"fuelAdjust"`
	FuelDollar             float64  `bson:"fuelDollar" json:"fuelDollar"`
	FuelLitre              float64  `bson:"fuelLitre" json:"fuelLitre"`
	OtherFuelDollar        float64  `bson:"otherFuelDollar" json:"otherFuelDollar"`
	OtherFuelLitre         float64  `bson:"otherFuelLitre" json:"otherFuelLitre"`
	Product                float64  `bson:"product" json:"product"`
	TotalNonFuel           float64  `bson:"totalNonFuel" json:"totalNonFuel"`
	TotalSales             float64  `bson:"totalSales" json:"totalSales"`
	TotalCash              float64  `bson:"cashTotal" json:"totalCash"`
	TotalCreditCardAndCash float64  `bson:"cashCCTotal" json:"totalCreditCardAndCash"`
	TotalCreditCard        float64  `bson:"creditCardTotal" json:"totalCreditCard"`
}

// Station struct
type Station struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

// StationNodes struct
type StationNodes struct {
	ID    primitive.ObjectID   `bson:"_id"`
	Name  string               `bson:"name"`
	Nodes []primitive.ObjectID `bson:"nodes"`
}
