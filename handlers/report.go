package handlers

import (
	"encoding/json"
	"time"

	pres "github.com/pulpfree/lambda-go-proxy-response"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/webbtech/gsales-xls-reports/config"
	"github.com/webbtech/gsales-xls-reports/model"
	"github.com/webbtech/gsales-xls-reports/report"
	"github.com/webbtech/gsales-xls-reports/validate"
)

// SignedURL struct
type SignedURL struct {
	URL string `json:"url"`
}

var cfg *config.Config

func init() {
	cfg = &config.Config{}
	err := cfg.Init()
	if err != nil {
		log.Fatal(err)
	}
}

// HandleRequest function
func HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// see following link for tips
	// https://stackoverflow.com/questions/53298478/has-been-blocked-by-cors-policy-response-to-preflight-request-doesn-t-pass-acce

	hdrs := make(map[string]string)
	hdrs["Content-Type"] = "application/json"
	hdrs["Access-Control-Allow-Origin"] = "*"
	hdrs["Access-Control-Allow-Methods"] = "GET,OPTIONS,POST,PUT"
	hdrs["Access-Control-Allow-Headers"] = "Authorization,Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token" // just added

	if req.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{Body: string("null"), Headers: hdrs, StatusCode: 200}, nil
	}

	t := time.Now()

	// If this is a ping test, intercept and return
	if req.HTTPMethod == "GET" {
		log.Info("Ping test in handleRequest")
		return pres.ProxyRes(pres.Response{
			Code:      200,
			Data:      "pong",
			Status:    "success",
			Timestamp: t.Unix(),
		}, hdrs, nil), nil
	}

	var r *model.RequestInput
	json.Unmarshal([]byte(req.Body), &r)

	// validate input
	reportRequest, err := validate.SetRequest(r)
	if err != nil {
		return pres.ProxyRes(pres.Response{
			Timestamp: t.Unix(),
		}, hdrs, err), nil
	}

	rpt, err := report.New(reportRequest, cfg)
	if err != nil {
		return pres.ProxyRes(pres.Response{
			Timestamp: t.Unix(),
		}, hdrs, err), nil
	}

	url, err := rpt.CreateSignedURL()
	if err != nil {
		return pres.ProxyRes(pres.Response{
			Timestamp: t.Unix(),
		}, hdrs, err), nil
	}

	urlStr := url[0:100]
	log.Infof("signed url created %s", urlStr)

	return pres.ProxyRes(pres.Response{
		Code:      201,
		Data:      SignedURL{URL: url},
		Status:    "success",
		Timestamp: t.Unix(),
	}, hdrs, nil), nil
}

func main() {
	lambda.Start(HandleRequest)
}
