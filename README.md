# gocronometer
gocronometer is an MIT licensed Go module that provides a client for exporting data from 
[Cronometer](https://cronometer.com). It utilizes the export features to retrieve the CSV data from the unpublished API.

**NOTE:** This module utilizes the same API the SPA uses. For that reason it should only be used by single users wanting 
to export their personal data for backup or other reasons. It should never be used for integrations that the enterprise 
plan would cover. 

## Example
```go
// Create the client.
c := gocronometer.NewClient()

// Login to cronometer.
err := c.Login(context.Background(), username, password)
if err != nil {
    t.Fatalf("failed to login with valid creds: %s", err)
}

// Retrieve the export data.
csvData, err = c.ExportServings(context.Background(), time.Date(2020, 06, 01, 0, 0, 0, 0, time.UTC), time.Date(2020, 06, 04, 0, 0, 0, 0, time.UTC))
if err != nil {
    t.Fatalf("failed to retrieve servings: %s", err)
}

fmt.Println(csvData)
```
## Exports Supported

|method|description|
|------|-----------|
|ExportDailyNutrition()|Exports daily nutrition information for the date range provided.|
|ExportServings()|Exports servings for the date range provided.|
|ExportExercises()|Exports exercises for the date range provided.|
|ExportBiometrics()|Exports biometrics for the date range provided.|
|ExportNotes(|Exports notes for the date range provided.|
