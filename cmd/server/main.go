package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/apaschout/tournaments"
	"github.com/cognicraft/event"
)

const (
	dsn           = "tournaments.db"
	dsnEventStore = "events.db"
	port          = ":8080"
)

func main() {
	store, err := event.NewStore(dsnEventStore)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	db, err := tournaments.NewDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	s := tournaments.NewServer(db, store)
	s.Init()
	fmt.Printf("Server running on %s\n", port)
	err = http.ListenAndServe(port, s)
	if err != nil {
		log.Fatal(err)
	}
}
