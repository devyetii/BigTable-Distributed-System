package main

import "bigtable_dist_sys/logger"


func main() {
	log := logger.NewLogger("out.txt")
	defer log.Close()

	log.WriteLine("Hi")
	log.WriteLine("I'm Ebrahim")
}