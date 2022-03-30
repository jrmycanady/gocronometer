package gocronometer

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type ServingRecord struct {
	RecordedTime     time.Time
	Group            string
	FoodName         string
	QuantityValue    float64
	QuantityUnits    string
	EnergyKcal       float64
	CaffeineMg       float64
	WaterG           float64
	B1Mg             float64
	B2Mg             float64
	B3Mg             float64
	B5Mg             float64
	B6Mg             float64
	B12Mg            float64
	BiotinUg         float64
	CholineMg        float64
	FolateUg         float64
	VitaminAUI       float64
	VitaminCMg       float64
	VitaminDUI       float64
	VitaminEMg       float64
	VitaminKMg       float64
	CalciumMg        float64
	ChromiumUg       float64
	CopperMg         float64
	FluorideUg       float64
	IodineUg         float64
	MagnesiumMg      float64
	ManganeseMg      float64
	PhosphorusMg     float64
	PotassiumMg      float64
	SeleniumUg       float64
	SodiumMg         float64
	ZincMg           float64
	CarbsG           float64
	FiberG           float64
	FructoseG        float64
	GalactoseG       float64
	GlucoseG         float64
	LactoseG         float64
	MaltoseG         float64
	StarchG          float64
	SucroseG         float64
	SugarsG          float64
	NetCarbsG        float64
	FatG             float64
	CholesterolMg    float64
	MonounsaturatedG float64
	PolyunsaturatedG float64
	SaturatedG       float64
	TransFatG        float64
	Omega3G          float64
	Omega6G          float64
	CystineG         float64
	HistidineG       float64
	IsoleucineG      float64
	LeucineG         float64
	LysineG          float64
	MethionineG      float64
	PhenylalanineG   float64
	ThreonineG       float64
	TryptophanG      float64
	TyrosineG        float64
	ValineG          float64
	ProtienG         float64
	IronMg           float64
	Category         string
}

type ServingRecords []ServingRecord

type ServingsExport struct {
	Records ServingRecords
}

func ParseServingsExport(rawCSVReader io.Reader, location *time.Location) (ServingRecords, error) {

	r := csv.NewReader(rawCSVReader)

	lineNum := 0
	headers := make(map[int]string)
	servings := make(ServingRecords, 0, 0)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Index all the headers.
		if lineNum == 0 {

			for i, v := range record {
				headers[i] = v
			}
			lineNum++
			continue
		}
		lineNum++

		var date string
		var timeStr string
		serving := ServingRecord{}
		for i, v := range record {
			columnName := headers[i]

			switch columnName {
			case "Day":
				date = v
			case "Time":
				timeStr = v
			case "Group":
				serving.Group = v
			case "Food Name":
				serving.FoodName = v
			case "Amount":
				s := strings.Split(v, " ")
				quantityValue, err := parseFloat(s[0], 64)
				if err != nil {
					return nil, fmt.Errorf("parsing quantity value: %s", err)
				}
				serving.QuantityValue = quantityValue
				serving.QuantityUnits = strings.Join(s[1:], "")

			case "Energy (kcal)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing energy: %s", err)
				}
				serving.EnergyKcal = f

			case "Caffeine (mg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing caffeine: %s", err)
				}
				serving.CaffeineMg = f

			case "Water (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing water: %s", err)
				}
				serving.WaterG = f

			case "B1 (Thiamine) (mg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing b1: %s", err)
				}
				serving.B1Mg = f

			case "B2 (Riboflavin) (mg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing b2: %s", err)
				}
				serving.B2Mg = f

			case "B3 (Niacin) (mg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing b3: %s", err)
				}
				serving.B3Mg = f

			case "B5 (Pantothenic Acid) (mg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing b5: %s", err)
				}
				serving.B5Mg = f

			case "B6 (Pyridoxine) (mg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing b6: %s", err)
				}
				serving.B6Mg = f

			case "B12 (Cobalamin) (µg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing b12: %s", err)
				}
				serving.B12Mg = f

			case "Biotin (µg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing biotin: %s", err)
				}
				serving.BiotinUg = f

			case "Choline (mg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing choline: %s", err)
				}
				serving.CholineMg = f

			case "Folate (µg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing folate: %s", err)
				}
				serving.FolateUg = f

			case "Vitamin A (IU)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman a: %s", err)
				}
				serving.VitaminAUI = f

			case "Vitamin C (mg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f

			case "Vitamin D (IU)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminDUI = f
			case "Vitamin E (mg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminEMg = f
			case "Vitamin K (µg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminKMg = f
			case "Calcium (mg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.CalciumMg = f
			case "Chromium (µg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.ChromiumUg = f
			case "Copper (mg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.CopperMg = f
			case "Fluoride (µg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.FluorideUg = f
			case "Iodine (µg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.IodineUg = f
			case "Iron (mg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.IronMg = f
			case "Magnesium (mg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.MagnesiumMg = f
			case "Manganese (mg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.ManganeseMg = f
			case "Phosphorus (mg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.PhosphorusMg = f
			case "Potassium (mg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.PotassiumMg = f
			case "Selenium (µg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.SeleniumUg = f
			case "Sodium (mg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.SodiumMg = f
			case "Zinc (mg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.ZincMg = f
			case "Carbs (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.CarbsG = f
			case "Fiber (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.FiberG = f
			case "Fructose (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.FructoseG = f
			case "Galactose (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.GalactoseG = f
			case "Glucose (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.GlucoseG = f
			case "Lactose (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.LactoseG = f
			case "Maltose (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.MaltoseG = f
			case "Starch (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.StarchG = f
			case "Sucrose (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.SucroseG = f
			case "Sugars (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.SugarsG = f
			case "Net Carbs (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.NetCarbsG = f
			case "Fat (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.FatG = f
			case "Cholesterol (mg)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.CholesterolMg = f
			case "Monounsaturated (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.MonounsaturatedG = f
			case "Polyunsaturated (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.PolyunsaturatedG = f
			case "Saturated (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.SaturatedG = f
			case "Trans-Fats (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.TransFatG = f
			case "Omega-3 (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.Omega3G = f
			case "Omega-6 (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.Omega6G = f
			case "Cystine (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.CystineG = f
			case "Histidine (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.HistidineG = f
			case "Isoleucine (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.IsoleucineG = f
			case "Leucine (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.LeucineG = f
			case "Lysine (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.LysineG = f
			case "Methionine (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.MethionineG = f
			case "Phenylalanine (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.PhenylalanineG = f
			case "Protein (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.ProtienG = f
			case "Threonine (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.ThreonineG = f
			case "Tryptophan (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.TryptophanG = f
			case "Tyrosine (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.TyrosineG = f
			case "Valine (g)":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.ValineG = f
			case "Category":
				serving.Category = v

			}
		}
		if timeStr == "" {
			timeStr = "00:00 AM"
		}

		if location == nil {
			location = time.UTC
		}
		serving.RecordedTime, err = time.ParseInLocation("2006-01-02 15:04 PM", date+" "+timeStr, location)
		if err != nil {
			return nil, fmt.Errorf("parsing record time: %s", err)
		}
		servings = append(servings, serving)
	}

	return servings, nil

}

// parseFloat wraps time.ParseFloat but interprites an empty string as 0.
func parseFloat(s string, bitSize int) (float64, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.ParseFloat(s, bitSize)
}

type ExerciseRecord struct {
	RecordedTime   time.Time
	Exercise       string
	Minutes        float64
	CaloriesBurned float64
}

type ExerciseRecords []ExerciseRecord

func ParseExerciseExport(rawCSVReader io.Reader, location *time.Location) (ExerciseRecords, error) {

	r := csv.NewReader(rawCSVReader)

	lineNum := 0
	headers := make(map[int]string)
	exercises := make(ExerciseRecords, 0, 0)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Index all the headers.
		if lineNum == 0 {

			for i, v := range record {
				headers[i] = v
			}
			lineNum++
			continue
		}
		lineNum++

		var date string
		var timeStr string
		exercise := ExerciseRecord{}
		for i, v := range record {
			columnName := headers[i]

			switch columnName {
			case "Day":
				date = v
			case "Time":
				timeStr = v
			case "Exercise":
				exercise.Exercise = v
			case "Minutes":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing energy: %s", err)
				}
				exercise.Minutes = f

			case "Calories Burned":
				f, err := parseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing caffeine: %s", err)
				}
				exercise.CaloriesBurned = f

			}
		}
		if timeStr == "" {
			timeStr = "00:00 AM"
		}

		if location == nil {
			location = time.UTC
		}
		exercise.RecordedTime, err = time.ParseInLocation("2006-01-02 15:04 PM", date+" "+timeStr, location)
		if err != nil {
			return nil, fmt.Errorf("parsing record time: %s", err)
		}
		exercises = append(exercises, exercise)
	}

	return exercises, nil

}

type BiometricRecord struct {
	RecordedTime time.Time
	Metric       string
	Unit         string
	Amount       float64
}

type BiometricRecords []BiometricRecord

func ParseBiometricRecordsExport(rawCSVReader io.Reader, location *time.Location) (BiometricRecords, error) {

	r := csv.NewReader(rawCSVReader)

	lineNum := 0
	headers := make(map[int]string)
	records := make(BiometricRecords, 0, 0)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Index all the headers.
		if lineNum == 0 {

			for i, v := range record {
				headers[i] = v
			}
			lineNum++
			continue
		}
		lineNum++

		var date string
		var timeStr string
		bioRecord := BiometricRecord{}
		for i, v := range record {
			columnName := headers[i]

			switch columnName {
			case "Day":
				date = v
			case "Time":
				timeStr = v
			case "Metric":
				bioRecord.Metric = v
			case "Unit":
				bioRecord.Unit = v
			case "Amount":
				if !strings.Contains(v, "/") {
					f, err := parseFloat(v, 64)
					if err != nil {
						return nil, fmt.Errorf("parsing energy: %s", err)
					}
					bioRecord.Amount = f
				}
			}
		}
		if timeStr == "" {
			timeStr = "00:00 AM"
		}

		if location == nil {
			location = time.UTC
		}
		bioRecord.RecordedTime, err = time.ParseInLocation("2006-01-02 15:04 PM", date+" "+timeStr, location)
		if err != nil {
			return nil, fmt.Errorf("parsing record time: %s", err)
		}
		records = append(records, bioRecord)
	}

	return records, nil

}
