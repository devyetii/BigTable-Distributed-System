package main

import (
	"io"
	"github.com/gofiber/fiber/v2"
  "github.com/gofiber/fiber/v2/middleware/logger"
	"strconv"
	"fmt"
	"bytes"
	"bufio"
	"strings"
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

	app.Put("/update-rows", func(c *fiber.Ctx) error {
		b := bytes.NewReader(c.Body())
		sc := bufio.NewScanner(b)
		var updates [][]string
		for sc.Scan() {
				words := strings.Fields(sc.Text())
				updates= append(updates,words)
		}
		updateTable(updates)
		return c.Status(200).SendString("ok")
	})
	app.Listen(addr)
}

func updateTable(updates [][]string) {
	var addRow [][]string
	var deleteRow [][]string
	var addCell [][]string
	var deleteCell [][]string
	var setCell [][]string
  for _, update := range updates {
		if update[0] == "add_row" {
			addRow = append(addRow,update)
		} else if update[0] == "delete_row" {
			deleteRow = append(deleteRow,update)
		} else if update[0] == "add_cell" {
			addCell = append(addCell,update)
		} else if update[0] == "delete_cell" {
			deleteCell = append(deleteCell,update)
		} else {
			setCell = append(setCell,update)
		}
	}
}