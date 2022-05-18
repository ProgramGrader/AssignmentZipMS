package main

import (
	"aws_lambda"
	"common"
	"context"
	"dynamoDAO"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// TODO - Test handler function using SAM
// TODO - Refactor code, reduce redundancy, use depinj

// basic skeleton for the redirect

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	cfg, _ := config.LoadDefaultConfig(context.TODO())
	dynamodbClient := dynamodb.NewFromConfig(cfg)

	fmt.Printf("event.HTTPMethod %v\n", request.HTTPMethod)
	fmt.Printf("event.QueryStringParameters %v\n", request.QueryStringParameters)
	fmt.Printf("event %v\n", request)

	UUID := request.QueryStringParameters["UUID"]
	if UUID == "" {
		// return status code
	}

	var bucket string
	var filename string

	if request.HTTPMethod == "GET" {
		bucket, _, filename = dynamoDAO.Get(dynamodbClient, common.TableName, UUID)
	}
	psUrl := aws_lambda.CreatePresignedURL(cfg, bucket, filename) // PreSigned URL
	url, err := aws_lambda.GetURLObject(psUrl, filename)
	if err != nil {
		fmt.Printf("Failed to download s3 object")
		return events.APIGatewayProxyResponse{}, err
	}

	//url := "https://iuscsg.org"
	return events.APIGatewayProxyResponse{
		//307: temporary redirect, 302: found, 301: moved permanently, 300: multiple location available.
		//Body:       fmt.Sprintf("{\"message\":\"Error occurred unmarshaling request: %v.\"}", url),
		StatusCode: 307,
		Headers: map[string]string{
			"Location": url,
		},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
