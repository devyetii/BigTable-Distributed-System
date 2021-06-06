package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type HttpClient struct {
	updateLogger *SafeUpdateLog
}

func (client *HttpClient) GetDataFromGFS(from RowKeyType, to RowKeyType) BigTablePartition {
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

func (client *HttpClient) SendUpdatesToGFS() bool {
	client.updateLogger.mu.Lock()
	defer client.updateLogger.mu.Unlock()
	a := fiber.Post(fmt.Sprintf("%v/updates", os.Getenv("GFS_ADDR")))
	defer fiber.ReleaseAgent(a)

	a.BodyStream(client.updateLogger.GetFileForRead(), -1)
	
	if err := a.Parse(); err != nil {
		log.Println(err)
		return false
	}
	
	code, _, _ := a.String()
	
	client.updateLogger.ClearLogs()
	log.Println("Sent POST /updates to gfs")

	return code >= 200 && code < 400 
}

func (client *HttpClient) SendRebalanceRequest() {
	a := fiber.Get(fmt.Sprintf("%v/load-balance-change", os.Getenv("MASTER_ADDR")))
	defer fiber.ReleaseAgent(a)

	
	if err := a.Parse(); err != nil {
		log.Println(err)
	}
	
	log.Println("Sent GET /load-balance-change to master")
}

func (client *HttpClient) SendServerIdRequest() int {
	a := fiber.Get(fmt.Sprintf("%v/server-id", os.Getenv("MASTER_ADDR")))
	defer fiber.ReleaseAgent(a)

	if err := a.Parse(); err != nil {
		log.Println(err)
	}
	log.Println("Sent GET /server-id to master")
	_, sn, _ := a.String()
	
	isn, _ := strconv.Atoi(sn)

	return isn
}