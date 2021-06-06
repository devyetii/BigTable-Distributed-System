package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Start/Stop of the server
var serving bool = false;
var max_tablet_cap int
var server_id int

func main() {
	// Load env vars
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// Get max cap
	max_tablet_cap, _ = strconv.Atoi(os.Getenv("MAX_TAB_CAP"))

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

	// Get server number
	server_id = SendServerIdRequest()
	fmt.Println(server_id)

	// Create the repository service and bind the update logger
	repo := Repository{data: BigTablePartition{}, keys: []RowKeyType{}, updateLogsFile: &update_logger}	
	
	log.Println("Tablet Server Started")
	
	// Create the API and bind the repo
	addr := fmt.Sprintf("localhost:%v", os.Getenv("PORT"))
	InitApi(addr, &repo, log_file)

}