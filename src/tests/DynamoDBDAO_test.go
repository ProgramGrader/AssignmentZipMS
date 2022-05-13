package tests

import (
	"common"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"testing"
)

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

var clientConfig, _ = CreateLocalClient()

type Item struct {
	UUID     string
	bucket   string
	region   string
	filename string

	expectedBucket   string
	expectedRegion   string
	expectedFilename string
}

func TestPutAndGet(t *testing.T) {

	putAnGetItems := []Item{
		{"uuid1", "a-bucket", "us-east-2", "file.txt", "a-bucket", "us-east-2", "file.txt"},
		{"uuid2", "the-bucket", "us-east-1", "file.txt", "the-bucket", "us-east-1", "file.txt"},
		{"uuid2", "the-bucket", "us-west-1", "file.txt", "the-bucket", "us-west-1", "file.txt"},
	}

	for _, test := range putAnGetItems {

		Put(clientConfig, common.TableName, test.UUID, test.bucket, test.region, test.filename)
		storedBucket, storedRegion, storedFilename := Get(clientConfig, common.TableName, test.UUID)
		if storedBucket != test.expectedBucket {
			t.Fatalf("TestPut(), Failed. Expected value was not found. Got %s expected %s", storedBucket, test.expectedBucket)
		} else if storedRegion != test.expectedRegion {
			t.Fatalf("TestPut(), Failed. Expected value was not found. Got %s expected %s", storedRegion, test.expectedRegion)
		} else if storedFilename != test.expectedFilename {
			t.Fatalf("TestPut(), Failed. Expected value was not found. Got %s expected %s", storedFilename, test.expectedFilename)
		}
	}

	DeleteAll(clientConfig, common.TableName)
}

func TestDelete(t *testing.T) {

	deleteItems := []Item{
		{"uuid1", "a_bucket", "us-east-1", "file.txt", "", "", ""},
		{"uuid2", "update_url", "us-east-1", "file.txt", "", "", ""},
		{"uuid2", "update_url", "us-east-2", "file.txt", "", "", ""},
	}

	for _, test := range deleteItems {
		Put(clientConfig, common.TableName, test.UUID, test.bucket, test.region, test.filename)
		deleteErr := Delete(clientConfig, common.TableName, test.UUID)
		if deleteErr != nil {
			t.Fatal("TestDelete(), Failed to delete. Expected error to be nil.")
		}
	}
}
