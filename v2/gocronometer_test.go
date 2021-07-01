package v2_test

import (
	"context"
	"fmt"
	"github.com/jrmycanady/gocronometer/v2"
	"os"
	"testing"
)

// setup perform some basic actions to setup testing.
func setup() (username string, password string, client *v2.Client, err error) {
	username = os.Getenv("GOCRONOMETER_TEST_USERNAME")
	password = os.Getenv("GOCRONOMETER_TEST_PASSWORD")

	if username == "" {
		return "", "", nil, fmt.Errorf("username is empty, is GOCRONOMETER_TEST_USERNAME set?")
	}

	if password == "" {
		return "", "", nil, fmt.Errorf("password is empty, is GOCRONOMETER_TEST_PASSWORD set?")
	}

	return username, password, v2.NewClient(nil), nil
}

func TestClient_ObtainAntiCSRF(t *testing.T) {
	_, _, client, err := setup()
	if err != nil {
		t.Fatal(err)
	}

	antiCSRF, err := client.ObtainAntiCSRF(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if antiCSRF == "" {
		t.Fatalf("the anticsrf value was found to be empty")
	}
}

func TestClient_Login(t *testing.T) {
	username, password, client, err := setup()
	if err != nil {
		t.Fatal(err)
	}

	if err := client.Login(context.Background(), username, password); err != nil {
		t.Fatalf("failed to login: %s", err)
	}
}

func TestClient_Login_BadCreds(t *testing.T) {
	username, _, client, err := setup()
	if err != nil {
		t.Fatal(err)
	}

	if err := client.Login(context.Background(), username, "BAD"); err == nil {
		t.Fatalf("logged in with bad credentials")
	}
}
