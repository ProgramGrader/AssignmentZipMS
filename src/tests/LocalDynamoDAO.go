package tests

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
)

// MAYBE: erase the Local from the name
// Turns out that to be able to use a endpoint in go it can't be https

// CreateLocalClient returns the config associated with local docker container
func CreateLocalClient() (*dynamodb.Client, error) {
	//hostName := fmt.Sprintf("http://%s:4566", os.Getenv("LOCALSTACK_HOSTNAME"))
	awsEndpoint := "http://localhost:4566"
	awsRegion := "us-east-2"

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: awsEndpoint}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "dummy cfg for localstack",
			},
		}),
	)
	if err != nil {
		panic(err)
	}

	return dynamodb.NewFromConfig(cfg), err
}

var localClientConfig, _ = CreateLocalClient()

// Get given hash returns value
func Get(tableName string, hash string) string {

	getItemInput := &dynamodb.GetItemInput{
		TableName:            aws.String(tableName),
		ConsistentRead:       aws.Bool(true),
		ProjectionExpression: aws.String("s3object"),

		Key: map[string]types.AttributeValue{
			"urlId": &types.AttributeValueMemberS{Value: hash},
		},
	}

	output, err := localClientConfig.GetItem(context.TODO(), getItemInput)
	if err != nil {
		log.Fatalf("Failed to get item, %v", err)
	}

	if output.Item == nil {
		log.Fatal("Item not found: ", hash)
	}

	var value string

	err = attributevalue.Unmarshal(output.Item["s3object"], &value)
	if err != nil {
		log.Fatalf("unmarshal failed, %v", err)
	}

	return value

}

// Put creates/update a new entry in the Dynamodb
func Put(tableName string, key string, value string) {

	var itemInput = dynamodb.PutItemInput{
		TableName: aws.String(tableName),

		Item: map[string]types.AttributeValue{
			"urlId":    &types.AttributeValueMemberS{Value: key},
			"s3object": &types.AttributeValueMemberS{Value: value},
		},
	}

	_, err := localClientConfig.PutItem(context.TODO(), &itemInput)
	if err != nil {
		log.Fatal("Error inserting value ", err)
	}
}

// Delete removes a item from the table given the key
func Delete(tableName string, key string) error {

	deleteInput := dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"urlId": &types.AttributeValueMemberS{Value: key},
		},
	}

	_, err := localClientConfig.DeleteItem(context.TODO(), &deleteInput)
	if err != nil {
		panic(err)
	}

	return err
}

// DeleteAll for testing purposes
func DeleteAll(tableName string) {
	scan := dynamodb.NewScanPaginator(localClientConfig, &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})

	for scan.HasMorePages() {
		out, err := scan.NextPage(context.TODO())
		if err != nil {
			print("Page error")
			panic(err)
		}

		for _, item := range out.Items {
			_, err = localClientConfig.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
				TableName: aws.String(tableName),
				Key: map[string]types.AttributeValue{
					"urlId": item["urlId"],
				},
			})
			if err != nil {
				print("Error Deleting Item")
				panic(err)
			}

		}
	}
}
