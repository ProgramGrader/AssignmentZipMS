package aws_lambda_test

import (
	"aws_lambda"
	"common"
	"context"
	"dynamoDAO"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"os"
	"testing"
)

// we don't want to have to change this struct everytime we run a test
// implement both versions

type testCfg struct {
}

var cfg, _ = config.LoadDefaultConfig(context.TODO())
var ddbClient = dynamodb.NewFromConfig(cfg)
var awsClient = s3.NewFromConfig(cfg)

// For LocalStacks
var testEnvDbClient, _ = common.CreateDynamoDbLocalClient()
var testEnvAwsClient = common.CreateAwsConfig()

func setup() {

	// PLACING S3OBJ INTO S3BUCKET  -- Needs work object becomes blank once uploaded//
	file, err := os.OpenFile("url.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		fmt.Printf("Failed to open file\n")
		panic(err)
	}

	_, err = file.Write([]byte("https://iuscsg.org\n"))
	if err != nil {
		fmt.Printf("Failed to write to file")
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	bucket := common.BucketName

	//fileInfo, _ := file.Stat()
	//var size = fileInfo.Size()
	//buffer := make([]byte, size)
	//file.Read(buffer)
	//
	//
	//_, err = s3.PutObject(context.TODO(), &s3.PutObjectInput{
	//	Bucket: aws.String(bucket),
	//	Key:    aws.String(file.Name()),
	//	Body:   file,
	//})
	//if err != nil {
	//	log.Fatalln("Failed to upload file")
	//}

	//
	//out, err := awsClient.PutObject(context.TODO(), &objInput)
	//if err != nil {
	//	log.Fatalf("Failed to put Object: %v, %v", err, out)
	//}

	// GETTING S3OBJECT URL AND STORING IT INTO DYNAMODB //
	region := "us-east-2"
	s3ObjUrl := "https://%s.amazonaws.com/%s/%s"
	s3ObjUrl = fmt.Sprintf(s3ObjUrl, region, bucket, file.Name())

	dynamoDAO.Put(ddbClient, common.TableName, "UUID", bucket, region, file.Name())

	err = os.Remove("url.txt")
	if err != nil {
		fmt.Printf("Failed to delete file: %v\n", err)
		return
	}
}

func TestCreatePresignedURL(t *testing.T) {
	setup()
	bucket, region, filename := dynamoDAO.Get(ddbClient, common.TableName, "UUID")

	fmt.Println(bucket)
	fmt.Println(region)
	fmt.Println(filename)

	// The test was appending3 dashes to the end of the output causing a bad signature
	psUrl := aws_lambda.CreatePresignedURL(cfg, bucket, filename)
	print(psUrl + "\n")

	object, err := aws_lambda.GetURLObject(psUrl, "url.txt")
	if err != nil {
		t.Errorf("Failed to get url from S3 Object")
	}

	if object != "https://iuscsg.org" {
		t.Errorf("Expected: https://iuscsg.org, Got: %v", object)

	}
}
