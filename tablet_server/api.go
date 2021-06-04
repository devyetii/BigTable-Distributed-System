package main

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)


func String(c *fiber.Ctx, status int, msg string) error {
    c.SendString(msg)
    return c.SendStatus(status)
}

func InitApi(addr string, repo *Repository, logFile io.Writer) {
	app := fiber.New()

    app.Use(logger.New(logger.Config{
        Output: logFile,
        TimeFormat: "2006/01/02 15:04:05",
        Format: "${time} ${status} - ${latency} ${method} ${path}\n",
    }))

    app.Get("/rows", func(c *fiber.Ctx) error {
        if rng := c.Query("range"); rng != "" {
            // Extract range string
            range_parts := strings.Split(rng, "-")

            // Cast it to the row key type
            from, efrom := RowKeyFromString(range_parts[0]) 
            to, eto := RowKeyFromString(range_parts[1])

            // Format validation 
            if (efrom != nil || eto != nil) {
                return String(c, 400, "Invalid Range")
            }

            // Get and send values
            if entries := repo.getByRange(from, to); entries == nil {
                return String(c, 400, "Invalid Range")
            } else {
                return c.JSON(entries)
            }
        } else if list := c.Query("list"); list != "" {
            // Cast keys to row key type
            keys_list := MapStringsToRowKeys(strings.Split(list, ","), func (v string) (RowKeyType,error) { return RowKeyFromString(v) })
            if keys_list == nil {
                return String(c, 400, "Invalid Keys")
            }
            return c.JSON(repo.getByKeysList(keys_list))
        }

        return String(c, 400, "Needs either range or list")
    })

    app.Post("/row/:key", func(c *fiber.Ctx) error {
        row_key := c.Params("key")
        cols := make(BigTableEntry)
        if err := json.Unmarshal(c.Body(), &cols); err != nil {
            return c.SendStatus(400)
        }

        rk, err := RowKeyFromString(row_key)
        if (err != nil) {
            return String(c, 400, "Invalid Row Key")
        }
        if row := repo.addRow(rk, cols); row != nil {
            return c.JSON(row)
        } else {
            return String(c, 400, "Row already exists")
        }
    })

    app.Put("/row/:key/cells", func(c *fiber.Ctx) error {
        row_key := c.Params("key")
        entry := make(BigTableEntry)
        if err := json.Unmarshal(c.Body(), &entry); err != nil {
            return c.SendStatus(400)
        }

        rk, err := RowKeyFromString(row_key)
        if (err != nil) {
            return String(c, 400, "Invalid Row Key")
        }

        if row := repo.setCells(rk, entry); row != nil {
            return c.JSON(row)
        } else {
            return String(c, 400, "Row not found")
        }
    })

    app.Put("/row/:key/cells/delete", func(c *fiber.Ctx) error {
        row_key := c.Params("key")
        var col_keys []ColKeyType

        if err := json.Unmarshal(c.Body(), &col_keys); err != nil {
            return String(c, 400, "Invalid Column Keys")
        }
        
        rk, err := RowKeyFromString(row_key)
        if (err != nil) {
            return String(c, 400, "Invalid Row Key")
        }

        if row := repo.deleteCells(rk, col_keys); row != nil {
            return c.JSON(row)
        } else {
            return String(c, 400, "Row Not found")
        }
    })

    app.Delete("/row/:key", func(c *fiber.Ctx) error {
        row_key := c.Params("key")

        rk, err := RowKeyFromString(row_key)
        if (err != nil) {
            return String(c, 400, "Invalid Row Key")
        }

        if repo.deleteRow(rk) {
            return String(c, 200, "Deleted")
        } else {
            return String(c, 400, "Row not found")
        }
    })

    app.Listen(addr)
}