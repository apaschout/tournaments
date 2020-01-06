package main

import (
	"log"
	"net/http"

	"github.com/apaschout/tournaments"
)

const (
	dsn = "tournaments.db"
)

func main() {

	db, err := tournaments.NewDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = tournaments.InitDB(db.DB)
	if err != nil {
		log.Fatal(err)
	}

	s := tournaments.NewServer(db)
	s.Init()
	err = http.ListenAndServe(":8080", s)
	if err != nil {
		log.Fatal(err)
	}
}
