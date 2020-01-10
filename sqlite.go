package tournaments

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func NewDB(dsn string) (*DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	store := &DB{db}
	return store, store.init()
}

func (db *DB) init() error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS "Seasons" 
	("id" TEXT PRIMARY KEY,"name" TEXT, "start" TEXT, "end" TEXT, "finished" INTEGER, "format" TEXT);`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "Players" ("id" TEXT PRIMARY KEY,"name" TEXT);`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "Decks" ("id" TEXT PRIMARY KEY,"name" TEXT, "link" TEXT);`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS "SeasonPlayers" ("seasonID" TEXT,"playerID" TEXT,
		FOREIGN KEY("playerID") REFERENCES "Players"("id"),
		FOREIGN KEY("seasonID") REFERENCES "Seasons"("id"),
		PRIMARY KEY("seasonID","playerID"));`)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) FindAllSeasons() ([]Season, error) {
	result := []Season{}
	query := "SELECT * FROM Seasons"
	rows, err := db.Query(query)
	if err != nil {
		err = fmt.Errorf("Query: %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		seas := Season{}
		err = rows.Scan(&seas.Id, &seas.Name, &seas.Start, &seas.End, &seas.Finished, &seas.Format)
		if err != nil {
			err = fmt.Errorf("Scan: %v", err)
			return nil, err
		}
		result = append(result, seas)
	}
	return result, nil
}

func (db *DB) FindSeasonByID(id string) (Season, error) {
	result := Season{}
	query := "SELECT * FROM Seasons WHERE id = ?"
	rows, err := db.Query(query, id)
	if err != nil {
		err = fmt.Errorf("Query: %v", err)
		return Season{}, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&result.Id, &result.Name, &result.Start, &result.End, &result.Finished, &result.Format)
		if err != nil {
			err = fmt.Errorf("Scan: %v", err)
			return Season{}, err
		}
	}
	return result, nil
}

func (db *DB) SaveSeason(seas Season) error {
	query := "INSERT OR REPLACE INTO Seasons (id, name, start, end, finished, format) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(query, seas.Id, seas.Name, seas.Start, seas.End, seas.Finished, seas.Format)
	if err != nil {
		err = fmt.Errorf("Exec: %v", err)
		return err
	}
	err = db.SavePlayersToSeason(seas.Id, seas.Players)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) FindAllPlayers() ([]Player, error) {
	result := []Player{}
	query := "SELECT * FROM Players"
	rows, err := db.Query(query)
	if err != nil {
		err = fmt.Errorf("Query: %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		plr := Player{}
		err = rows.Scan(&plr.Id, &plr.Name)
		if err != nil {
			err = fmt.Errorf("Scan: %v", err)
			return nil, err
		}
		result = append(result, plr)
	}
	return result, nil
}

func (db *DB) FindPlayerByID(id string) (Player, error) {
	result := Player{}
	query := "SELECT * FROM Players WHERE id = ?"
	rows, err := db.Query(query, id)
	if err != nil {
		err = fmt.Errorf("Query: %v", err)
		return Player{}, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&result.Id, &result.Name)
		if err != nil {
			err = fmt.Errorf("Scan: %v", err)
			return Player{}, err
		}
	}
	return result, nil
}

func (db *DB) SavePlayer(plr Player) error {
	query := "INSERT OR REPLACE INTO Players (id, name) VALUES (?, ?)"
	_, err := db.Exec(query, plr.Id, plr.Name)
	if err != nil {
		err = fmt.Errorf("Exec: %v", err)
		return err
	}
	log.Printf("Successfully saved Player with id:%s and name:%s", plr.Id, plr.Name)
	return nil
}

func (db *DB) FindPlayersInSeason(seasID string) ([]Player, error) {
	result := []Player{}
	query := `
			SELECT id, name
			FROM Players
			INNER JOIN SeasonPlayers ON SeasonPlayers.playerID = Players.id
			WHERE seasonID = ?`
	rows, err := db.Query(query, seasID)
	if err != nil {
		err = fmt.Errorf("Query: %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		plr := Player{}
		err = rows.Scan(&plr.Id, &plr.Name)
		if err != nil {
			err = fmt.Errorf("Scan: %v", err)
			return nil, err
		}
		result = append(result, plr)
	}
	return result, nil
}

func (db *DB) SavePlayersToSeason(seasID string, plrs []Player) error {
	query := "INSERT OR REPLACE INTO SeasonPlayers (seasonID, playerID) VALUES (?,?);"
	for _, p := range plrs {
		_, err := db.Exec(query, seasID, p.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) FindAllDecks() ([]Deck, error) {
	result := []Deck{}
	query := "SELECT * FROM Decks ;"
	rows, err := db.Query(query)
	if err != nil {
		err = fmt.Errorf("Query : %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		deck := Deck{}
		err = rows.Scan(&deck.Id, &deck.Name, &deck.Link)
		if err != nil {
			err = fmt.Errorf("Scan: %v", err)
			return nil, err
		}
		result = append(result, deck)
	}
	return result, nil
}

func (db *DB) FindDeckByID(id string) (Deck, error) {
	result := Deck{}
	query := "SELECT * FROM Decks WHERE id = ?;"
	rows, err := db.Query(query, id)
	if err != nil {
		err = fmt.Errorf("Query: %v", err)
		return Deck{}, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&result.Id, &result.Name, &result.Link)
		if err != nil {
			err = fmt.Errorf("Scan: %v", err)
			return Deck{}, err
		}
	}
	return result, nil
}

func (db *DB) SaveDeck(deck Deck) error {
	query := "INSERT OR REPLACE INTO Decks (id, name, link) VALUES (?, ?, ?);"
	_, err := db.Exec(query, deck.Id, deck.Name, deck.Link)
	if err != nil {
		err = fmt.Errorf("Exec: %v", err)
		return err
	}
	return nil
}

func (db *DB) CheckForDuplicateName(table string, name string) error {
	query := fmt.Sprintf("SELECT id FROM %s WHERE name = ?", table)

	rows, err := db.Query(query, name)
	if err != nil {
		err = fmt.Errorf("Query: %v", err)
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return fmt.Errorf("Name already exists")
	}
	return nil
}
