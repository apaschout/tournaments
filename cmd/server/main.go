package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/apaschout/tournaments"
)

const (
	dsn  = "tournaments.db"
	port = ":8080"
)

func main() {

	db, err := tournaments.NewDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	s := tournaments.NewServer(db)
	s.Init()
	fmt.Printf("Server running on %s\n", port)
	err = http.ListenAndServe(port, s)
	if err != nil {
		log.Fatal(err)
	}
}
