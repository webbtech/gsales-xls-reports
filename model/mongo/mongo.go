package mongo

import (
	"context"
	"log"

	"github.com/pulpfree/gsales-xls-reports/model"
)

// MDB struct
type MDB struct {
	client *mongo.Client
	dbName string
	db     *mongo.Database
}

// NewDB sets up new MDB struct
func NewDB(connection string, dbNm string) (model.DBHandler, error) {

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

	// defer suite.db.Close()

	return &MDB{
		client: client,
		dbName: dbNm,
		db:     client.Database(dbNm),
	}, err
}
