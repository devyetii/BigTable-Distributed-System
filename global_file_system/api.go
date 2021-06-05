package main

import (
	"io"
	"github.com/gofiber/fiber/v2"
  "github.com/gofiber/fiber/v2/middleware/logger"
	"strconv"
)

func InitApi(addr string, logFile io.Writer, result BigTablePartition) {
	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Output:     logFile,
		TimeFormat: "2006/01/02 15:04:05",
		Format:     "${time} ${status} - ${latency} ${method} ${path}",
	}))

	//NOTE: coupled with RowKeyType
	app.Get("/rows", func(c *fiber.Ctx) error {
		from := c.Query("from")
		to := c.Query("to")

		start,errStart := RowKeyFromString(from)
		end,errEnd := RowKeyFromString(to)
		if errStart != nil || errEnd != nil || start.UpperBound(end){
			return c.Status(400).SendString("Invalid from or to")
		}

		entries := make(BigTablePartition)
		for i:= start;i<=end;i++ {
			if v,ok := result[RowKeyType(i)]; ok {
				entries[RowKeyType(i)] = v
			}
		}
		return c.Status(200).JSON(entries)
	})

	app.Get("/rows-count", func(c *fiber.Ctx) error {
		return c.Status(200).SendString(strconv.Itoa(len(result)))
	})

	app.Listen(addr)
}