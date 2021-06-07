package main

import (
	"fmt"
	"io"
	"log"

	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var metaData []Server
var currServerId int = -1
var hashIPMap map[int]string = make(map[int]string)

func InitApi(addr string, logFile io.Writer) {
	app := fiber.New()

	// compute iniital metadata and serve data
	log.Println("metadata initial calculation")
	tablets := assignDataToTablets()
	metaData = assignTabletsToServers(tablets)

	app.Use(logger.New(logger.Config{
		Output:     logFile,
		TimeFormat: "2006/01/02 15:04:05",
		Format:     "${time} ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New())

	app.Get("/metadata", func(c *fiber.Ctx) error {

		data, _ := json.Marshal(metaData)
		return c.SendString(string(data))
	})

	app.Get("/load-balance-change", func(c *fiber.Ctx) error {
		// recompute metadata
		c1 := make(chan int, 1)
		tablets := assignDataToTablets()
		metaData = assignTabletsToServers(tablets)
		data, _ := json.Marshal(metaData)
		// send serve request to all servers with new metadata
		for _, element := range hashIPMap {
			go serveRequestServer(element, data, c1)
		}
		return c.SendString(string(data))
	})

	app.Get("/server-id", func(c *fiber.Ctx) error {
		c1 := make(chan int, 1)
		currServerId = (currServerId + 1) % noOfServers
		// add to hash map id
		hashIPMap[currServerId] = c.Query("serverAddress")
		// send serve request
		fmt.Println(hashIPMap[currServerId])
		data, _ := json.Marshal(metaData[currServerId].Tablets)
		go serveRequestServer(hashIPMap[currServerId], data, c1)
		return c.SendString(fmt.Sprint(currServerId))
	})
	app.Get("/serve/:id", func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		data, _ := json.Marshal(metaData[id].Tablets)
		return c.SendString(string(data))
	})
	app.Listen(addr)
}
