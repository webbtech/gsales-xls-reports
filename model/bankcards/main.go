package bankcards

import (
	"github.com/pulpfree/gsales-xls-reports/config"
	"github.com/pulpfree/gsales-xls-reports/model"
	"github.com/pulpfree/gsales-xls-reports/model/db"
)

// Records struct
type Records struct {
	cfg     *model.Config
	DB      model.DBHandler
	dates   *model.RequestDates
	Records []*model.BankCardRecord
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
func (bc *Records) GetRecords() ([]*model.BankCardRecord, error) {

	var err error
	err = bc.setRecords()
	if err != nil {
		return nil, err
	}

	return bc.Records, err
}

// ======================== Un-exported Methods ================================================ //

func (bc *Records) setRecords() (err error) {

	sales, err := bc.DB.GetBankCards(bc.dates)
	if err != nil {
		return err
	}

	stationMap, err := bc.DB.GetStationMap()
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
		bc.Records = append(bc.Records, record)
	}

	return err
}
