package lamba_test_test

import (
	"URLShortener/tests"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/net/context"
	"log"
	"os"
	"testing"
)

// First things first we need to set up the active environment to make sense with our tests.

// Need to place a s3_bucket_obj into the s3 bucket, store its object url into the dynamodb table then
// After this is complete we can start: getPreSignedUrl(get(hash))

// To run these tests and not incur and real aws usage, deploy your terraform and after make sure
// you change you're aws credential to a dummy or test profile

// just run the awslocal commands to create these in our localstacks...
// but then we wouldn't solve the issue which is we can't access our local dynamodb table from our codebase

const BUCKET_NAME = "url-s3-bucket"

func TestingConfig() aws.Config {
	awsEndpoint := "http://localhost:4566"
	awsRegion := "us-east-2"

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
						PartitionID:       "aws",
						URL:               awsEndpoint,
						HostnameImmutable: true, // without this you won't be able to find the
					},
					nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "dummy cfg for localstack",
			},
		}),
	)
	if err != nil {
		log.Fatalf("Failed to get a configuration for aws: %v", err)
	}

	return cfg
}

func setup() {

	// You can use noSql and connect to local stacks to see these changes easier
	// cfg will need to be a parameter for these functions because we'll need to use the real creds eventually

	// PLACING S3OBJ INTO S3BUCKET //
	cfg := TestingConfig()
	client := s3.NewFromConfig(cfg)

	file, err := os.OpenFile("url.txt", os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		err := os.Remove("url.txt")
		if err != nil {
			fmt.Printf("Failed to delete file: %v\n", err)
			return
		}
	}

	_, err = file.Write([]byte("www.google.com\n"))
	if err != nil {
		panic(err)
	}

	bucket := BUCKET_NAME
	objInput := s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(file.Name()),
		Body:   file,
	}
	_, err = client.PutObject(context.TODO(), &objInput)
	if err != nil {
		log.Fatal("Failed to put Object: ", err)
	}

	// GETTING S3OBJECT URL AND STORING IT INTO DYNAMODB //
	region := "us-east-2"
	s3ObjUrl := "https://%s.amazonaws.com/%s/%s"
	s3ObjUrl = fmt.Sprintf(s3ObjUrl, region, bucket, file.Name())

	tests.Put("S3URLS", "hash", s3ObjUrl)
}

func TestGetPresignedURL(t *testing.T) {
	url := tests.Get("S3URLS", "hash")
	if url == "" {
		setup()
		url = tests.Get("S3URLS", "hash")
	}

}
