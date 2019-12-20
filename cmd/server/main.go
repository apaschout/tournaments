package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/apaschout/tournaments"
)

const (
	dsn = "tournaments.db"
)

func main() {

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = initDB(db)
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

func initDB(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS "Seasons" ("id" INTEGER PRIMARY KEY,"name" TEXT UNIQUE)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "Players" ("id" INTEGER PRIMARY KEY,"name" TEXT UNIQUE)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "SeasonPlayers" ("seasonID" INTEGER,"playerID" INTEGER,
		FOREIGN KEY("playerID") REFERENCES "Players"("id"),
		FOREIGN KEY("seasonID") REFERENCES "Seasons"("id"),
		PRIMARY KEY("seasonID","playerID")
	)`)
	if err != nil {
		return err
	}
	return nil
}
