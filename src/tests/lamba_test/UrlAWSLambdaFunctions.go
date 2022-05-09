package lamba_test

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3PresignGetObjectAPI defines the interface for the PresignGetObject function.
// We use this interface to test the function using a mocked service.

type S3PresignGetObjectAPI interface {
	PresignGetObject(
		ctx context.Context,
		params *s3.GetObjectInput,
		optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

func GetPresignedURL(c context.Context, api S3PresignGetObjectAPI, input *s3.GetObjectInput) (*v4.PresignedHTTPRequest, error) {
	return api.PresignGetObject(c, input)
}

// given s3 url

func DownloadS3Object(config aws.Config, key string) {

}

// Remember we're accessing the s3 bucket via the url in the dyanmodb not directly from s3

func GetPresignedURL(config aws.Config, bucket string, key string) string {

	client := s3.NewFromConfig(config)
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	psClient := s3.NewPresignClient(client)

	response, err := GetPresignedURL(context.TODO(), psClient, input)
	if err != nil {
		fmt.Println("Error retrieving pre-signed object:")
		panic(err)
	}

	fmt.Println("The URL: ")
	fmt.Println(response.URL)

	return response.URL
}
