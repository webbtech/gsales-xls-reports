package report

import (
	"github.com/pulpfree/gsales-xls-reports/model"
)

// MonthlySales struct
type MonthlySales struct {
	cfg     *model.Config
	db      model.DBHandler
	dates   *model.RequestDates
	records []*model.MonthlySaleRecord
}

// ======================== Exported Methods =================================================== //

// GetRecords method
func (ms *MonthlySales) GetRecords() ([]*model.MonthlySaleRecord, error) {

	var err error

	err = ms.setRecords()
	if err != nil {
		return nil, err
	}

	err = ms.setSalesProducts()
	if err != nil {
		return nil, err
	}

	err = ms.setCarWashes()
	if err != nil {
		return nil, err
	}

	err = ms.setGalesLoyalty()

	return ms.records, err
}

// ======================== Un-exported Methods ================================================ //

func (ms *MonthlySales) setRecords() (err error) {

	sales, err := ms.db.GetMonthlySales(ms.dates)
	if err != nil {
		return err
	}

	hst := (float64(ms.cfg.HST)/100 + 1)
	stationMap, err := ms.db.GetStationMap()
	if err != nil {
		return err
	}

	for _, s := range sales {

		employee, err := ms.db.GetEmployee(s.Attendant.ID)
		if err != nil {
			return err
		}

		// set HST
		fuelSalesNoHST := s.Summary.FuelDollar / hst
		fuelSalesHST := s.Summary.FuelDollar - fuelSalesNoHST

		record := &model.MonthlySaleRecord{
			BankAmex:               model.SetFloat(s.CreditCard.Amex),
			BankDiscover:           model.SetFloat(s.CreditCard.Discover),
			BankGales:              model.SetFloat(s.CreditCard.Gales),
			BankMC:                 model.SetFloat(s.CreditCard.MC),
			BankVisa:               model.SetFloat(s.CreditCard.Visa),
			BobsGiftCertificates:   model.SetFloat(s.OtherNonFuelBobs.BobsGiftCerts),
			BobsNonFuelAdjustments: model.SetFloat(s.Summary.BobsFuelAdj),
			BobsSales:              model.SetFloat(s.OtherNonFuel.Bobs),
			CashBills:              model.SetFloat(s.Cash.Bills),
			CashDebit:              model.SetFloat(s.Cash.Debit),
			CashDieselDiscount:     model.SetFloat(s.Cash.DieselDiscount),
			CashDriveOffNSF:        model.SetFloat(s.Cash.DriveOffNSF),
			CashGalesLoyaltyRedeem: model.SetFloat(s.Cash.GalesLoyaltyRedeem),
			CashGiftCertRedeem:     model.SetFloat(s.Cash.GiftCertRedeem),
			CashLotteryPayout:      model.SetFloat(s.Cash.LotteryPayout),
			CashOSAdjusted:         model.SetFloat(s.Cash.OSAdjusted),
			CashOther:              model.SetFloat(s.Cash.Other),
			CashPayout:             model.SetFloat(s.Cash.Payout),
			CashWriteOff:           model.SetFloat(s.Cash.WriteOff),
			Employee:               employee,
			FuelSales:              fuelSalesNoHST,
			FuelSalesHST:           fuelSalesHST,
			FuelSalesOther:         model.SetFloat(s.Summary.OtherFuelDollar),
			FuelSalesTotal:         model.SetFloat(s.Summary.FuelDollar),
			GiftCertificates:       model.SetFloat(s.OtherNonFuel.GiftCerts),
			NonFuelTotal:           model.SetFloat(s.Summary.TotalNonFuel),
			RecordNumber:           s.RecordNum,
			ShiftOvershort:         model.SetFloat(s.Overshort.Amount),
			StationID:              s.StationID,
			StationName:            stationMap[s.StationID].Name,
		}
		ms.records = append(ms.records, record)
	}

	return err
}

// setSalesProducts method
// set values for cigarettes and oil category products
func (ms *MonthlySales) setSalesProducts() (err error) {

	products, err := ms.db.GetMonthlyProducts(ms.dates)
	if err != nil {
		return err
	}
	for _, s := range ms.records {
		for _, p := range products {
			if s.StationID == p.ID.Station && s.RecordNumber == p.ID.RecordNum {
				if p.ID.ProductCategory == "cigarettes" {
					s.ProductCigarettesQty = p.Qty
					s.ProductCigarettesSales = p.Sales
				}
				if p.ID.ProductCategory == "oil" {
					s.ProductOilQty = p.Qty
					s.ProductOilSales = p.Sales
				}
			}
		}
	}

	return err
}

// setCarWashes method
func (ms *MonthlySales) setCarWashes() (err error) {
	products, err := ms.db.GetCarWash(ms.dates)
	for _, s := range ms.records {
		for _, p := range products {
			if s.StationID == p.StationID && s.RecordNumber == p.RecordNum {
				// fmt.Printf("p.Qty %+v\n", p.Qty)
				s.CarWash = p.Qty.Sold
			}
		}
	}
	return err
}

// setGalesLoyalty method
func (ms *MonthlySales) setGalesLoyalty() (err error) {
	docs, err := ms.db.GetMonthlyGalesLoyalty(ms.dates)
	for _, s := range ms.records {
		for _, p := range docs {
			if s.StationID == p.StationID && s.RecordNumber == p.RecordNum {
				// fmt.Printf("p.Qty %+v\n", p.Qty)
				s.GalesLoyalty = p.Qty.Sold
			}
		}
	}
	return err
}
