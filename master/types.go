package main

type Tablet struct {
	Id   int
	From int
	To   int
}
type Server struct {
	Id      int
	Tablets []Tablet
}
