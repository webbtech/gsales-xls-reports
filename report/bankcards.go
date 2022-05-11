package report

import (
	"github.com/webbtech/gsales-xls-reports/model"
)

// BankCard struct
type BankCard struct {
	db      model.DBHandler
	dates   *model.RequestDates
	records []*model.BankCardRecord
}

// ======================== Exported Methods =================================================== //

// GetRecords method
func (bc *BankCard) GetRecords() ([]*model.BankCardRecord, error) {

	var err error
	err = bc.setRecords()
	if err != nil {
		return nil, err
	}

	return bc.records, err
}

// ======================== Un-exported Methods ================================================ //

func (bc *BankCard) setRecords() (err error) {

	sales, err := bc.db.GetBankCards(bc.dates)
	if err != nil {
		return err
	}

	stationMap, err := bc.db.GetStationMap()
	if err != nil {
		return err
	}

	for _, s := range sales {
		record := &model.BankCardRecord{
			BankAmex:     model.SetFloat(s.CreditCard.Amex),
			BankDiscover: model.SetFloat(s.CreditCard.Discover),
			BankGales:    model.SetFloat(s.CreditCard.Gales),
			BankMC:       model.SetFloat(s.CreditCard.MC),
			BankVisa:     model.SetFloat(s.CreditCard.Visa),
			CashDebit:    model.SetFloat(s.Cash.Debit),
			CashOther:    model.SetFloat(s.Cash.Other),
			RecordNumber: s.RecordNum,
			StationID:    s.StationID,
			StationName:  stationMap[s.StationID].Name,
		}
		bc.records = append(bc.records, record)
	}

	return err
}
