package handlers

import "github.com/aws/aws-lambda-go/events"

var headers map[string]string = map[string]string{
	"Content-Type":                 "application/json",
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Methods": "*",
	"Access-Control-Allow-Headers": "*",
}

type Handler interface {
	Response(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	process()
}

type response struct {
	Body       string
	Headers    map[string]string
	StatusCode int
}

type responseBody struct {
	Code    string `json:"code"`
	Data    string `json:"data"`
	Message string `json:"message"`
}
