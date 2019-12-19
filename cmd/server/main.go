package main

import (
	"log"
	"net/http"

	"github.com/apaschout/tournaments"
)

func main() {

	s := tournaments.NewServer()
	s.Init()
	err := http.ListenAndServe(":8080", s)
	if err != nil {
		log.Fatal(err)
	}
}
