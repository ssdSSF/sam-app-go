package main

import (
	"bytes"
	"crud/model"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	buf := bytes.NewBufferString("")
	encoder := json.NewEncoder(buf)
	encoder.SetIndent("", "  ")
	encoder.Encode(model.Student{
		Id:        2345,
		FirstName: "Charels",
		LastName:  "Shi",
		Email:     "charles@srvusd.net",
	})

	return events.APIGatewayProxyResponse{
		Body: buf.String(),
		Headers: map[string]string{
			"content-type":                 "application-json",
			"customer-header":              "very-custom",
			"Access-Control-Allow-Headers": "*",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "*",
		},
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
