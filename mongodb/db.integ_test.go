package mongodb

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/webbtech/gsales-xls-reports/config"
	"github.com/webbtech/gsales-xls-reports/handlers"
	"github.com/webbtech/gsales-xls-reports/model"
	"github.com/webbtech/gsales-xls-reports/utils"
)

const (
	cfgHST        = int32(13)
	dateMonth     = "2020-07"
	dateDayStart  = "2022-02-18"
	dateDayEnd    = "2022-02-19"
	defaultsFP    = "../../config/defaults.yml"
	employeeIDStr = "5733c671982d828347021ed7"
	employeeName  = "Grimstead, Kevin"
	timeForm      = "2006-01-02"
)

// MongoSuite struct
type MongoSuite struct {
	cfg       *config.Config
	dateMonth *model.RequestDates
	dateDays  *model.RequestDates
	db        *MDB
	suite.Suite
	rpt *handlers.Report
}

// SetupTest method
func (suite *MongoSuite) SetupTest() {
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

// ===================== Exported Functions =============================================== //

// TestNewDB method
func (suite *MongoSuite) TestNewDB() {
	_, err := NewDB(suite.cfg.GetMongoConnectURL(), suite.cfg.DbName)
	suite.NoError(err)
}

// TestGetConfig method
func (suite *MongoSuite) TestGetConfig() {
	cfg, err := suite.db.GetConfig()
	suite.NoError(err)
	suite.Equal(cfgHST, cfg.HST)
}

// TestGetStationMap method
func (suite *MongoSuite) TestGetStationMap() {
	sm, err := suite.db.GetStationMap()
	suite.NoError(err)
	suite.True(len(sm) > 0)
}

// TestGetEmployeeOS
func (suite *MongoSuite) TestGetEmployeeOS() {
	sales, err := suite.db.GetEmployeeOS(suite.dateMonth)
	suite.NoError(err)
	suite.True(len(sales) > 0)
}

// TestGetFuelSales
func (suite *MongoSuite) TestGetFuelSales() {
	sales, err := suite.db.GetFuelSales(suite.dateMonth)
	suite.NoError(err)
	suite.True(len(sales) > 0)
}

// TestGetCarWash
func (suite *MongoSuite) TestGetCarWash() {
	sales, err := suite.db.GetCarWash(suite.dateMonth)
	suite.NoError(err)
	suite.True(len(sales) > 0)
}

// TestGetMonthlySales
func (suite *MongoSuite) TestGetMonthlySales() {
	sales, err := suite.db.GetMonthlySales(suite.dateMonth)
	suite.NoError(err)
	suite.True(len(sales) > 0)
}

// TestGetMonthlyProducts
func (suite *MongoSuite) TestGetMonthlyProducts() {
	sales, err := suite.db.GetMonthlyProducts(suite.dateMonth)
	suite.NoError(err)
	suite.True(len(sales) > 0)
}

// TestGetNonFuelSales
func (suite *MongoSuite) TestGetNonFuelSales() {
	sales, err := suite.db.GetNonFuelSales(suite.dateMonth)
	suite.NoError(err)
	suite.True(len(sales) > 0)
}

// TestGetPayPeriodSales
func (suite *MongoSuite) TestGetPayPeriodSales() {
	sales, err := suite.db.GetPayPeriodSales(suite.dateMonth)
	suite.NoError(err)
	suite.True(len(sales) > 0)
}

// TestGetPropaneSales
func (suite *MongoSuite) TestGetPropaneSales() {
	sales, err := suite.db.GetPropaneSales(suite.dateMonth)
	suite.NoError(err)
	suite.True(len(sales) > 0)
}

// ===================== Utility Functions ===================================================== //

// TestSetEmployee method
func (suite *MongoSuite) TestSetEmployee() {
	employeeID, _ := primitive.ObjectIDFromHex(employeeIDStr)
	employee, err := suite.db.setEmployee(employeeID)
	suite.NoError(err)
	suite.Equal(employee, employeeName)
}

// TestGetEmployee method
func (suite *MongoSuite) TestGetEmployee() {
	employeeID, _ := primitive.ObjectIDFromHex(employeeIDStr)
	employee, err := suite.db.GetEmployee(employeeID)
	suite.NoError(err)
	suite.Equal(employee, employeeName)
}

// TestSetStationMap method
func (suite *MongoSuite) TestSetStationMap() {
	err := suite.db.setStationMap()
	suite.NoError(err)
}

// TestFetchStationNodes method
func (suite *MongoSuite) TestFetchStationNodes() {
	nodes, err := suite.db.fetchStationNodes()
	suite.NoError(err)
	suite.True(len(nodes) > 10)
}

// TestIntegrationSuite function
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(MongoSuite))
}
