package main

import (
	"fmt"
	"math"
	"net/http"
)

const tabletSize = 100
const noOfServers = 2

var GFSEndPoint string = "localhost:3033"

func getRowsCount() int {
	fmt.Println(GFSEndPoint)
	response, _ := http.Get(GFSEndPoint + "/rows-count")
	fmt.Println(response.Body)
	rowsCount := 500 //strconv.Atoi( response.
	//print(rowsCount)
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

func serveRequestServer(serverAddress string, data string) {

}