package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"math/rand"
)
type InputTable []BigTableEntry
func changedata() {
	//Open the dataset file
	jsonFile, err := os.Open("dataset.json")

	//Check for errors
	if err != nil {
		fmt.Println("File reading error", err)
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
  
	var result InputTable
	json.Unmarshal([]byte(byteValue), &result)

	newTable :=make(BigTablePartition)

	for i,v := range result {
		c := i%5
    min := c*1000
    max := (c+1)*1000
		newTable[RowKeyType(rand.Intn(max - min) + min)]=v
	}

	dataFile, err := os.OpenFile("data.json", os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
	str, err := json.Marshal(newTable)

	dataFile.Write(str)

}
