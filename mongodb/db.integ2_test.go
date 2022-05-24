package mongodb

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	pkgerrors "github.com/pulpfree/go-errors"
	"github.com/stretchr/testify/suite"
	"github.com/webbtech/gsales-xls-reports/config"
	"github.com/webbtech/gsales-xls-reports/handlers"
	"github.com/webbtech/gsales-xls-reports/model"
	"github.com/webbtech/gsales-xls-reports/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// These test comprise primarily the un-exported methods

// MongoSuite2 struct
type MongoSuite2 struct {
	cfg       *config.Config
	dateMonth *model.RequestDates
	dateDays  *model.RequestDates
	db        *MDB
	suite.Suite
	rpt *handlers.Report
}

// SetupTest method
func (suite *MongoSuite2) SetupTest() {
	// setup config
	os.Setenv("Stage", "test")
	suite.cfg = &config.Config{}
	err := suite.cfg.Init()
	suite.NoError(err)

	// Set client options
	clientOptions := options.Client().ApplyURI(suite.cfg.GetMongoConnectURL())

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	suite.NoError(err)

	suite.db = &MDB{
		client: client,
		dbName: suite.cfg.DbName,
		db:     client.Database(suite.cfg.DbName),
	}

	startDte, endDte, _ := utils.DatesFromMonth(dateMonth)
	suite.dateMonth = &model.RequestDates{
		DateFrom: startDte,
		DateTo:   endDte,
	}

	startDte, endDte, _ = utils.DatesFromDays(dateDayStart, dateDayEnd)
	suite.dateDays = &model.RequestDates{
		DateFrom: startDte,
		DateTo:   endDte,
	}
}

// ===================== Un-exported Functions ============================================ //

// TestGetBankCardsError method
// use this with modifications to the method to test error
func (suite *MongoSuite) TestGetBankCardsError() {

	dateMonthStr := "2200-12"

	startDte, endDte, _ := utils.DatesFromMonth(dateMonthStr)
	dts := &model.RequestDates{
		DateFrom: startDte,
		DateTo:   endDte,
	}

	_, err := suite.db.GetBankCards(dts)
	suite.Error(err)

	var mongoError *pkgerrors.MongoError
	if ok := errors.As(err, &mongoError); ok {
		// fmt.Printf("As mongoError: %v\n", mongoError)
		// fmt.Printf("mongoError.Err: %v\n", mongoError.Err)
		// handle gracefully
		// fmt.Printf("mongoError.Msg: %v\n", mongoError.Msg)
		suite.Equal(mongoError.Msg, noRecordsMsg)
	}
	if errors.Is(err, mongoError) {
	} else if err != nil {
		fmt.Printf("error: %+v\n", err)
		suite.Equal(mongoError.Msg, noRecordsMsg)
	}
}

// TestSetConfig
func (suite *MongoSuite2) TestSetConfig() {
	err := suite.db.setConfig()
	suite.NoError(err)
	suite.Equal(cfgHST, suite.db.cfg.HST)
}

// TestFetchBankCards method
func (suite *MongoSuite2) TestFetchBankCards() {
	records, err := suite.db.fetchBankCards(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	suite.NoError(err)
	suite.True(len(records) > 10)
}

// TestFetchEmployeeOS method
func (suite *MongoSuite2) TestFetchEmployeeOS() {
	records, err := suite.db.fetchEmployeeOS(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	suite.NoError(err)
	suite.True(len(records) > 10)
}

// TestFetchMonthlyGalesLoyalty
func (suite *MongoSuite2) TestFetchMonthlyGalesLoyalty() {
	_, err := suite.db.fetchGalesLoyalty(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	suite.NoError(err)
}

// TestFetchGalesLoyalty
func (suite *MongoSuite2) TestFetchGalesLoyalty() {
	docs, err := suite.db.fetchGalesLoyalty(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	suite.NoError(err)
	suite.True(len(docs) > 0)
}

// TestFetchFuelSales method
func (suite *MongoSuite2) TestFetchFuelSales() {
	docs, err := suite.db.fetchFuelSales(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	suite.NoError(err)
	suite.True(len(docs) > 0)
}

// TestFetchMonthlyNonFuel method
func (suite *MongoSuite2) TestFetchMonthlyNonFuel() {
	nf, err := suite.db.fetchMonthlyNonFuel(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	suite.NoError(err)
	suite.True(len(nf) > 0)
}

// TestFetchMonthlySales method
func (suite *MongoSuite2) TestFetchMonthlySales() {
	sales, err := suite.db.fetchMonthlySales(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	suite.NoError(err)
	suite.True(len(sales) > 0)
}

// TestFetchNonFuelSales method
func (suite *MongoSuite2) TestFetchNonFuelSales() {
	docs, err := suite.db.fetchNonFuelSales(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	suite.NoError(err)
	suite.True(len(docs) > 0)
}

// TestFetchPayPeriodSales method
func (suite *MongoSuite2) TestFetchPayPeriodSales() {
	records, err := suite.db.fetchPayPeriodSales(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	suite.NoError(err)
	suite.True(len(records) > 10)
}

// TestFetchProductNumbers method
func (suite *MongoSuite2) TestFetchProductNumbers() {
	records, err := suite.db.fetchProductNumbers(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	suite.NoError(err)
	suite.True(len(records) > 10)
}

// TestFetchCarWash method
func (suite *MongoSuite2) TestFetchCarWash() {
	records, err := suite.db.fetchCarWash(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	suite.NoError(err)
	suite.True(len(records) > 10)
}

// TestFetchNonFuelCommission method
func (suite *MongoSuite2) TestFetchNonFuelCommission() {
	records, _ := suite.db.fetchPayPeriodSales(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	// recordNum := records[0].RecordNum
	// stationID := records[0].StationID
	for _, r := range records {
		_, err := suite.db.fetchNonFuelCommission(r.RecordNum, r.StationID)
		suite.NoError(err)
	}
}

// TestFetchPropaneProducts method
func (suite *MongoSuite2) TestFetchPropaneProducts() {
	docs, err := suite.db.fetchPropaneProducts(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	suite.NoError(err)
	suite.True(len(docs) > 0)
}

func TestIntegrationSuite2(t *testing.T) {
	suite.Run(t, new(MongoSuite2))
}
