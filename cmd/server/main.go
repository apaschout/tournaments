package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/apaschout/tournaments"
	"github.com/cognicraft/event"
)

const (
	storeDSN  = "test.db"
	eventsDSN = "events.db"
	port      = ":5000"
)

func main() {
	es, err := event.NewStore(eventsDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer es.Close()

	db, err := tournaments.NewStore(storeDSN)
	if err != nil {
		log.Fatal(err)
	}

	s := tournaments.NewServer(db, es)
	fmt.Printf("Server running on %s\n", port)
	err = http.ListenAndServe(port, s)
	if err != nil {
		log.Fatal(err)
	}
}
