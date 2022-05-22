package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"test/m/v2/aws_lambda"
)

// basic skeleton for the redirect

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	cfg, _ := config.LoadDefaultConfig(context.TODO())

	fmt.Printf("event.HTTPMethod %v\n", request.HTTPMethod)
	fmt.Printf("event.QueryStringParameters %v\n", request.QueryStringParameters)
	fmt.Printf("event %v\n", request)

	psUrl := aws_lambda.CreatePresignedURL(cfg, "assignment_doc_bucket", "How_to_Install_Jetbrains_Toolbox_and_IDEs") // PreSigned URL
	url, err := aws_lambda.GetURLObject(psUrl, "How_to_Install_Jetbrains_Toolbox_and_IDEs")
	if err != nil {
		fmt.Printf("Failed to get url s3 object")
		return events.APIGatewayProxyResponse{}, err
	}

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
