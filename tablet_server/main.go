package main

import (
	"fmt"
	"log"
	"os"
)

func SetLogger() {
	fd, err := os.OpenFile("tablet_server.log", os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
	if (err != nil) {
		panic(err)
	}
	log.SetOutput(fd)
}

func main() {
	SetLogger()
	log.Println("Tablet Server Started")
	fmt.Println(data)
	InitApi(fmt.Sprintf("localhost:%v", os.Getenv("PORT")))
}