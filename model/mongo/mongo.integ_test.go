package mongo

import (
	"context"
	"os"
	"testing"

	"github.com/pulpfree/config"
	"github.com/stretchr/testify/suite"
)

// IntegSuite struct
type IntegSuite struct {
	suite.Suite
	cfg *config.Config
	db  *MDB
	// q   *model.Quote
}

const (
	defaultsFP = "../../config/defaults.yml"
	// quoteID    = "5cb7a1b5b413522c173b3cde"
)

// SetupTest method
func (suite *IntegSuite) SetupTest() {
	// setup config
	os.Setenv("Stage", "test")
	suite.cfg = &config.Config{DefaultsFilePath: defaultsFP}
	err := suite.cfg.Load()
	suite.NoError(err)
	// fmt.Printf("suite.cfg.GetMongoConnectURL() %s\n", suite.cfg.GetMongoConnectURL())

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
	// suite.q = &model.Quote{}
	// suite.q.Items = &model.Items{}
}

// TestNewDB method
func (suite *IntegSuite) TestNewDB() {
	_, err := NewDB(suite.cfg.GetMongoConnectURL(), suite.cfg.DBName)
	suite.NoError(err)
}

// TestIntegrationSuite function
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}
