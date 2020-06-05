package gocronometer_test

import (
	"context"
	"gocronometer"
	"os"
	"testing"
)

func TestCRSFRetrieval(t *testing.T) {

	c := gocronometer.NewClient()

	_, err := c.RetrieveAntiCSRF(context.Background())
	if err != nil {
		t.Fatalf("failed to retrieve csrf: %s", err)
	}
}

func TestLogin(t *testing.T) {
	username := os.Getenv("GOCRONOMETER_TEST_USERNAME")
	password := os.Getenv("GOCRONOMETER_TEST_PASSWORD")

	if username == "" {
		t.Fatalf("username is empty, is GOCRONOMETER_TEST_USERNAME set?")
	}

	if password == "" {
		t.Fatalf("password is empty, is GOCRONOMETER_TEST_PASSWORD set?")
	}

	c := gocronometer.NewClient()

	err := c.Login(context.Background(), username, password)
	if err != nil {
		t.Fatalf("failed to login with valid creds: %s", err)
	}

}
