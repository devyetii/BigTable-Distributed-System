package main

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
)


func InitApi(addr string) {
	app := fiber.New()

    app.Get("/rows", func(c *fiber.Ctx) error {
        if rng := c.Query("range"); rng != "" {
            range_parts := strings.Split(rng, "-")
            if entries := getByRange(range_parts[0], range_parts[1]); entries == nil {
                c.SendString("Invalid range")
                return c.SendStatus(400)
            } else {
                c.JSON(entries)
                return c.SendStatus(200)
            }
        } else if list := c.Query("list"); list != "" {
            keys_list := strings.Split(list, ",")
            c.JSON(getByKeysList(keys_list))
            return c.SendStatus(200)
        }

        c.SendString("Needs either range or list")
        return c.SendStatus(400)
    })

    app.Post("/row/:key", func(c *fiber.Ctx) error {
        row_key := c.Params("key")
        cols := make(BigTableEntry)
        if err := json.Unmarshal(c.Body(), &cols); err != nil {
            return c.SendStatus(400)
        }

        if row := addRow(row_key, cols); row != nil {
            c.JSON(row)
            return c.SendStatus(200)
        } else {
            c.SendString("Row already exists")
            return c.SendStatus(400)
        }
    })

    app.Put("/row/:key/cells", func(c *fiber.Ctx) error {
        row_key := c.Params("key")
        entry := make(BigTableEntry)
        if err := json.Unmarshal(c.Body(), &entry); err != nil {
            return c.SendStatus(400)
        }

        if row := setCells(row_key, entry); row != nil {
            c.JSON(row)
            return c.SendStatus(200)
        } else {
            c.SendString("Row not found")
            return c.SendStatus(404)
        }
    })

    app.Put("/row/:key/cells/delete", func(c *fiber.Ctx) error {
        row_key := c.Params("key")
        var col_keys []string

        if err := json.Unmarshal(c.Body(), &col_keys); err != nil {
            return c.SendStatus(400)
        }
        
        if row := deleteCells(row_key, col_keys); row != nil {
            c.JSON(row)
            return c.SendStatus(200)
        } else {
            c.SendString("Row not found")
            return c.SendStatus(404)
        }
    })

    app.Delete("/row/:key", func(c *fiber.Ctx) error {
        row_key := c.Params("key")

        if deleteRow(row_key) {
            c.SendString("Deleted")
            return c.SendStatus(200)
        } else {
            c.SendString("Row not found")
            return c.SendStatus(404)
        }
    })

    app.Listen(addr)
}