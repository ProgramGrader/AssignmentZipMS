package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// basic skeleton for the redirect

// Handler testing to see if this lambda implementation really works

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// logging
	fmt.Printf("event.HTTPMethod %v\n", request.HTTPMethod)
	fmt.Printf("event.Body %v\n", request.Body)
	fmt.Printf("event.QueryStringParameters %v\n", request.QueryStringParameters)
	fmt.Printf("event %v\n", request)

	url := "www.google.com"

	// this would be how we get the hash from url
	//if request.HTTPMethod == "GET" {
	//	hash = request.QueryStringParameters["hash"]
	//}
	//

	//  getS3URL(){}
	// 	createTempURL(){}

	// the rest of the error codes are to be handled in terraform, specifically through aws_api_integration_response
	return events.APIGatewayProxyResponse{
		//302: found, 301: moved permanently, 300: multiple location available.
		StatusCode: 301,
		Headers: map[string]string{
			"Location": url,
		},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
