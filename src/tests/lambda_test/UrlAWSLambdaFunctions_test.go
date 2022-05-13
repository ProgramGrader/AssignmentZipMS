package lambda_test_test

import (
	"common"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"os"
	"testing"
	"tests"
	"tests/lambda_test"
)

var dynamodbClient, _ = common.CreateDynamoDbLocalClient()
var awsCfg = common.CreateAwsConfig()

func setup() {

	// PLACING S3OBJ INTO S3BUCKET //
	client := s3.NewFromConfig(awsCfg)

	file, err := os.OpenFile("url", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		fmt.Printf("Failed to open file\n")
		panic(err)
	}

	_, err = file.Write([]byte("https://iuscsg.org\n"))
	if err != nil {
		fmt.Printf("Failed to write to file")
		panic(err)
	}

	bucket := common.BucketName
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

	tests.Put(dynamodbClient, common.TableName, "UUID", region, bucket, file.Name())

	err = os.Remove("url")
	if err != nil {
		fmt.Printf("Failed to delete file: %v\n", err)
		return
	}
}

func TestCreatePresignedURL(t *testing.T) {
	setup()
	bucket, region, filename := tests.Get(dynamodbClient, common.TableName, "UUID")

	fmt.Println(bucket)
	fmt.Println(region)
	fmt.Println(filename)

	psUrl := lambda_test.CreatePresignedURL(awsCfg, bucket, filename)
	print(psUrl)

}

func TestDownloadS3Object(t *testing.T) {

}
