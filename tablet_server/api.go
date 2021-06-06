package main

import (
	"encoding/json"
	"fmt"
	// "fmt"
	"io"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/cors"
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
        Format: "${time} Handled ${status} - ${latency} ${method} ${path}\n",
    }))

    app.Use(cors.New())

    clientApi := app.Group("/row", func(c *fiber.Ctx) error {
        if !serving {
            return c.Status(fiber.StatusMethodNotAllowed).SendString("Not Ready")
        } else {
            return c.Next()
        }
    })

    app.Post("/serve", func(c *fiber.Ctx) error {
        // Recieve serve query
        var serveQuery ServeQueryType
        if err := json.Unmarshal(c.Body(), &serveQuery); err != nil {
            return String(c, 400, "Error in serve query")
        }

        // Initialize tablets, get data from GFS and add it
        for _, t := range serveQuery {
            // Get data
            data := repo.httpClient.GetDataFromGFS(t["From"], t["To"])
            if (data == nil) {
                return String(c, 400, "Error in GFS")
            }
            repo.AddData(data)

            // Assign tablet
            if (t["From"] == 0 && t["To"] == 0) {
                continue
            }
            currentTablet := Tablet{ from : t["From"], to: t["To"] }
            
            // Set tablet count
            currentTablet.count = len(data)

            // Assign new tablet
            repo.tablets = append(repo.tablets, &currentTablet)
        }
        serving = true
        return c.Status(201).JSON("Started Serving")
    })

    clientApi.Get("/", func(c *fiber.Ctx) error {
        if list := c.Query("list"); list != "" {
            // Cast keys to row key type
            keys_list := MapStringsToRowKeys(strings.Split(list, ","), func (v string) (RowKeyType,error) { return RowKeyFromString(v) })
            if keys_list == nil {
                return String(c, 400, "Invalid Keys")
            }
            return c.JSON(repo.getByKeysList(keys_list))
        }

        return String(c, 400, "Needs list")
    })

    clientApi.Post("/:key", func(c *fiber.Ctx) error {
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
            return String(c, 400, "Row range invalid")
        }
    })

    clientApi.Put("/:key/cells", func(c *fiber.Ctx) error {
        row_key := c.Params("key")
        entry := make(BigTableEntry)
        if err := json.Unmarshal(c.Body(), &entry); err != nil {
            return String(c, 400, "Invalid Column Keys")
        }

        rk, err := RowKeyFromString(row_key)
        if (err != nil) {
            return String(c, 400, "Invalid Row Key")
        }

        if row := repo.setCells(rk, entry); row != nil {
            return c.JSON(row)
        } else {
            return String(c, 404, "Row not found")
        }
    })

    clientApi.Put("/:key/cells/delete", func(c *fiber.Ctx) error {
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
            return String(c, 404, "Row Not found")
        }
    })

    clientApi.Delete("/", func(c *fiber.Ctx) error {
        if list := c.Query("list"); list != "" {
            // Cast keys to row key type
            keys_list := MapStringsToRowKeys(strings.Split(list, ","), func (v string) (RowKeyType,error) { return RowKeyFromString(v) })
            if keys_list == nil {
                return String(c, 400, "Invalid Keys")
            }
            delCount := repo.deleteRows(keys_list)
            return String(c, 200, fmt.Sprintf("Deleted %v rows", delCount))
        }

        return String(c, 400, "Needs list")
    })

    app.Listen(addr)
}