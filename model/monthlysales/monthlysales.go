package monthlysales

import (
	"github.com/pulpfree/gsales-xls-reports/config"
	"github.com/pulpfree/gsales-xls-reports/model"
	"github.com/pulpfree/gsales-xls-reports/model/db"
)

// Sales struct
type Sales struct {
	cfg   *model.Config
	dates *model.RequestDates
	DB    model.DBHandler
	Sales []*model.MonthlySales
}

// ======================== Exported Functions ================================================= //

// Init function
func Init(dates *model.RequestDates, config *config.Config) (*Sales, error) {

	var err error
	db, err := db.NewDB(config.GetMongoConnectURL(), config.DBName)
	if err != nil {
		return nil, err
	}

	cfg, err := db.GetConfig()
	if err != nil {
		return nil, err
	}

	return &Sales{
		cfg:   cfg,
		dates: dates,
		DB:    db,
	}, err
}

// ======================== Exported Methods =================================================== //

// GetRecords method
func (ms *Sales) GetRecords() ([]*model.MonthlySales, error) {

	var err error

	err = ms.setRecords()
	if err != nil {
		return nil, err
	}

	err = ms.setMonthlySalesProducts()
	return ms.Sales, err
}

// ======================== Un-exported Methods ================================================ //

func (ms *Sales) setRecords() (err error) {

	sales, err := ms.DB.GetMonthlySales(ms.dates)
	if err != nil {
		return err
	}

	hst := (float64(ms.cfg.HST)/100 + 1)
	stationMap, err := ms.DB.GetStationMap()
	if err != nil {
		return err
	}

	for _, s := range sales {

		employee, err := ms.DB.GetEmployee(s.Attendant.ID)
		if err != nil {
			return err
		}

		// set HST
		fuelSalesNoHST := s.Summary.FuelDollar / hst
		fuelSalesHST := s.Summary.FuelDollar - fuelSalesNoHST

		record := &model.MonthlySales{
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
		ms.Sales = append(ms.Sales, record)
	}

	return err
}

//

// examples see: https://github.com/simagix/mongo-go-examples/blob/master/examples/aggregate_array_test.go
/* func (db *MDB) fetchNonFuel(startDate, endDate time.Time) (docs []*productDoc, err error) {

	col := db.db.Collection(colNonFuelSales)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{
			primitive.E{
				Key: "$match",
				Value: bson.D{
					primitive.E{
						Key: "recordDate",
						Value: bson.D{
							primitive.E{
								Key:   "$gte",
								Value: startDate,
							},
							primitive.E{
								Key:   "$lte",
								Value: endDate,
							},
						},
					},
				},
			},
		},
		{
			primitive.E{
				Key: "$lookup",
				Value: bson.D{
					primitive.E{
						Key:   "from",
						Value: colProducts,
					},
					primitive.E{
						Key:   "localField",
						Value: "productID",
					},
					primitive.E{
						Key:   "foreignField",
						Value: "_id",
					},
					primitive.E{
						Key:   "as",
						Value: "product",
					},
				},
			},
		},
		{
			primitive.E{
				Key:   "$unwind",
				Value: "$product",
			},
		},
		{
			primitive.E{
				Key: "$group",
				Value: bson.D{
					primitive.E{
						Key: "_id",
						Value: bson.D{
							primitive.E{
								Key:   "recordNum",
								Value: "$recordNum",
							},
							primitive.E{
								Key:   "station",
								Value: "$stationID",
							},
							primitive.E{
								Key:   "product",
								Value: "$product.category",
							},
						},
					},
					primitive.E{
						Key: "sales",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$sales",
							},
						},
					},
					primitive.E{
						Key: "qty",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$qty.sold",
							},
						},
					},
				},
			},
		},
		{
			primitive.E{
				Key: "$sort",
				Value: bson.D{
					primitive.E{
						Key:   "_id.recordNum",
						Value: 1,
					},
				},
			},
		},
	}

	cur, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var doc productDoc
		cur.Decode(&doc)
		docs = append(docs, &doc)
	}

	return docs, nil
} */

// setSalesProducts
// set values for cigarettes and oil category products
func (ms *Sales) setMonthlySalesProducts() (err error) {

	products, err := ms.DB.GetMonthlyProducts(ms.dates)
	if err != nil {
		return err
	}
	for _, s := range ms.Sales {
		for _, p := range products {
			if s.StationID == p.ID.Station && s.RecordNumber == p.ID.RecordNum {
				if p.ID.Product == "cigarettes" {
					s.ProductCigarettesQty = p.Qty
					s.ProductCigarettesSales = p.Sales
				}
				if p.ID.Product == "oil" {
					s.ProductOilQty = p.Qty
					s.ProductOilSales = p.Sales
				}
			}
		}
	}

	return err
}
