package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)
type InputTable []BigTableEntry
func main() {
	//Open the dataset file
	jsonFile, err := os.Open("data.json")

	//Check for errors
	if err != nil {
		fmt.Println("File reading error", err)
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result BigTablePartition
	json.Unmarshal([]byte(byteValue), &result)
	gfs_log_file, err := os.OpenFile("gfs.log", os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
	InitApi("localhost:3033",gfs_log_file,result)
}
