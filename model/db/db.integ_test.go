package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/pulpfree/gsales-xls-reports/config"
	"github.com/pulpfree/gsales-xls-reports/model"
)

const (
	cfgHST  = int32(13)
	dateStr = "2019-08-01"
	// dateStr       = "2019-09-01"
	defaultsFP    = "../../config/defaults.yml"
	employeeIDStr = "5733c671982d828347021ed7"
	employeeName  = "Grimstead, Kevin"
	timeForm      = "2006-01-02"
)

/* var (
	endDate   time.Time
	startDate time.Time
) */

// IntegSuite struct
type IntegSuite struct {
	cfg   *config.Config
	dates *model.RequestDates
	db    *MDB
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

	// Set start and end dates for monthly reports
	t, err := time.Parse(timeForm, dateStr)
	if err != nil {
		panic(err)
	}
	currentYear, currentMonth, _ := t.Date()
	currentLocation := t.Location()
	dte := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	suite.dates = &model.RequestDates{
		DateFrom: dte,
		DateTo:   dte.AddDate(0, 1, -1),
	}
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

// TestGetMonthlySales
func (suite *IntegSuite) TestGetMonthlySales() {
	sales, err := suite.db.GetMonthlySales(suite.dates)
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
	records, err := suite.db.fetchBankCards(suite.dates.DateFrom, suite.dates.DateTo)
	suite.NoError(err)
	suite.True(len(records) > 10)
	// fmt.Printf("records %+v\n", records[0])
}

// TestGetBankCardsError method
// use this with modifications to the method to test error
func (suite *IntegSuite) TestGetBankCardsError() {

	_, err := suite.db.GetBankCards(suite.dates)
	suite.Error(err)

	var mongoError *MongoError
	if ok := errors.As(err, &mongoError); ok {
		fmt.Printf("As mongoError: %v\n", mongoError)
		fmt.Printf("mongoError.err: %v\n", mongoError.err)
		// handle gracefully
		fmt.Printf("mongoError.more: %v\n", mongoError.more)
		return
	}
	if errors.Is(err, mongoError) {
		fmt.Printf("Is mongoError %+v\n", mongoError.more)
	} else if err != nil {
		fmt.Printf("error: %+v\n", err)
	}

}

// TestfetchMonthlyNonFuel method
func (suite *IntegSuite) TestfetchMonthlyNonFuel() {
	nf, err := suite.db.fetchMonthlyNonFuel(suite.dates.DateFrom, suite.dates.DateTo)
	suite.NoError(err)
	suite.True(len(nf) > 0)
	// fmt.Printf("nf %+v\n", nf[0])
}

// TestfetchPayPeriodSales method
func (suite *IntegSuite) TestfetchPayPeriodSales() {
	records, err := suite.db.fetchPayPeriodSales(suite.dates.DateFrom, suite.dates.DateTo)
	suite.NoError(err)
	suite.True(len(records) > 10)
	// fmt.Printf("records[0] %+v\n", records[0])
}

// TestfetchPayPeriodSales method
func (suite *IntegSuite) TestfetchCarWash() {
	records, err := suite.db.fetchCarWash(suite.dates.DateFrom, suite.dates.DateTo)
	suite.NoError(err)
	suite.True(len(records) > 10)
}

// TestfetchNonFuelCommission method
func (suite *IntegSuite) TestfetchNonFuelCommission() {
	records, _ := suite.db.fetchPayPeriodSales(suite.dates.DateFrom, suite.dates.DateTo)
	// recordNum := records[0].RecordNum
	// stationID := records[0].StationID
	for _, r := range records {
		_, err := suite.db.fetchNonFuelCommission(r.RecordNum, r.StationID)
		suite.NoError(err)
	}
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
