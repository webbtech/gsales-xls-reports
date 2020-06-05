package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	log "github.com/sirupsen/logrus"

	"github.com/pulpfree/gsales-xls-reports/model"
	"github.com/pulpfree/gsales-xls-reports/pkgerrors"
)

// MDB struct
type MDB struct {
	client      *mongo.Client
	cfg         *model.Config
	dbName      string
	db          *mongo.Database
	employeeMap map[primitive.ObjectID]string
	stationMap  map[primitive.ObjectID]*model.Station
}

// DB and Table constants
const (
	colConfig       = "config"
	colEmployees    = "employees"
	colNonFuelSales = "non-fuel-sales"
	colProducts     = "products"
	colSales        = "sales"
	colStations     = "stations"
)

// misc constants
const (
	noRecordsMsg     = "No records found matching criteria"
	carWashProductID = "574707198bba4f0100582b83"
)

// slice of Gales Loyalty product ids
var galesLoyaltyProductIDs = []string{"5e2080472dbbd30008721739", "5e1f63167934140007cc6c98"}

// slice of Propane Tanks type product ids
var propaneTankProductIDs = []string{"56cf4bfe982d82b41d000019", "56cf4bfe982d82b41d00001a"}

// ======================== Exported Functions ================================================= //

// NewDB sets up new MDB struct
func NewDB(connection string, dbNm string) (*MDB, error) {

	clientOptions := options.Client().ApplyURI(connection)
	err := clientOptions.Validate()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	return &MDB{
		client: client,
		dbName: dbNm,
		db:     client.Database(dbNm),
	}, err
}

// ======================== Exported Methods =================================================== //

// Close method
func (db *MDB) Close() {
	err := db.client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

// GetBankCards method
func (db *MDB) GetBankCards(dates *model.RequestDates) (records []*model.Sales, err error) {

	records, err = db.fetchBankCards(dates.DateFrom, dates.DateTo)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, &pkgerrors.MongoError{Err: "", Caller: "db.GetBankCards", Msg: noRecordsMsg}
	}

	return records, err
}

// GetConfig method
func (db *MDB) GetConfig() (cfg *model.Config, err error) {

	if db.cfg != nil {
		return db.cfg, nil
	}

	err = db.setConfig()
	return db.cfg, err
}

// GetEmployee method
func (db *MDB) GetEmployee(employeeID primitive.ObjectID) (employee string, err error) {
	employee, err = db.setEmployee(employeeID)
	return employee, err
}

// GetStationMap method
func (db *MDB) GetStationMap() (map[primitive.ObjectID]*model.Station, error) {

	if len(db.stationMap) > 0 {
		return db.stationMap, nil
	}

	err := db.setStationMap()
	return db.stationMap, err
}

// GetEmployeeOS method
func (db *MDB) GetEmployeeOS(dates *model.RequestDates) (records []*model.Sales, err error) {
	records, err = db.fetchEmployeeOS(dates.DateFrom, dates.DateTo)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, &pkgerrors.MongoError{Err: "", Caller: "db.GetEmployeeOS", Msg: noRecordsMsg}
	}

	return records, err
}

// GetGalesLoyalty method
func (db *MDB) GetGalesLoyalty(dates *model.RequestDates) (records []*model.NonFuelSale, err error) {
	records, err = db.fetchGalesLoyalty(dates.DateFrom, dates.DateTo)
	if err != nil {
		return nil, err
	}

	return records, err
}

// GetMonthlySales method
func (db *MDB) GetMonthlySales(dates *model.RequestDates) (records []*model.Sales, err error) {
	records, err = db.fetchMonthlySales(dates.DateFrom, dates.DateTo)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, &pkgerrors.MongoError{Err: "", Caller: "db.GetMonthlySales", Msg: noRecordsMsg}
	}

	return records, err
}

// GetMonthlyProducts method
func (db *MDB) GetMonthlyProducts(dates *model.RequestDates) (records []*model.NonFuelProduct, err error) {

	records, err = db.fetchMonthlyNonFuel(dates.DateFrom, dates.DateTo)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, &pkgerrors.MongoError{Err: "", Caller: "db.GetMonthlyProducts", Msg: noRecordsMsg}
	}

	return records, err
}

// GetNonFuelSales method
func (db *MDB) GetNonFuelSales(dates *model.RequestDates) (records []*model.NonFuelProduct, err error) {
	records, err = db.fetchNonFuelSales(dates.DateFrom, dates.DateTo)
	if err != nil {
		return nil, err
	}

	return records, err
}

// GetPayPeriodSales method
func (db *MDB) GetPayPeriodSales(dates *model.RequestDates) (records []*model.Sales, err error) {
	records, err = db.fetchPayPeriodSales(dates.DateFrom, dates.DateTo)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, &pkgerrors.MongoError{Err: "", Caller: "db.GetPayPeriodSales", Msg: noRecordsMsg}
	}

	return records, err
}

// GetProductNumbers method
func (db *MDB) GetProductNumbers(dates *model.RequestDates) (records []*model.ProductNumberRecord, err error) {
	records, err = db.fetchProductNumbers(dates.DateFrom, dates.DateTo)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, &pkgerrors.MongoError{Err: "", Caller: "db.GetProductNumbers", Msg: noRecordsMsg}
	}

	return records, err
}

// GetPropaneSales method
func (db *MDB) GetPropaneSales(dates *model.RequestDates) (records []*model.NonFuelProduct, err error) {
	records, err = db.fetchPropaneProducts(dates.DateFrom, dates.DateTo)
	if err != nil {
		return nil, err
	}

	return records, err
}

// GetCarWash method
func (db *MDB) GetCarWash(dates *model.RequestDates) (records []*model.NonFuelSale, err error) {
	records, err = db.fetchCarWash(dates.DateFrom, dates.DateTo)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, &pkgerrors.MongoError{Err: "", Caller: "db.GetCarWash", Msg: noRecordsMsg}
	}

	return records, err
}

// GetNonFuelCommission method
func (db *MDB) GetNonFuelCommission(recordNum string, stationID primitive.ObjectID) (com *model.CommissionSale, err error) {
	com, err = db.fetchNonFuelCommission(recordNum, stationID)
	return com, err
}

// ======================== Un-exported Methods ================================================ //

// setConfig method
func (db *MDB) setConfig() (err error) {

	if db.cfg != nil {
		return nil
	}
	col := db.db.Collection(colConfig)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{primitive.E{Key: "_id", Value: 1}}
	err = col.FindOne(ctx, filter).Decode(&db.cfg)

	return err
}

// setEmployee method
func (db *MDB) setEmployee(employeeID primitive.ObjectID) (employee string, err error) {

	if db.employeeMap[employeeID] != "" {
		return db.employeeMap[employeeID], nil
	}

	db.employeeMap = make(map[primitive.ObjectID]string)

	col := db.db.Collection(colEmployees)
	emp := new(model.Employee)
	filter := bson.D{primitive.E{Key: "_id", Value: employeeID}}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = col.FindOne(ctx, filter).Decode(&emp)
	if err != nil {
		return employee, err
	}

	db.employeeMap[emp.ID] = fmt.Sprintf("%s, %s", emp.NameLast, emp.NameFirst)
	return db.employeeMap[emp.ID], nil
}

// setStationMap method
func (db *MDB) setStationMap() (err error) {

	if len(db.stationMap) > 0 {
		fmt.Printf("db.stationMap has len %+v\n", db.stationMap)
		return nil
	}
	db.stationMap = make(map[primitive.ObjectID]*model.Station)

	col := db.db.Collection(colStations)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := col.Find(ctx, bson.D{{}})
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result model.Station
		err := cur.Decode(&result)
		if err != nil {
			return err
		}

		db.stationMap[result.ID] = &result
	}
	return err
}

// fetchBankCards method
func (db *MDB) fetchBankCards(startDate, endDate time.Time) (sales []*model.Sales, err error) {

	col := db.db.Collection(colSales)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{primitive.E{Key: "stationID", Value: 1}, primitive.E{Key: "recordNum", Value: 1}})
	filter := bson.M{"recordDate": bson.M{"$gte": startDate, "$lte": endDate}}
	cur, err := col.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &sales); err != nil {
		return nil, err
	}

	return sales, err
}

// fetchEmployeeOS method
func (db *MDB) fetchEmployeeOS(startDate, endDate time.Time) (sales []*model.Sales, err error) {

	col := db.db.Collection(colSales)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{primitive.E{Key: "attendant.ID", Value: 1}, primitive.E{Key: "recordNum", Value: 1}})
	filter := bson.M{"recordDate": bson.M{"$gte": startDate, "$lte": endDate}}
	cur, err := col.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &sales); err != nil {
		return nil, err
	}

	return sales, err
}

func (db *MDB) fetchGalesLoyalty(startDate, endDate time.Time) (docs []*model.NonFuelSale, err error) {

	col := db.db.Collection(colNonFuelSales)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pIDs := make([]primitive.ObjectID, len(galesLoyaltyProductIDs))
	for i := range galesLoyaltyProductIDs {
		pIDs[i], _ = primitive.ObjectIDFromHex(galesLoyaltyProductIDs[i])
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{primitive.E{Key: "stationID", Value: 1}, primitive.E{Key: "recordNum", Value: 1}})
	filter := bson.M{"productID": bson.M{"$in": pIDs}, "recordDate": bson.M{"$gte": startDate, "$lte": endDate}}
	cur, err := col.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &docs); err != nil {
		return nil, err
	}

	return docs, err
}

// examples see: https://github.com/simagix/mongo-go-examples/blob/master/examples/aggregate_array_test.go
func (db *MDB) fetchMonthlyNonFuel(startDate, endDate time.Time) (docs []*model.NonFuelProduct, err error) {

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
								Key:   "productCategory",
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
		var doc model.NonFuelProduct
		cur.Decode(&doc)
		docs = append(docs, &doc)
	}

	return docs, nil
}

func (db *MDB) fetchMonthlySales(startDate, endDate time.Time) (sales []*model.Sales, err error) {
	col := db.db.Collection(colSales)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{primitive.E{Key: "stationID", Value: 1}})
	filter := bson.M{"recordDate": bson.M{"$gte": startDate, "$lte": endDate}}
	cur, err := col.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &sales); err != nil {
		return nil, err
	}

	return sales, err
}

func (db *MDB) fetchNonFuelSales(startDate, endDate time.Time) (docs []*model.NonFuelProduct, err error) {

	col := db.db.Collection(colNonFuelSales)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pIDs := make([]primitive.ObjectID, len(propaneTankProductIDs))
	for i := range propaneTankProductIDs {
		pIDs[i], _ = primitive.ObjectIDFromHex(propaneTankProductIDs[i])
	}

	pipeline := mongo.Pipeline{
		{
			primitive.E{
				Key: "$match",
				Value: bson.D{
					primitive.E{
						Key:   "productID",
						Value: bson.M{"$nin": pIDs},
					},
					primitive.E{
						Key:   "sales",
						Value: bson.M{"$gt": 0},
					},
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
	}

	cur, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &docs); err != nil {
		return nil, err
	}

	return docs, err
}

func (db *MDB) fetchPayPeriodSales(startDate, endDate time.Time) (sales []*model.Sales, err error) {

	col := db.db.Collection(colSales)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{primitive.E{Key: "attendant.ID", Value: 1}, primitive.E{Key: "recordNum", Value: 1}})
	filter := bson.M{"recordDate": bson.M{"$gte": startDate, "$lte": endDate}}
	cur, err := col.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &sales); err != nil {
		return nil, err
	}

	return sales, err
}

func (db *MDB) fetchProductNumbers(startDate, endDate time.Time) (products []*model.ProductNumberRecord, err error) {

	col := db.db.Collection(colNonFuelSales)
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
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
						},
					},
					primitive.E{
						Key: "recordDate",
						Value: bson.D{
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
				Key: "$match",
				Value: bson.D{
					primitive.E{
						Key:   "product.type",
						Value: "report1",
					},
				},
			},
		},
		{
			primitive.E{
				Key: "$group",
				Value: bson.D{
					primitive.E{
						Key:   "_id",
						Value: "$product.name",
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
						Key:   "_id",
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
		var record model.ProductNumberRecord
		cur.Decode(&record)
		products = append(products, &record)
	}

	return products, err
}

func (db *MDB) fetchPropaneProducts(startDate, endDate time.Time) (docs []*model.NonFuelProduct, err error) {

	col := db.db.Collection(colNonFuelSales)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pIDs := make([]primitive.ObjectID, len(propaneTankProductIDs))
	for i := range propaneTankProductIDs {
		pIDs[i], _ = primitive.ObjectIDFromHex(propaneTankProductIDs[i])
	}

	pipeline := mongo.Pipeline{
		{
			primitive.E{
				Key: "$match",
				Value: bson.D{
					primitive.E{
						Key:   "productID",
						Value: bson.M{"$in": pIDs},
					},
					primitive.E{
						Key:   "sales",
						Value: bson.M{"$gt": 0},
					},
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
	}

	cur, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &docs); err != nil {
		return nil, err
	}

	return docs, err
}

func (db *MDB) fetchCarWash(startDate, endDate time.Time) (sales []*model.NonFuelSale, err error) {

	col := db.db.Collection(colNonFuelSales)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cwProductID, _ := primitive.ObjectIDFromHex(carWashProductID)
	findOptions := options.Find()
	findOptions.SetSort(bson.D{primitive.E{Key: "stationID", Value: 1}, primitive.E{Key: "recordNum", Value: 1}})
	filter := bson.M{"productID": cwProductID, "recordDate": bson.M{"$gte": startDate, "$lte": endDate}}
	cur, err := col.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &sales); err != nil {
		return nil, err
	}

	return sales, err
}

func (db *MDB) fetchNonFuelCommission(recordNum string, stationID primitive.ObjectID) (com *model.CommissionSale, err error) {

	cfg, _ := db.GetConfig()
	hst := (float64(cfg.HST)/100 + 1)
	comm := (float64(cfg.Commission) / 100)

	var saleF float64
	com = new(model.CommissionSale)

	col := db.db.Collection(colNonFuelSales)
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{
			primitive.E{
				Key: "$match",
				Value: bson.D{
					primitive.E{
						Key:   "recordNum",
						Value: recordNum,
					},
					primitive.E{
						Key:   "stationID",
						Value: stationID,
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
				Key: "$match",
				Value: bson.D{
					primitive.E{
						Key: "product.commissionEligible",
						Value: bson.D{
							primitive.E{
								Key:   "$eq",
								Value: true,
							},
						},
					},
					primitive.E{
						Key: "sales",
						Value: bson.D{
							primitive.E{
								Key:   "$gt",
								Value: 0,
							},
						},
					},
				},
			},
		},
		{
			primitive.E{
				Key: "$group",
				Value: bson.D{
					primitive.E{
						Key:   "_id",
						Value: nil,
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
						Key: "numSold",
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
	}

	cur, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []bson.M
	var doc bson.M
	if err = cur.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	// as there is only one result, we need to extract
	if len(results) > 0 {
		doc = results[0]
	}

	if doc["sales"] != nil {
		switch v := doc["sales"].(type) {
		case int:
			saleF = float64(doc["sales"].(int))
		case int32:
			saleF = float64(doc["sales"].(int32))
		case float64:
			saleF = doc["sales"].(float64)
		default:
			fmt.Printf("unexpected type in db: %T val %v\n", v, doc["sales"])
		}
		com.Qty = doc["numSold"].(int32)
		com.Sales = saleF
		com.Commission = (saleF / hst * comm)
	}

	return com, nil
}
