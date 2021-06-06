package main

import (
	"io"

	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var metaData []Server

func InitApi(addr string, logFile io.Writer) {
	app := fiber.New()

	// compute iniital metadata and serve data
	tablets := assignDataToTablets()
	metaData = assignTabletsToServers(tablets)

	app.Use(logger.New(logger.Config{
		Output:     logFile,
		TimeFormat: "2006/01/02 15:04:05",
		Format:     "${time} ${status} - ${latency} ${method} ${path}\n",
	}))

	app.Get("/metadata", func(c *fiber.Ctx) error {

		data, _ := json.Marshal(metaData)
		return c.SendString(string(data))
	})

	app.Get("/load-balance-change", func(c *fiber.Ctx) error {
		// recompute metadata
		tablets := assignDataToTablets()
		metaData = assignTabletsToServers(tablets)
		data, _ := json.Marshal(metaData)
		return c.SendString(string(data))
	})

	app.Listen(addr)
}
