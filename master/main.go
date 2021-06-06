package main

import (
	"os"
)

func main() {
	master_log_file, _ := os.OpenFile("master.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer master_log_file.Close()
	InitApi(":3001", master_log_file)

}
