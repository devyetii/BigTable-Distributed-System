package main

import (
	"log"

	"math"

	"github.com/gofiber/fiber/v2"
)

/*
read data file json
divide data based on object count
assign to different tablets according to policy
write metadata
*/

const tabletSize = 100
const noOfServers = 2

type Tablet struct {
	id   int
	from int
	to   int
}

func getRowsCount() int {
	rowsCount := 5000
	return rowsCount
}

func assignDataToTablets() []Tablet {
	dataRowsCount := getRowsCount()
	noOfTablets := math.Ceil(float64(dataRowsCount) / float64(tabletSize))
	tabletsArr := make([]Tablet, 0, int(noOfTablets))
	for i := 0; i < int(noOfTablets); i++ {
		tabletsArr[i] = Tablet{
			id:   i,
			from: i * tabletSize,
			to:   int(math.Max(float64((i+1)*tabletSize-1), float64(dataRowsCount-1))),
		}
	}
	return tabletsArr

}

func assignTabletsToServers(tablets []Tablet) {
	noOfTabletsPerServer := int(math.Ceil(float64(len(tablets)) / noOfServers))
	j := 0
	for i := 0; i < noOfServers; i++ {
		tabletsArr := make([]Tablet, 0, noOfTabletsPerServer)
		k := 0
		if j < len(tablets) {
			tabletsArr[k] = tablets[j]
			j = j + 1
			k = k + 1
		}
	}
}
func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Fatal(app.Listen(":3000"))
}
