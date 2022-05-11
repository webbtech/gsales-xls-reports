package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"github.com/webbtech/gsales-xls-reports/config"
	"github.com/webbtech/gsales-xls-reports/handlers"
)

var (
	cfg *config.Config
	// db     model.DbHandler
	// client *mongo.Client
)

// init isn't called for each invocation, so we take advantage and only setup cfg and db for (I'm assuming) cold starts
func init() {
	log.Info("calling config.Config.init in main")
	cfg = &config.Config{}
	err := cfg.Init()
	if err != nil {
		log.Fatal(err)
		return
	}

	// db, err = mongodb.NewDb(cfg.GetMongoConnectString(), cfg.GetDbName())
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var h handlers.Handler

	switch request.Path {
	// case "/pdf":
	// 	h = &handlers.Pdf{Cfg: cfg, Db: db}
	default:
		h = &handlers.Ping{}
	}

	return h.Response(request)
}

func main() {
	lambda.Start(handler)
}
