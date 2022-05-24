package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"github.com/webbtech/gsales-xls-reports/config"
	"github.com/webbtech/gsales-xls-reports/handlers"
	"github.com/webbtech/gsales-xls-reports/model"
	"github.com/webbtech/gsales-xls-reports/mongodb"
)

var (
	cfg *config.Config
	db  model.DbHandler
	// client *mongo.Client
	err error
)

// init isn't called for each invocation, so we take advantage and only setup cfg and db for (I'm assuming) cold starts
func init() {

	// For local unit testing, we need the Stage env var set to test
	// When run as a lambda function this is already set to the appropriate environment
	_, exists := os.LookupEnv("Stage")
	if !exists {
		os.Setenv("Stage", "test")
	}

	log.Info("calling config.Config.Init in main")
	cfg = &config.Config{}
	err = cfg.Init()
	if err != nil {
		log.Fatal(err)
		return
	}

	db, err = mongodb.NewDB(cfg.GetMongoConnectURL(), cfg.DbName)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var h handlers.Handler

	switch request.Path {
	case "/report":
		h = &handlers.Report{Cfg: cfg, Db: db}
	default:
		h = &handlers.Ping{}
	}

	return h.Response(request)
}

func main() {
	lambda.Start(handler)
}
