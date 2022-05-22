package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

// application/zip
func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf("Processing Lambda req %s\n", req.RequestContext.RequestID)

	log.Println(req)
	res := events.APIGatewayProxyResponse{
		//302: found, 301: moved permanently, 300: multiple location available.
		StatusCode: 302,
		Headers: map[string]string{
			"Location": "https://www.google.com",
		},
	}
	log.Println("-------")
	log.Println(res)
	return res, nil
}

func main() {
	lambda.Start(handler)
}
