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
	ProteinG         float64
	ThreonineG       float64
	TryptophanG      float64
	TyrosineG        float64
	ValineG          float64
	Category         string
}

type ServingRecords []ServingRecord

type ServingsExport struct {
	Records ServingRecords
}

func ParseServingsExport(rawCSVReader io.Reader) (ServingRecords, error) {

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
			continue
		}

		for i, v := range record {
			columnName := headers[i]
			serving := ServingRecord{}
			var date string
			var timeStr string

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
				quantityValue, err := strconv.ParseFloat(s[0], 64)
				if err != nil {
					return nil, fmt.Errorf("parsing quantity value: %s", err)
				}
				serving.QuantityValue = quantityValue
				serving.QuantityUnits = strings.Join(s[1:], "")

			case "Energy (kcal)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing energy: %s", err)
				}
				serving.EnergyKcal = f

			case "Caffeine (mg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing caffeine: %s", err)
				}
				serving.CaffeineMg = f

			case "Water (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing water: %s", err)
				}
				serving.WaterG = f

			case "B1 (Thiamine) (mg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing b1: %s", err)
				}
				serving.B1Mg = f

			case "B2 (Riboflavin) (mg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing b2: %s", err)
				}
				serving.B2Mg = f

			case "B3 (Niacin) (mg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing b3: %s", err)
				}
				serving.B3Mg = f

			case "B5 (Pantothenic Acid) (mg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing b5: %s", err)
				}
				serving.B5Mg = f

			case "B6 (Pyridoxine) (mg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing b6: %s", err)
				}
				serving.B6Mg = f

			case "B12 (Cobalamin) (µg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing b12: %s", err)
				}
				serving.B12Mg = f

			case "Biotin (µg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing biotin: %s", err)
				}
				serving.BiotinUg = f

			case "Choline (mg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing choline: %s", err)
				}
				serving.CholineMg = f

			case "Folate (µg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing folate: %s", err)
				}
				serving.FolateUg = f

			case "Vitamin A (IU)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman a: %s", err)
				}
				serving.VitaminAUI = f

			case "Vitamin C (mg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f

			case "Vitamin D (IU)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Vitamin E (mg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Vitamin K (µg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Calcium (mg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Chromium (µg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Copper (mg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Fluoride (µg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Iodine (µg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Iron (mg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Magnesium (mg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Manganese (mg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Phosphorus (mg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Potassium (mg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Selenium (µg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Sodium (mg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Zinc (mg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Carbs (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Fiber (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Fructose (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Galactose (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Glucose (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Lactose (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Maltose (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Starch (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Sucrose (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Sugars (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Net Carbs (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Fat (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Cholesterol (mg)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Monounsaturated (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Polyunsaturated (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Saturated (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Trans-Fats (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Omega-3 (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Omega-6 (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Cystine (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Histidine (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Isoleucine (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Leucine (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Lysine (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Methionine (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Phenylalanine (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Protein (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Threonine (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Tryptophan (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Tyrosine (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Valine (g)":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f
			case "Category":
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("parsing vitiman c: %s", err)
				}
				serving.VitaminCMg = f

			}

			serving.RecordedTime, err = time.Parse("2006-01-02 15:04 PM", date+" "+timeStr)
			if err != nil {
				return nil, fmt.Errorf("parsing record time: %s", err)
			}

			servings = append(servings, serving)
		}
	}

	return servings, nil

}
