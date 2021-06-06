package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetDataFromGFS(from RowKeyType, to RowKeyType) BigTablePartition {
	a := fiber.Get(fmt.Sprintf("%v/rows", os.Getenv("GFS_ADDR")))
	defer fiber.ReleaseAgent(a)
	a.QueryString(fmt.Sprintf("from=%v&to=%v", from, to))


	if err := a.Parse(); err != nil {
		return nil
	}
	log.Println(fmt.Sprintf("Sent GET /rows to GFS from %v to %v", from, to))

	var result BigTablePartition
	code, _, _ := a.Struct(&result)
	
	if (code != fiber.StatusOK) {
		return nil
	}
	return result	
}

func SendUpdatesToGFS() bool {
	a := fiber.Post(fmt.Sprintf("%v/updates", os.Getenv("GFS_ADDR")))
	defer fiber.ReleaseAgent(a)

	update_log_file, err := os.OpenFile("updates.log", os.O_RDWR, 0644)
	if (err != nil) {
		log.Panic(err)
		return false
	}
	defer update_log_file.Close()

	a.BodyStream(update_log_file, -1)

	if err := a.Parse(); err != nil {
		log.Println(err)
		return false
	}

	code, _, _ := a.String()

	update_log_file.Truncate(100)

	return code >= 200 && code < 400 
}

func SendRebalanceRequest() {
	a := fiber.Get(fmt.Sprintf("%v/load-balance-change", os.Getenv("MASTER_ADDR")))

	if err := a.Parse(); err != nil {
		log.Println(err)
	}
}

func SendServerIdRequest() int {
	a := fiber.Get(fmt.Sprintf("%v/load-balance-change", os.Getenv("MASTER_ADDR")))

	if err := a.Parse(); err != nil {
		log.Println(err)
	}

	_, sn, _ := a.String()
	
	isn, _ := strconv.Atoi(sn)

	return isn
}