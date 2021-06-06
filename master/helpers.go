package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
)

const tabletSize = 100
const noOfServers = 2

var GFSEndPoint string = "localhost:3033"
var TabletServerEndPoint string = "localhost:3036"

func getRowsCount() int {
	log.Println("get max row index from GFS")
	response, err := http.Get(GFSEndPoint + "/rows-count")
	if err != nil {
		log.Fatal(err)
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	rowsCount, _ := strconv.Atoi(bodyString)
	return rowsCount
}

func assignDataToTablets() []Tablet {
	log.Println("assign data to tablets")
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
	log.Println("assign tablets to availabe tablet servers")
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

func serveRequestServer(serverAddress string, data []byte) {
	log.Println("send serve request to tablet server " + serverAddress)

	responseBody := bytes.NewBuffer(data)
	//Leverage Go's HTTP Post function to make request
	_, err := http.Post(serverAddress+"/serve", "application/json", responseBody)
	if err != nil {
		log.Println(err)
	}

}
