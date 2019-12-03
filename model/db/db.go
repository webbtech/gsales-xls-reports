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

// GetMonthlySales method
// TODO: create un-exported method to do actual fetch
func (db *MDB) GetMonthlySales(dates *model.RequestDates) (sales []*model.Sales, err error) {

	col := db.db.Collection(colSales)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{primitive.E{Key: "stationID", Value: 1}})
	filter := bson.M{"recordDate": bson.M{"$gte": dates.DateFrom, "$lte": dates.DateTo}}
	cur, err := col.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result model.Sales
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		sales = append(sales, &result)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return sales, err
}

// GetMonthlyProducts method
func (db *MDB) GetMonthlyProducts(dates *model.RequestDates) (docs []*model.NonFuelProduct, err error) {
	docs, err = db.fetchMonthlyNonFuel(dates.DateFrom, dates.DateTo)
	return docs, err
}

// GetPayPeriodSales method
func (db *MDB) GetPayPeriodSales(dates *model.RequestDates) (sales []*model.Sales, err error) {
	sales, err = db.fetchPayPeriodSales(dates.DateFrom, dates.DateTo)
	return sales, err
}

// GetCarWash method
func (db *MDB) GetCarWash(dates *model.RequestDates) (nfSales []*model.NonFuelSale, err error) {
	nfSales, err = db.fetchCarWash(dates.DateFrom, dates.DateTo)
	return nfSales, err
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
		var doc model.NonFuelProduct
		cur.Decode(&doc)
		docs = append(docs, &doc)
	}

	return docs, nil
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

	for cur.Next(ctx) {
		var result model.Sales
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		sales = append(sales, &result)
		// fmt.Printf("result %+v\n", result.Summary.Product)
	}
	return sales, err
}

func (db *MDB) fetchCarWash(startDate, endDate time.Time) (nfSales []*model.NonFuelSale, err error) {

	col := db.db.Collection(colNonFuelSales)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{primitive.E{Key: "stationID", Value: 1}, primitive.E{Key: "recordNum", Value: 1}})
	filter := bson.M{"recordDate": bson.M{"$gte": startDate, "$lte": endDate}}
	cur, err := col.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result model.NonFuelSale
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		nfSales = append(nfSales, &result)
	}

	return nfSales, err
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
