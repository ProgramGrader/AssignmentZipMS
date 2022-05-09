package lamba_test

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Unit Testing functions that will be used in aws lambda

// get s3url from bucket and verify that it exists.
// The dynamodb table has the path to the s3object that stores the file containing the url that we need to use
// to get temp url and redirect the user, to do this we need to create a pre-resigned url to bypass authentication

// create a temporary url --> s3url takes you to the obj which is a file with the url inside it, you need to convert
// the s3 object to a file then extract its contents. Then create the Url

// redirect user to said temp url

// Pre-signed Url needs to be created then stored into the local dynamodb because the url is temporarily available
// meaning if one is created manually testing will eventually not work because the URL will expire

// Creates and returns a Pre-signed URL given the s3Url
//func createPresignedURL(s3url string) {
//	// before getting a Pre-signed URL you must first create a Pre-signed client
//	fmt.Println("Create a pre-signed Client")
//	presignClient := s3.NewPre
//}

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

// looks like we have to use a presigned url to bypass access is denied :(
func downloadS3object(key string) {
	// The session the S3 Downloader will use
	//sess := session.Must(session.NewSession())
	//bucket, key := getBucketAnKey(get(key))
	//downloader := s3manager.NewDownloader(sess)
	//
	//downloadFile, error := os.Create("s3Url.txt")
	//if error != nil {
	//	log.Fatal("Failed to create new file ", error)
	//}
	//
	//_, error = downloader.Download(downloadFile, &s3.GetObjectInput{
	//	Bucket: aws.String(bucket),
	//	Key:    aws.String(key),
	//})
	//if error != nil {
	//	log.Fatal("Failed to download s3 object", error)
	//}
	//
	//err := downloadFile.Close()
	//if err != nil {
	//	return
	//}
}
