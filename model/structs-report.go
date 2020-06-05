package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// BankCardRecord struct
type BankCardRecord struct {
	BankAmex     float64
	BankDiscover float64
	BankGales    float64
	BankMC       float64
	BankVisa     float64
	CashDebit    float64
	CashOther    float64
	RecordNumber string
	StationID    primitive.ObjectID
	StationName  string
}

// EmployeeOSRecord struct
type EmployeeOSRecord struct {
	DiscrepancyDescription string
	Employee               string
	OvershortAttendant     float64
	OvershortDiff          float64
	OvershortShift         float64
	RecordNumber           string
	StationID              primitive.ObjectID
	StationName            string
}

// MonthlySaleRecord struct
type MonthlySaleRecord struct {
	BankAmex               float64
	BankDiscover           float64
	BankGales              float64
	BankMC                 float64
	BankVisa               float64
	BobsGiftCertificates   float64
	BobsNonFuelAdjustments float64
	BobsSales              float64
	CarWash                int
	CashBills              float64
	CashDebit              float64
	CashDieselDiscount     float64
	CashDriveOffNSF        float64
	CashGalesLoyaltyRedeem float64
	CashGiftCertRedeem     float64
	CashLotteryPayout      float64
	CashOSAdjusted         float64
	CashOther              float64
	CashPayout             float64
	CashWriteOff           float64
	Employee               string
	FuelSales              float64
	FuelSalesHST           float64
	FuelSalesOther         float64
	FuelSalesTotal         float64
	GalesLoyalty           int
	GiftCertificates       float64
	NonFuelSales           float64
	NonFuelTotal           float64
	ProductCigarettesQty   int
	ProductCigarettesSales float64
	ProductOilQty          int
	ProductOilSales        float64
	PropaneSales           float64
	PropaneQty             int
	RecordNumber           string
	ShiftOvershort         float64
	StationID              primitive.ObjectID
	StationName            string
}

// PayPeriodRecord struct
type PayPeriodRecord struct {
	AttendantAdjustment string
	CarwashNumber       int
	Commission          *CommissionSale
	Employee            string
	GalesLoyaltyQty     int
	NonFuelSales        float64
	ProductSales        float64
	RecordNumber        string
	ShiftOvershort      float64
	StationID           primitive.ObjectID
	StationName         string
}

// ProductNumberRecord struct
type ProductNumberRecord struct {
	Product string `bson:"_id"`
	Qty     int    `bson:"qty"`
}
