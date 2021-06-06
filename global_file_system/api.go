package main

import (
	"io"
	"github.com/gofiber/fiber/v2"
  "github.com/gofiber/fiber/v2/middleware/logger"
	"strconv"
	"bytes"
	"bufio"
	"strings"
)

func InitApi(addr string, logFile io.Writer, result BigTablePartition) {
	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Output:     logFile,
		TimeFormat: "2006/01/02 15:04:05",
		Format:     "${time} ${status} - ${latency} ${method} ${path}\n",
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

	app.Put("/update-rows", func(c *fiber.Ctx) error {
		b := bytes.NewReader(c.Body())
		sc := bufio.NewScanner(b)
		var updates [][]string
		for sc.Scan() {
				words := strings.Fields(sc.Text())
				updates= append(updates,words)
		}
		updateTable(updates,result)
		return c.Status(200).SendString("ok")
	})
	app.Listen(addr)
}

func updateTable(updates [][]string,result BigTablePartition) {
  for _, update := range updates {
		key,_ := RowKeyFromString(update[1])
		if update[0] == "add_row" {

			result[key] = make(BigTableEntry)

		} else if update[0] == "delete_row" {

				delete(result, key)

		} else if update[0] == "delete_cell" {

			colKey := ColKeyType(update[2])
			result[key][colKey]=nil

		} else if update[0] == "set_cell" {

			colKey := ColKeyType(update[2])
			result[key][colKey]=update[3]

		}
	}
}