package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Persons struct {
	Persons []Person `json:"data"`
}

type Person struct {
	ID              int     `json:"id"`
	Gender          string  `json:"gender"`
	Age             int     `json:"age"`
	Hypertension    int     `json:"hypertension"`
	HeartDisease    int     `json:"heart_disease"`
	EverMarried     string  `json:"ever_married"`
	WorkType        string  `json:"work_type"`
	ResidenceType   string  `json:"Residence_type"`
	AvgGlucoseLevel float32 `json:"avg_glucose_level"`
	Bmi             string  `json:"bmi"`
	SmokingStatus   string  `json:"smoking_status"`
	Stroke          int     `json:"stroke"`
}

func main() {
	//Open the dataset file
	jsonFile, err := os.Open("dataset.json")

	//Check for errors
	if err != nil {
		fmt.Println("File reading error", err)
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	//Save the data in array of Person struct
	var persons Persons
	json.Unmarshal(byteValue, &persons)

}
