package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"math/rand"
)
type InputData []BigTableEntry
var maxInd RowKeyType
func generateRandomKey(i int) RowKeyType{
	c := i%5
	min := c*10000
	max := (c+1)*10000
	return RowKeyType(rand.Intn(max - min) + min)
}
func changeDataFormat() {
	//Open the dataset file
	jsonFile, err := os.Open("dataset.json")

	//Check for errors
	if err != nil {
		fmt.Println("File reading error", err)
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
  
	var result InputData
	json.Unmarshal([]byte(byteValue), &result)

	newTable :=make(BigTablePartition)

	for i,v := range result {
		x := generateRandomKey(i)
		for  {
			if _,ok := newTable[x]; ok {
				x = generateRandomKey(i)
			} else {
				break;
			}
		}
		if x>maxInd {
			maxInd = x
		}
		newTable[x]=v
	}

	dataFile, err := os.OpenFile("data.json", os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
	str, err := json.Marshal(newTable)

	dataFile.Write(str)
}
