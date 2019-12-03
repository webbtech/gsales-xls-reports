package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// MonthlySales struct
type MonthlySales struct {
	BankAmex               float64
	BankDiscover           float64
	BankGales              float64
	BankMC                 float64
	BankVisa               float64
	BobsGiftCertificates   float64
	BobsNonFuelAdjustments float64
	BobsSales              float64
	CashBills              float64
	CashDebit              float64
	CashDieselDiscount     float64
	CashDriveOffNSF        float64
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
	GiftCertificates       float64
	NonFuelTotal           float64
	ProductCigarettesQty   int
	ProductCigarettesSales float64
	ProductOilQty          int
	ProductOilSales        float64
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
	NonFuelSales        float64
	ProductSales        float64
	RecordNumber        string
	ShiftOvershort      float64
	StationID           primitive.ObjectID
	StationName         string
}
