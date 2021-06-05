package main

import (
	"encoding/json"
	"fmt"
)

/*
	Client logical steps
	1- fetch METADATA table from master at route ...
	2- schedule step 1 to run every X seconds
	3- ask user for which operation to do and the corresponding params
	4- apply mutex lock before every POST/PUT/DELETE request
	5- present data to user in any readable format
*/

func main() {
	res := getRows(3031, []RowKeyType{1, 2, 5}, "list")
	var resJSON BigTableEntry
	json.Unmarshal([]byte(res), &resJSON)
	person := resJSON["5"]
	fmt.Print(person)
}
