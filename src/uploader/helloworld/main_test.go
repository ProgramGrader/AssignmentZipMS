package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {

	t.Run("Non Post Request - Empty", func(t *testing.T) {

		test, err := handler(events.APIGatewayProxyRequest{})
		if err != nil {
			t.Fatalf("Error failed to trigger with an invalid HTTP response: %v", err)
		}
		assert.Equalf(t, test.StatusCode, 501, "Should return method not supported.")
		assert.Equalf(t, test.Body, "Method '' not allowed. Only Post is supported on this endpoint.", "Body should be the same")

	})

	t.Run("Non Post Request - Get", func(t *testing.T) {

		test, err := handler(events.APIGatewayProxyRequest{
			HTTPMethod: "GET",
		})

		if err != nil {
			t.Fatalf("Error failed to trigger with an invalid HTTP response: %v", err)
		}
		assert.Equalf(t, test.StatusCode, 501, "Should return method not supported.")
		assert.Equalf(t, test.Body, "Method 'GET' not allowed. Only Post is supported on this endpoint.", "Body should be the same")
	})

	t.Run("Non Post Request - Patch", func(t *testing.T) {

		test, err := handler(events.APIGatewayProxyRequest{
			HTTPMethod: "PATCH",
		})

		if err != nil {
			t.Fatalf("Error failed to trigger with an invalid HTTP response: %v", err)
		}
		assert.Equalf(t, test.StatusCode, 501, "Should return method not supported.")
		assert.Equalf(t, test.Body, "Method 'PATCH' not allowed. Only Post is supported on this endpoint.", "Body should be the same")
	})

	t.Run("Post Request - Without content", func(t *testing.T) {

		// Open our jsonFile
		jsonFile, err := os.Open("resources/requestWithoutContent.json")
		// if we os.Open returns an error then handle it
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Successfully Opened users.json")
		// defer the closing of our jsonFile so that we can parse it later on
		defer jsonFile.Close()

		var request events.APIGatewayProxyRequest

		byteValue, _ := ioutil.ReadAll(jsonFile)

		err = json.Unmarshal(byteValue, &request)
		if err != nil {
			t.Fatalf("Unable to marshal, %v", err)
		}

		test, err := handler(request)

		assert.Error(t, err)

		assert.Equalf(t, test.StatusCode, 400, "Should return method not supported.")
		assert.Equalf(t, test.Body, "Error occurred getting file from request: content type header missing.", "Body should be the same")
	})

	t.Run("Post Request - Without content", func(t *testing.T) {

		// Open our jsonFile
		jsonFile, err := os.Open("resources/request.json")
		// if we os.Open returns an error then handle it
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Successfully Opened users.json")
		// defer the closing of our jsonFile so that we can parse it later on
		defer jsonFile.Close()

		var request events.APIGatewayProxyRequest

		byteValue, _ := ioutil.ReadAll(jsonFile)

		err = json.Unmarshal(byteValue, &request)
		if err != nil {
			t.Fatalf("Unable to marshal, %v", err)
		}

		test, err := handler(request)

		assert.Error(t, err)
		assert.Equalf(t, test.StatusCode, 400, "Should return method not supported.")
		assert.Equalf(t, test.Body, "Error occurred getting file from request: content type header missing.", "Body should be the same")
	})

	t.Run("Post Request - Without content", func(t *testing.T) {

		fileDir, _ := os.Getwd()
		fileName := "resources/test.zip"
		filePath := path.Join(fileDir, fileName)

		file, _ := os.Open(filePath)
		defer file.Close()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
		io.Copy(part, file)
		writer.Close()

		r, _ := http.NewRequest("POST", "http://example.com", body)
		//r.Header.Add("Content-Type", writer.FormDataContentType())
		r.Header.Add("Content-Type", "application/zip")
		fmt.Println(r)
		//var test events.APIGatewayProxyRequest
		//
		//json.Unmarshal(r, &test)

		//var headers map[string]string
		//for Key, value := range r.Header {
		//	headers[Key] = value
		//}
		//
		//test := events.APIGatewayProxyRequest{
		//	HTTPMethod: "POST",
		//	Headers:    r.Header,
		//	Body:       r.Body,
		//}

		//fmt.Println(test)
		//client := &http.Client{}
		//client.Do(r)

	})

}
