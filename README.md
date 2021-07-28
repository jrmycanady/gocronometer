# gocronometer
gocronometer is an GPLv2 licensed Go module that provides a client for exporting data from 
[Cronometer](https://cronometer.com). It utilizes the export features to retrieve the CSV data from the unpublished API.

**NOTE:** This module utilizes the same API the SPA uses. For that reason it should only be used by single users wanting 
to export their personal data for backup or other reasons. It should never be used for integrations that the enterprise 
plan would cover. The library is licensed under the GPLv2 to help prevent the unacceptable usage.

## Basic Example
```go
// Create the client.
c := gocronometer.NewClient(nil)

// Login to cronometer.
err := c.Login(context.Background(), username, password)
if err != nil {
    t.Fatalf("failed to login with valid creds: %s", err)
}

// Retrieve the export data.
rawCSVData, err = c.ExportServings(context.Background(), time.Date(2020, 06, 01, 0, 0, 0, 0, time.UTC), time.Date(2020, 06, 04, 0, 0, 0, 0, time.UTC))
if err != nil {
    t.Fatalf("failed to retrieve servings: %s", err)
}

fmt.Println(rawCSVData)
```

## Exports Supported

|method|description|
|------|-----------|
|ExportDailyNutrition()|Exports daily nutrition information for the date range provided.|
|ExportServings()|Exports servings for the date range provided.|
|ExportExercises()|Exports exercises for the date range provided.|
|ExportBiometrics()|Exports biometrics for the date range provided.|
|ExportNotes(|Exports notes for the date range provided.|

## API Magic Values

This library mimics the GWT HTTP requests to perform the export of data. The GWT API exposed by Cronometer is not 
designed to be accessed from anything besides their deployed GWT application. For that reason, there are several values 
that can only be obtained from loading the application itself. These values change over time with application updates,
and the library must use those new values.

The library includes the values as of the last push of the library. The new values can be provided to the client via
the ClientOptions parameter of the NewClient function.

### Magic Values

|Name|Location|Changes|
|----|------|------|
|GWTContentType|Retrieve from request header.|false|
|GWTModuleBase|Retrieve from request header.|false|
|GWTPermutation|Retrieve from request header.|true|
|GWTHeader|Retrieve from GWT request body.|true|

