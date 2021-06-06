package main

import (
	"log"
	"os"
)

func main() {
	GFSEndPoint = os.Args[2]
	port := os.Args[1]
	master_log_file, _ := os.OpenFile("master.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer master_log_file.Close()
	InitApi(":"+port, master_log_file)
	log.Println("Master Server Started")
	getRowsCount()

}
