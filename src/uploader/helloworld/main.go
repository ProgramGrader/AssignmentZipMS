package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/grokify/go-awslambda"
	"github.com/spf13/afero"
	"io"
	"log"
	"path/filepath"
)

type FileRequest struct {
	fileName string
}

type customStruct struct {
	Content       string
	FileName      string
	FileExtension string
}

var (
	fs  = afero.NewOsFs()
	afs = &afero.Afero{Fs: fs}
)

func createReturnMessage(message string, code int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body: fmt.Sprintf("{\"message\":\"%v\"}", message),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: code,
	}
}

// application/zip
func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf("Processing Lambda req %s\n", req.RequestContext.RequestID)

	log.Println(req)

	ApiResponse := events.APIGatewayProxyResponse{}

	if req.HTTPMethod != "POST" {
		return createReturnMessage(
			fmt.Sprintf("Method '%v' not allowed. Only Post is supported on this endpoint.",
				req.HTTPMethod),
			405,
		), nil
	}

	// -----------------------
	// Create new reader
	log.Println("----!----! starting parse")
	r, err := awslambda.NewReaderMultipart(req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: fmt.Sprintf("{\"message\":\"Error occurred getting file from request: %v.\"}", err),
		}, nil
	}
	// Get the part
	part, err := r.NextPart()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: fmt.Sprintf("{\"message\":\"Error occurred getting file from request: %v.\"}", err),
		}, err
	}

	//Get file data
	content, err := io.ReadAll(part)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: fmt.Sprintf("{\"message\":\"Error occurred getting file from request: %v.\"}", err),
		}, err
	}

	custom := customStruct{
		Content:       string(content),
		FileName:      part.FileName(),
		FileExtension: filepath.Ext(part.FileName())}

	customBytes, err := json.Marshal(custom)
	log.Println(customBytes)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: fmt.Sprintf("{\"message\":\"Error occurred unmarshaling request: %v.\"}", err),
		}, err
	}

	if custom.FileExtension != "zip" {
		return createReturnMessage("Only zip files supported.", 400), nil
	}

	exists, err := afs.DirExists("/temp")

	if err != nil || !exists {
		return events.APIGatewayProxyResponse{}, err
	}

	// Get file into lambda

	err = afs.WriteFile("/temp/incomingzip.zip", []byte(custom.Content), 0777)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	log.Println("Stopping at current point.")
	// encode and hash file

	// put file in S3

	// Put hash, short hash and location into Dynamodb

	// Response
	return ApiResponse, nil

	//resp, err := http.Get(DefaultHTTPGetAddress)
	//if err != nil {
	//	return events.APIGatewayProxyResponse{}, err
	//}
	//
	//if resp.StatusCode != 200 {
	//	return events.APIGatewayProxyResponse{}, ErrNon200Response
	//}
	//
	//ip, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return events.APIGatewayProxyResponse{}, err
	//}
	//
	//if len(ip) == 0 {
	//	return events.APIGatewayProxyResponse{}, ErrNoIP
	//}
	//
	//return events.APIGatewayProxyResponse{
	//	Body:       fmt.Sprintf("Hello, %v", string(ip)),
	//	StatusCode: 200,
	//}, nil
}

func main() {
	lambda.Start(handler)
}
