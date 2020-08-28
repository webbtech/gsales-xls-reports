package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/pulpfree/gsales-xls-reports/config"
	"github.com/pulpfree/gsales-xls-reports/model"
	"github.com/pulpfree/gsales-xls-reports/util"
	"github.com/pulpfree/pkgerrors"
)

const (
	cfgHST = int32(13)
	// dateMonth     = "2019-11"
	dateMonth     = "2020-04"
	dateDayStart  = "2020-02-01"
	dateDayEnd    = "2019-02-28"
	defaultsFP    = "../../config/defaults.yml"
	employeeIDStr = "5733c671982d828347021ed7"
	employeeName  = "Grimstead, Kevin"
	timeForm      = "2006-01-02"
)

// IntegSuite struct
type IntegSuite struct {
	cfg *config.Config
	// dates *model.RequestDates
	dateMonth *model.RequestDates
	dateDays  *model.RequestDates
	db        *MDB
	suite.Suite
}

// SetupTest method
func (suite *IntegSuite) SetupTest() {
	// setup config
	os.Setenv("Stage", "test")
	suite.cfg = &config.Config{DefaultsFilePath: defaultsFP}
	err := suite.cfg.Load()
	suite.NoError(err)

	// Set client options
	clientOptions := options.Client().ApplyURI(suite.cfg.GetMongoConnectURL())

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	suite.NoError(err)

	suite.db = &MDB{
		client: client,
		dbName: suite.cfg.DBName,
		db:     client.Database(suite.cfg.DBName),
	}

	inputMonth := &model.RequestInput{
		Date: dateMonth,
	}
	suite.dateMonth, _ = util.CreateDates(inputMonth)
	inputDayRange := &model.RequestInput{
		DateFrom: dateDayStart,
		DateTo:   dateDayEnd,
	}
	suite.dateDays, _ = util.CreateDates(inputDayRange)
}

// ===================== Exported Functions =============================================== //

// TestNewDB method
func (suite *IntegSuite) TestNewDB() {
	_, err := NewDB(suite.cfg.GetMongoConnectURL(), suite.cfg.DBName)
	suite.NoError(err)
}

// TestGetConfig method
func (suite *IntegSuite) TestGetConfig() {
	cfg, err := suite.db.GetConfig()
	suite.NoError(err)
	suite.Equal(cfgHST, cfg.HST)
}

// TestGetStationMap method
func (suite *IntegSuite) TestGetStationMap() {
	sm, err := suite.db.GetStationMap()
	suite.NoError(err)
	suite.True(len(sm) > 0)
}

// TestGetEmployeeOS
func (suite *IntegSuite) TestGetEmployeeOS() {
	sales, err := suite.db.GetEmployeeOS(suite.dateMonth)
	suite.NoError(err)
	suite.True(len(sales) > 0)
}

// TestGetCarWash
func (suite *IntegSuite) TestGetCarWash() {
	sales, err := suite.db.GetCarWash(suite.dateMonth)
	suite.NoError(err)
	suite.True(len(sales) > 0)
}

// TestGetMonthlySales
func (suite *IntegSuite) TestGetMonthlySales() {
	sales, err := suite.db.GetMonthlySales(suite.dateMonth)
	suite.NoError(err)
	suite.True(len(sales) > 0)
}

// TestGetMonthlyProducts
func (suite *IntegSuite) TestGetMonthlyProducts() {
	sales, err := suite.db.GetMonthlyProducts(suite.dateMonth)
	suite.NoError(err)
	suite.True(len(sales) > 0)
}

// TestGetNonFuelSales
func (suite *IntegSuite) TestGetNonFuelSales() {
	sales, err := suite.db.GetNonFuelSales(suite.dateMonth)
	suite.NoError(err)
	suite.True(len(sales) > 0)
}

// TestGetPayPeriodSales
func (suite *IntegSuite) TestGetPayPeriodSales() {
	sales, err := suite.db.GetPayPeriodSales(suite.dateMonth)
	suite.NoError(err)
	suite.True(len(sales) > 0)
}

// TestGetPropaneSales
func (suite *IntegSuite) TestGetPropaneSales() {
	sales, err := suite.db.GetPropaneSales(suite.dateMonth)
	suite.NoError(err)
	suite.True(len(sales) > 0)
}

// ===================== Un-exported Functions ============================================ //

// TestsetConfig
func (suite *IntegSuite) TestsetConfig() {
	err := suite.db.setConfig()
	suite.NoError(err)
	suite.Equal(cfgHST, suite.db.cfg.HST)
}

// TestfetchBankCards method
func (suite *IntegSuite) TestfetchBankCards() {
	records, err := suite.db.fetchBankCards(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	suite.NoError(err)
	suite.True(len(records) > 10)
	// fmt.Printf("records %+v\n", records[0])
}

// TestfetchEmployeeOS method
func (suite *IntegSuite) TestfetchEmployeeOS() {
	records, err := suite.db.fetchEmployeeOS(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	suite.NoError(err)
	suite.True(len(records) > 10)
}

// TestfetchMonthlyGalesLoyalty
func (suite *IntegSuite) TestfetchMonthlyGalesLoyalty() {
	_, err := suite.db.fetchGalesLoyalty(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	suite.NoError(err)
	// suite.True(len(nf) > 0)
}

// TestfetchGalesLoyalty
func (suite *IntegSuite) TestfetchGalesLoyalty() {
	fmt.Printf("suite.dateDays %+v\n", suite.dateDays)
	docs, err := suite.db.fetchGalesLoyalty(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	suite.NoError(err)
	fmt.Printf("docs %+v\n", docs[0])
	// _, err := suite.db.fetchGalesLoyalty(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	// suite.NoError(err)
	// suite.True(len(nf) > 0)
}

// TestGetBankCardsError method
// use this with modifications to the method to test error
func (suite *IntegSuite) TestGetBankCardsError() {

	inputMonth := &model.RequestInput{
		Date: "2200-12",
	}
	dateMonth, _ := util.CreateDates(inputMonth)

	_, err := suite.db.GetBankCards(dateMonth)
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

// TestfetchMonthlyNonFuel method
func (suite *IntegSuite) TestfetchMonthlyNonFuel() {
	nf, err := suite.db.fetchMonthlyNonFuel(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	suite.NoError(err)
	suite.True(len(nf) > 0)
	// fmt.Printf("nf %+v\n", nf[0])
}

// TestfetchMonthlySales method
func (suite *IntegSuite) TestfetchMonthlySales() {
	sales, err := suite.db.fetchMonthlySales(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	suite.NoError(err)
	suite.True(len(sales) > 0)
}

// TestfetchNonFuelSales method
func (suite *IntegSuite) TestfetchNonFuelSales() {
	docs, err := suite.db.fetchNonFuelSales(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	suite.NoError(err)
	suite.True(len(docs) > 0)
}

// TestfetchPayPeriodSales method
func (suite *IntegSuite) TestfetchPayPeriodSales() {
	records, err := suite.db.fetchPayPeriodSales(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	suite.NoError(err)
	suite.True(len(records) > 10)
}

// TestfetchProductNumbers method
func (suite *IntegSuite) TestfetchProductNumbers() {
	records, err := suite.db.fetchProductNumbers(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	suite.NoError(err)
	suite.True(len(records) > 10)
}

// TestfetchPayPeriodSales method
func (suite *IntegSuite) TestfetchCarWash() {
	records, err := suite.db.fetchCarWash(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	suite.NoError(err)
	suite.True(len(records) > 10)
}

// TestfetchNonFuelCommission method
func (suite *IntegSuite) TestfetchNonFuelCommission() {
	records, _ := suite.db.fetchPayPeriodSales(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	// recordNum := records[0].RecordNum
	// stationID := records[0].StationID
	for _, r := range records {
		_, err := suite.db.fetchNonFuelCommission(r.RecordNum, r.StationID)
		suite.NoError(err)
	}
}

// TestfetchPropaneProducts method
func (suite *IntegSuite) TestfetchPropaneProducts() {
	// fmt.Printf("suite.dateDays %+v\n", suite.dateDays)
	docs, err := suite.db.fetchPropaneProducts(suite.dateMonth.DateFrom, suite.dateMonth.DateTo)
	suite.NoError(err)
	suite.True(len(docs) > 0)
}

// ===================== Utility Functions ===================================================== //

// TestsetEmployee method
func (suite *IntegSuite) TestsetEmployee() {
	employeeID, _ := primitive.ObjectIDFromHex(employeeIDStr)
	employee, err := suite.db.setEmployee(employeeID)
	suite.NoError(err)
	suite.Equal(employee, employeeName)
}

// TestGetEmployee method
func (suite *IntegSuite) TestGetEmployee() {
	employeeID, _ := primitive.ObjectIDFromHex(employeeIDStr)
	employee, err := suite.db.GetEmployee(employeeID)
	suite.NoError(err)
	suite.Equal(employee, employeeName)
}

// TestsetStationMap method
func (suite *IntegSuite) TestsetStationMap() {
	err := suite.db.setStationMap()
	suite.NoError(err)
}

// func (suite *IntegSuite) TestInitPayPeriod() {
// pp, err :=
// }

// TestIntegrationSuite function
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}
