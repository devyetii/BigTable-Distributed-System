package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
)

// Start/Stop of the server
var serving bool = false
var max_tablet_cap int
var server_id int = -1

var update_logger SafeUpdateLog
var httpClient HttpClient

func cron(httpClient *HttpClient) {
	gocron.Every(1).Minute().Do(httpClient.SendUpdatesToGFS)
	<- gocron.Start()
}

func justPrintErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

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
	update_logs_file, err := os.OpenFile("updates.log", os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0644)
	if (err != nil) {
		panic(err)
	}
	defer update_logs_file.Close()
	update_logger = SafeUpdateLog{file: update_logs_file}
	httpClient = HttpClient{ updateLogger: &update_logger }

	// Get server number
	for server_id = httpClient.SendServerIdRequest(); server_id < 0; {
		log.Println("Waiting for master")
		time.Sleep(5 * time.Second)
	}
	log.Println(fmt.Sprintf("Got id %v from master", server_id))

	// Create the repository service and bind the update logger
	repo := Repository{data: BigTablePartition{}, keys: []RowKeyType{}, httpClient: &httpClient, updateLogsFile: &update_logger}	
	
	// Run the update logger sender cron
	go cron(&httpClient)

	log.Println("Tablet Server Started")
	
	// Create the API and bind the repo
	addr := fmt.Sprintf("localhost:%v", os.Getenv("PORT"))
	InitApi(addr, &repo, log_file)
}