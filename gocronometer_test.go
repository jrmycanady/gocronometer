package gocronometer_test

import (
	"context"
	"fmt"
	"github.com/jrmycanady/gocronometer"
	"os"
	"testing"
	"time"
)

// setup perform some basic actions to setup testing.
func setup() (username string, password string, client *gocronometer.Client, err error) {
	username = os.Getenv("GOCRONOMETER_TEST_USERNAME")
	password = os.Getenv("GOCRONOMETER_TEST_PASSWORD")

	if username == "" {
		return "", "", nil, fmt.Errorf("username is empty, is GOCRONOMETER_TEST_USERNAME set")
	}

	if password == "" {
		return "", "", nil, fmt.Errorf("password is empty, is GOCRONOMETER_TEST_PASSWORD set")
	}

	return username, password, gocronometer.NewClient(nil), nil
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

	defer client.Logout(context.Background())
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

func TestClient_GenerateAuthToken(t *testing.T) {
	username, password, client, err := setup()
	if err != nil {
		t.Fatal(err)
	}

	if err := client.Login(context.Background(), username, password); err != nil {
		t.Fatalf("failed to login: %s", err)
	}

	defer client.Logout(context.Background())

	token, err := client.GenerateAuthToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if token == "" {
		t.Fatalf("GWT auth token was empty")
	}

	if len(token) > len("2f90aabe5493a07a6d9ab4a17b9ea65e") {
		t.Fatalf("token was %s", token)
	}
}

func TestClient_ExportBiometrics(t *testing.T) {
	username, password, client, err := setup()
	if err != nil {
		t.Fatal(err)
	}

	if err := client.Login(context.Background(), username, password); err != nil {
		t.Fatalf("failed to login: %s", err)
	}

	defer client.Logout(context.Background())

	startTime := time.Date(2021, 6, 1, 0, 0, 0, 0, time.Local)
	endTime := time.Date(2021, 6, 10, 0, 0, 0, 0, time.Local)

	_, err = client.ExportBiometrics(context.Background(), startTime, endTime)
	if err != nil {
		t.Fatalf("failed to export bio: %s", err)
	}

}

func TestClient_ExportDailyNutrition(t *testing.T) {
	username, password, client, err := setup()
	if err != nil {
		t.Fatal(err)
	}

	if err := client.Login(context.Background(), username, password); err != nil {
		t.Fatalf("failed to login: %s", err)
	}

	defer client.Logout(context.Background())

	startTime := time.Date(2021, 6, 1, 0, 0, 0, 0, time.Local)
	endTime := time.Date(2021, 6, 10, 0, 0, 0, 0, time.Local)

	_, err = client.ExportDailyNutrition(context.Background(), startTime, endTime)
	if err != nil {
		t.Fatalf("failed to export bio: %s", err)
	}
}

func TestClient_ExportNotes(t *testing.T) {
	username, password, client, err := setup()
	if err != nil {
		t.Fatal(err)
	}

	if err := client.Login(context.Background(), username, password); err != nil {
		t.Fatalf("failed to login: %s", err)
	}

	defer client.Logout(context.Background())

	startTime := time.Date(2021, 6, 1, 0, 0, 0, 0, time.Local)
	endTime := time.Date(2021, 6, 10, 0, 0, 0, 0, time.Local)

	_, err = client.ExportNotes(context.Background(), startTime, endTime)
	if err != nil {
		t.Fatalf("failed to export bio: %s", err)
	}

}

func TestClient_ExportServings(t *testing.T) {
	username, password, client, err := setup()
	if err != nil {
		t.Fatal(err)
	}

	if err := client.Login(context.Background(), username, password); err != nil {
		t.Fatalf("failed to login: %s", err)
	}

	defer client.Logout(context.Background())

	startTime := time.Date(2021, 6, 1, 0, 0, 0, 0, time.Local)
	endTime := time.Date(2021, 6, 10, 0, 0, 0, 0, time.Local)

	_, err = client.ExportServings(context.Background(), startTime, endTime)
	if err != nil {
		t.Fatalf("failed to export bio: %s", err)
	}

}

func TestClient_ExportExercises(t *testing.T) {
	username, password, client, err := setup()
	if err != nil {
		t.Fatal(err)
	}

	if err := client.Login(context.Background(), username, password); err != nil {
		t.Fatalf("failed to login: %s", err)
	}

	defer client.Logout(context.Background())

	startTime := time.Date(2021, 6, 1, 0, 0, 0, 0, time.Local)
	endTime := time.Date(2021, 6, 10, 0, 0, 0, 0, time.Local)

	_, err = client.ExportExercises(context.Background(), startTime, endTime)
	if err != nil {
		t.Fatalf("failed to export bio: %s", err)
	}
}

func TestClient_ExportServingsParsed(t *testing.T) {
	username, password, client, err := setup()
	if err != nil {
		t.Fatal(err)
	}

	if err := client.Login(context.Background(), username, password); err != nil {
		t.Fatalf("failed to login: %s", err)
	}

	defer client.Logout(context.Background())

	startTime := time.Date(2021, 6, 1, 0, 0, 0, 0, time.Local)
	endTime := time.Date(2021, 6, 10, 0, 0, 0, 0, time.Local)

	_, err = client.ExportServingsParsed(context.Background(), startTime, endTime)
	if err != nil {
		t.Fatalf("failed to export bio: %s", err)
	}

}

func TestClient_ExportExercisesParsed(t *testing.T) {
	username, password, client, err := setup()
	if err != nil {
		t.Fatal(err)
	}

	if err := client.Login(context.Background(), username, password); err != nil {
		t.Fatalf("failed to login: %s", err)
	}

	defer client.Logout(context.Background())

	startTime := time.Date(2021, 6, 1, 0, 0, 0, 0, time.Local)
	endTime := time.Date(2021, 6, 10, 0, 0, 0, 0, time.Local)

	_, err = client.ExportExercisesParsedWithLocation(context.Background(), startTime, endTime, time.UTC)
	if err != nil {
		t.Fatalf("failed to export bio: %s", err)
	}

}

func TestClient_ExportBiometricRecordsParsed(t *testing.T) {
	username, password, client, err := setup()
	if err != nil {
		t.Fatal(err)
	}

	if err := client.Login(context.Background(), username, password); err != nil {
		t.Fatalf("failed to login: %s", err)
	}

	defer client.Logout(context.Background())

	startTime := time.Date(2021, 6, 1, 0, 0, 0, 0, time.Local)
	endTime := time.Date(2021, 6, 10, 0, 0, 0, 0, time.Local)

	_, err = client.ExportBiometricRecordsParsedWithLocation(context.Background(), startTime, endTime, time.UTC)
	if err != nil {
		t.Fatalf("failed to export bio: %s", err)
	}

}
