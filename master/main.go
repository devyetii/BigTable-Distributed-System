package main

import (
	"encoding/json"
	"log"
	"math"

	"github.com/gofiber/fiber/v2"
)

const tabletSize = 100
const noOfServers = 2

type Tablet struct {
	Id   int
	From int
	To   int
}
type Server struct {
	Id      int
	Tablets []Tablet
}

func getRowsCount() int {
	rowsCount := 500
	return rowsCount
}

func assignDataToTablets() []Tablet {
	dataRowsCount := getRowsCount()
	noOfTablets := math.Ceil(float64(dataRowsCount) / float64(tabletSize))
	tabletsArr := make([]Tablet, int(noOfTablets))
	for i := 0; i < int(noOfTablets); i++ {
		tabletsArr[i] = Tablet{
			Id:   i,
			From: i * tabletSize,
			To:   int(math.Min(float64((i+1)*tabletSize-1), float64(dataRowsCount-1))),
		}
	}
	return tabletsArr

}

func assignTabletsToServers(tablets []Tablet) []Server {
	noOfTabletsPerServer := int(math.Ceil(float64(len(tablets)) / noOfServers))
	servers := make([]Server, noOfServers)
	j := 0
	for i := 0; i < noOfServers; i++ {
		tabletsArr := make([]Tablet, noOfTabletsPerServer)

		for k := 0; k < noOfTabletsPerServer; k++ {
			if j >= len(tablets) {
				break
			}
			tabletsArr[k] = tablets[j]
			j = j + 1

		}

		servers[i] = Server{
			Tablets: tabletsArr,
			Id:      i,
		}
	}
	return servers
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		Tablets := assignDataToTablets()
		metaData := assignTabletsToServers(Tablets)
		data, _ := json.Marshal(metaData)
		return c.SendString(string(data))
	})

	log.Fatal(app.Listen(":3001"))

}
