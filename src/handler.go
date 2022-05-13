package main

import (
	_ "common"
	"context"
	_ "fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/aws/aws-sdk-go-v2/aws"
	_ "github.com/aws/aws-sdk-go-v2/config"
	_ "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	_ "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	_ "os"
	_ "tests"
	_ "tests/lambda_test"
)

// TODO - Fix Presign calcuation error when using presigned URL in local stacks
// TODO - Test handler function using SAM

// basic skeleton for the redirect

// Handler testing to see if this lambda implementation really works

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	//cfg, _ := config.LoadDefaultConfig(context.TODO())
	//dynamodbClient := dynamodb.NewFromConfig(cfg)

	// logging
	//fmt.Printf("event.HTTPMethod %v\n", request.HTTPMethod)
	//fmt.Printf("event.Body %v\n", request.Body)
	//fmt.Printf("event.QueryStringParameters %v\n", request.QueryStringParameters)
	//fmt.Printf("event %v\n", request)

	//UUID := request.QueryStringParameters["UUID"]

	//var bucket string
	//var filename string

	//if request.HTTPMethod == "GET" {
	//	bucket, _, filename = tests.Get(dynamodbClient, common.TableName, UUID)
	//}
	//psUrl := lambda_test.CreatePresignedURL(cfg, bucket, filename) // PreSigned URL
	//err := lambda_test.DownloadS3Object(psUrl, filename)
	//if err != nil {
	//	fmt.Printf("Failed to download s3 object")
	//	return events.APIGatewayProxyResponse{}, err
	//}
	//
	//object, err := os.ReadFile(filename)
	//if err != nil{
	//	fmt.Printf("Failed to Read file, %v", err)
	//}
	//
	//url = string(object)

	url := "https://iuscsg.org"
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
