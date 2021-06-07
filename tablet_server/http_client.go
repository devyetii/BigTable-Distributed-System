package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type HttpClient struct {
	updateLogger *SafeUpdateLog
}

func checkResponseErrors(code int, errors []error) bool {
	if (code != fiber.StatusOK || len(errors) > 0) {
		log.Println(fmt.Sprintf("Errors with code: %v", code))
		for err := range errors {
			log.Println(err)
		}
		return true
	}
	return false
}

func (client *HttpClient) GetDataFromGFS(from RowKeyType, to RowKeyType) BigTablePartition {
	log.Println(fmt.Sprintf("Sending GET /rows to GFS from %v to %v", from, to))
	a := fiber.Get(fmt.Sprintf("%v/rows", os.Getenv("GFS_ADDR")))
	defer fiber.ReleaseAgent(a)
	a.QueryString(fmt.Sprintf("from=%v&to=%v", from, to))

	if err := a.Parse(); err != nil {
		log.Println(fmt.Sprintf("Errors in GFS Request: %v", err))
		return nil
	}
	log.Println(fmt.Sprintf("Sent GET /rows to GFS from %v to %v", from, to))

	var result BigTablePartition
	code, body, errors := a.Bytes()

	if checkResponseErrors(code, errors) {
		return nil
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Println(fmt.Sprintf("Errors in GFS Unmarshal: %v", err))
	}
	return result
}

func (client *HttpClient) SendUpdatesToGFS() bool {
	if (!serving) {
		return false
	}

	log.Println("Sending POST /updates to gfs")
	client.updateLogger.mu.Lock()
	defer client.updateLogger.mu.Unlock()
	a := fiber.Post(fmt.Sprintf("%v/updates", os.Getenv("GFS_ADDR")))
	defer fiber.ReleaseAgent(a)

	a.BodyStream(client.updateLogger.GetFileForRead(), -1)
	
	if err := a.Parse(); err != nil {
		log.Println(err)
		return false
	}
	log.Println("Sent POST /updates to gfs")

	code, _, errors := a.String()

	if checkResponseErrors(code, errors) {
		return false
	}
	
	client.updateLogger.ClearLogs()
	
	return code >= 200 && code < 400 
}

func (client *HttpClient) SendRebalanceRequest() {
	log.Println("Sending GET /load-balance-change to master")
	a := fiber.Get(fmt.Sprintf("%v/load-balance-change", os.Getenv("MASTER_ADDR")))
	defer fiber.ReleaseAgent(a)

	
	if err := a.Parse(); err != nil {
		log.Println(err)
		return
	}
	log.Println("Sent GET /load-balance-change to master")
	
	code, _, errors := a.Bytes()
	
	if checkResponseErrors(code, errors) {
		return
	}

}

func (client *HttpClient) SendServerIdRequest() int {
	log.Println("Sending GET /server-id to master")
	a := fiber.Get(fmt.Sprintf("%v/server-id", os.Getenv("MASTER_ADDR")))
	a.QueryString(fmt.Sprintf("serverAddress=%v", os.Getenv("SELF_ADDR")))
	defer fiber.ReleaseAgent(a)

	if err := a.Parse(); err != nil {
		log.Println(err)
		return -1
	}
	log.Println("Sent GET /server-id to master")

	code, sn, errors := a.String()
	
	if checkResponseErrors(code, errors) {
		return -1
	}

	isn, _ := strconv.Atoi(sn)

	return isn
}