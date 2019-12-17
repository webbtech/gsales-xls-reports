package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pulpfree/gsales-xls-reports/config"
	"github.com/pulpfree/gsales-xls-reports/model"
	"github.com/pulpfree/gsales-xls-reports/report"
	"github.com/pulpfree/gsales-xls-reports/validate"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response data format
type Response struct {
	Code      int         `json:"code"`      // HTTP status code
	Data      interface{} `json:"data"`      // Data payload
	Message   string      `json:"message"`   // Error or status message
	Status    string      `json:"status"`    // Status code (error|fail|success)
	Timestamp int64       `json:"timestamp"` // Machine-readable UTC timestamp in nanoseconds since EPOCH
}

// SignedURL struct
type SignedURL struct {
	URL string `json:"url"`
}

var cfg *config.Config

func init() {
	cfg = &config.Config{}
	err := cfg.Load()
	if err != nil {
		log.Fatal(err)
	}
}

// HandleRequest function
func HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	hdrs := make(map[string]string)
	hdrs["Content-Type"] = "application/json"
	hdrs["Access-Control-Allow-Origin"] = "*"
	t := time.Now()

	// If this is a ping test, intercept and return
	if req.HTTPMethod == "GET" {
		log.Info("Ping test in handleRequest")
		return gatewayResponse(Response{
			Code:      200,
			Data:      "pong",
			Status:    "success",
			Timestamp: t.Unix(),
		}, hdrs), nil
	}

	var r *model.RequestInput
	json.Unmarshal([]byte(req.Body), &r)

	// validate input
	reportRequest, err := validate.SetRequest(r)
	if err != nil {
		return gatewayResponse(Response{
			Code:      500,
			Message:   fmt.Sprintf("error: %s", err.Error()),
			Status:    "error",
			Timestamp: t.Unix(),
		}, hdrs), nil
	}
	rpt := report.New(reportRequest, cfg)
	url, err := rpt.CreateSignedURL()
	if err != nil {
		return gatewayResponse(Response{
			Code:      500,
			Message:   fmt.Sprintf("error: %s", err.Error()),
			Status:    "error",
			Timestamp: t.Unix(),
		}, hdrs), nil
	}
	log.Infof("signed url created %s", url)

	return gatewayResponse(Response{
		Code:      201,
		Data:      SignedURL{URL: url},
		Status:    "success",
		Timestamp: t.Unix(),
	}, hdrs), nil
}

func main() {
	lambda.Start(HandleRequest)
}

func gatewayResponse(resp Response, hdrs map[string]string) events.APIGatewayProxyResponse {

	body, _ := json.Marshal(&resp)
	if resp.Status == "error" {
		log.Errorf("Error: status: %s, code: %d, message: %s", resp.Status, resp.Code, resp.Message)
	}

	return events.APIGatewayProxyResponse{Body: string(body), Headers: hdrs, StatusCode: resp.Code}
}
