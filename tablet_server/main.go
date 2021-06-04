package main

import (
	"fmt"
	"log"
	"os"
)

// MOCK DATA
// TODO Remove 
var keys []string = []string{
    "a.person.com",
    "b.person.com",
    "b.persona.com",
}

var data BigTablePartition = BigTablePartition{
	"a.person.com" : BigTableEntry{
		"name" : "ebrahim",
		"age" : "22",
	},
	"b.person.com" : BigTableEntry{
		"name" : "Farha",
		"age" : "23",
	},
	"b.persona.com" : BigTableEntry{
		"name" : "Mahmoud",
		"age" : "21",
	},
};

func main() {
	// Setup logger
	log_file, err := os.OpenFile("tablet_server.log", os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
	if (err != nil) {
		panic(err)
	}
	defer log_file.Close()
	log.SetOutput(log_file)

	// Setup update logs
	update_logs_file, err := os.OpenFile("updates.log", os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
	if (err != nil) {
		panic(err)
	}
	defer update_logs_file.Close()
	update_logger := SafeUpdateLog{file: *update_logs_file}

	// Create the repository service and bind the update logger
	repo := Repository{data: data, keys: keys, updateLogsFile: &update_logger}
	
	// Create the API and bind the repo
	addr := fmt.Sprintf("localhost:%v", os.Getenv("PORT"))
	InitApi(addr, &repo, log_file)

	log.Println("Tablet Server Started")
}