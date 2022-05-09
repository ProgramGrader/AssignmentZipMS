package tests

// All of this code is most likely useless, because since we switched to localstacks and
// after our discovery of the sam cli we can test all of these function as needed with
// a mock aws service

import (
	"testing"
)

const TABLENAME = "S3URLS"

func TestPutAndGet(t *testing.T) {

	type putItem struct {
		key      string
		value    string
		expected string
	}

	putItems := []putItem{
		{"key", "url", "url"},
		{"asefay", "update_url", "update_url"},
		{"urlId", "update_url", "update_url"},
	}

	for _, test := range putItems {
		Put(TABLENAME, test.key, test.value)
		storedItem := Get(TABLENAME, test.key)
		if storedItem != test.expected {
			t.Fatalf("TestPut(), Failed. Expected value was not found. Got %s expected %s", test.value, test.expected)
		}
	}

	DeleteAll(TABLENAME)
}

func TestDelete(t *testing.T) {
	type deleteItem struct {
		key      string
		value    string
		expected error
	}

	deleteItems := []deleteItem{
		{"key", "url", nil},
		{"asefay", "update_url", nil},
		{"urlId", "update_url", nil},
	}

	for _, test := range deleteItems {
		Put(TABLENAME, test.key, test.value)
		deleteErr := Delete(TABLENAME, test.key)
		if deleteErr != nil {
			t.Fatalf("TestDelete(), Failed. Expected error to be nil. Got %v", deleteErr)
		}
	}
}
