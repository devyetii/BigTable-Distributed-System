package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type RowKeyType int
type ColKeyType string
type ValType interface{}
type BigTableEntry map[ColKeyType]ValType
type BigTablePartition map[RowKeyType]BigTableEntry

func getRows(port int, searchKeys []RowKeyType, searchType string) string {
	agent := fiber.AcquireAgent()
	req := agent.Request()
	sendArray := strings.Trim(strings.Replace(fmt.Sprint(searchKeys), " ", ",", -1), "[]")
	URL := "http://localhost:" + strconv.Itoa(port) + "/rows?" + searchType + "=" + sendArray
	req.SetRequestURI(URL)
	errorSend := agent.Parse()
	if errorSend != nil {
		fmt.Print(errorSend)
	}

	_, body, errorBytes := agent.Bytes()
	if errorBytes != nil {
		fmt.Print(errorBytes)
	}

	return string(body)
}

func addRow(port int, searchKeys []RowKeyType, searchType string) string {
	agent := fiber.AcquireAgent()
	req := agent.Request()
	sendArray := strings.Trim(strings.Replace(fmt.Sprint(searchKeys), " ", ",", -1), "[]")
	URL := "http://localhost:" + strconv.Itoa(port) + "/rows?" + searchType + "=" + sendArray
	req.SetRequestURI(URL)
	errorSend := agent.Parse()
	if errorSend != nil {
		fmt.Print(errorSend)
	}

	_, body, errorBytes := agent.Bytes()
	if errorBytes != nil {
		fmt.Print(errorBytes)
	}

	return string(body)
}
