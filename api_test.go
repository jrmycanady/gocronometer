package gocronometer_test

import (
	"context"
	"github.com/jrmycanady/gocronometer"
	"os"
	"testing"
	"time"
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

	t.Errorf("%d:%d", len(username), len(password))
	t.Errorf("%s:%s", username, password)

	c := gocronometer.NewClient()

	err := c.Login(context.Background(), username, password)
	if err != nil {
		t.Fatalf("failed to login with valid creds: %s", err)
	}
}

func TestAuthTokenGeneration(t *testing.T) {
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

	token, err := c.GenerateAuthToken(context.Background())
	if err != nil {
		t.Fatalf("failed to generate auth token: %s", err)
	}

	if token == "" {
		t.Fatalf("token was empty")
	}
}

func TestExportServings(t *testing.T) {
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

	_, err = c.ExportServings(context.Background(), time.Date(2020, 06, 01, 0, 0, 0, 0, time.UTC), time.Date(2020, 06, 04, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("failed to retrieve servings: %s", err)
	}
}

func TestExportDailyNutrition(t *testing.T) {
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

	_, err = c.ExportDailyNutrition(context.Background(), time.Date(2020, 06, 01, 0, 0, 0, 0, time.UTC), time.Date(2020, 06, 04, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("failed to retrieve servings: %s", err)
	}
}

func TestExportBiometrics(t *testing.T) {
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

	_, err = c.ExportBiometrics(context.Background(), time.Date(2020, 06, 01, 0, 0, 0, 0, time.UTC), time.Date(2020, 06, 04, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("failed to retrieve servings: %s", err)
	}
}

func TestExportExercises(t *testing.T) {
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

	_, err = c.ExportExercises(context.Background(), time.Date(2020, 06, 01, 0, 0, 0, 0, time.UTC), time.Date(2020, 06, 04, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("failed to retrieve servings: %s", err)
	}
}

func TestExportNotes(t *testing.T) {
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

	_, err = c.ExportNotes(context.Background(), time.Date(2020, 06, 01, 0, 0, 0, 0, time.UTC), time.Date(2020, 06, 04, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("failed to retrieve servings: %s", err)
	}
}
