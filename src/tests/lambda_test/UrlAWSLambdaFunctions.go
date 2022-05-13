package lamba_test

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"os"
	"time"
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
	return api.PresignGetObject(c, input, s3.WithPresignExpires(300*time.Second))
}

// given s3 url

// Remember we're accessing the s3 bucket via the url in the DynamoDB not directly from s3

func CreatePresignedURL(config aws.Config, bucket string, key string) string {

	client := s3.NewFromConfig(config, func(options *s3.Options) {
		options.UsePathStyle = true
	})

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
	return response.URL
}

func DownloadS3Object(config aws.Config, key string) {
	//client := s3.NewFromConfig(config)
	//downloader := manager.NewDownloader(client)
	//
	//numBytes, err := downloader.Download(context.TODO(), downloadFile, &s3.GetObjectInput{
	//
	//})
}

func CreateTempURL(file os.File) {

}
