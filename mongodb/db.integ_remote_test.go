package mongodb

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/webbtech/gsales-xls-reports/config"
	"github.com/webbtech/gsales-xls-reports/handlers"
	"github.com/webbtech/gsales-xls-reports/model"
	"github.com/webbtech/gsales-xls-reports/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRemoteSuite struct
type MongoRemoteSuite struct {
	cfg       *config.Config
	dateMonth *model.RequestDates
	dateDays  *model.RequestDates
	db        *MDB
	suite.Suite
	rpt *handlers.Report
}

// SetupTest method
func (suite *MongoRemoteSuite) SetupTest() {
	// setup config
	os.Setenv("Stage", "prod")
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

// TestNewDB2 method
func (suite *MongoSuite) TestNewDB2() {
	_, err := NewDB(suite.cfg.GetMongoConnectURL(), suite.cfg.DbName)
	suite.NoError(err)
}

// TestGetStationMap2 method
func (suite *MongoSuite) TestGetStationMap2() {
	sm, err := suite.db.GetStationMap()
	suite.NoError(err)
	suite.True(len(sm) > 0)
}

func TestMongoRemoteSuite(t *testing.T) {
	suite.Run(t, new(MongoSuite2))
}
