package handlers

import (
	"github.com/aws/aws-lambda-go/events"
)

type Options struct {
	response events.APIGatewayProxyResponse
}

func (c *Options) Response(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	c.process()
	return c.response, nil
}

func (c *Options) process() {
	c.response = events.APIGatewayProxyResponse{
		Body:       string("null"),
		Headers:    headers,
		StatusCode: 200,
	}
}
