package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/apaschout/tournaments"
	"github.com/cognicraft/event"
)

const (
	dsn           = ":memory:"
	dsnEventStore = "events.db"
	port          = ":8080"
)

func main() {
	es, err := event.NewStore(dsnEventStore)
	if err != nil {
		log.Fatal(err)
	}
	defer es.Close()
	sub := es.SubscribeToStreamFromCurrent(event.All)
	defer sub.Cancel()

	db, err := tournaments.NewDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	sub.On(db.On)

	s := tournaments.NewServer(db, es)
	s.Init()
	fmt.Printf("Server running on %s\n", port)
	err = http.ListenAndServe(port, s)
	if err != nil {
		log.Fatal(err)
	}
}
