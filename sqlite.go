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
		err = rows.Scan(&seas.ID, &seas.Name, &seas.Start, &seas.End, &seas.Ongoing, &seas.Finished, &seas.Format)
		if err != nil {
			err = fmt.Errorf("Scan: %v", err)
			return nil, err
		}
		result = append(result, seas)
	}
	return result, nil
}

func (db *DB) FindSeasonByID(id SeasonID) (Season, error) {
	result := Season{}
	query := "SELECT * FROM Seasons WHERE id = ?"
	rows, err := db.Query(query, id)
	if err != nil {
		err = fmt.Errorf("Query: %v", err)
		return Season{}, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&result.ID, &result.Name, &result.Start, &result.End, &result.Ongoing, &result.Finished, &result.Format)
		if err != nil {
			err = fmt.Errorf("Scan: %v", err)
			return Season{}, err
		}
	}
	return result, nil
}

func (db *DB) SaveSeason(seas Season) error {
	query := "INSERT OR REPLACE INTO Seasons (id, name, start, end, ongoing, finished, format) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(query, seas.ID, seas.Name, seas.Start, seas.End, seas.Ongoing, seas.Finished, seas.Format)
	if err != nil {
		err = fmt.Errorf("Exec: %v", err)
		return err
	}
	err = db.SavePlayersToSeason(seas.ID, seas.Players)
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
		err = rows.Scan(&plr.ID, &plr.Name)
		if err != nil {
			err = fmt.Errorf("Scan: %v", err)
			return nil, err
		}
		result = append(result, plr)
	}
	return result, nil
}

func (db *DB) FindPlayerByID(id PlayerID) (Player, error) {
	result := Player{}
	query := "SELECT * FROM Players WHERE id = ?"
	rows, err := db.Query(query, id)
	if err != nil {
		err = fmt.Errorf("Query: %v", err)
		return Player{}, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&result.ID, &result.Name)
		if err != nil {
			err = fmt.Errorf("Scan: %v", err)
			return Player{}, err
		}
	}
	return result, nil
}

func (db *DB) SavePlayer(plr Player) error {
	query := "INSERT OR REPLACE INTO Players (id, name) VALUES (?, ?)"
	_, err := db.Exec(query, plr.ID, plr.Name)
	if err != nil {
		err = fmt.Errorf("Exec: %v", err)
		return err
	}
	log.Printf("Successfully saved Player with id:%s and name:%s", plr.ID, plr.Name)
	return nil
}

func (db *DB) FindPlayersInSeason(seasID SeasonID) ([]SeasonPlayer, error) {
	result := []SeasonPlayer{}
	query := `
			SELECT id
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
		var plrID PlayerID
		err = rows.Scan(&plrID)
		if err != nil {
			err = fmt.Errorf("Scan: %v", err)
			return nil, err
		}
		result = append(result, SeasonPlayer{ID: plrID})
	}
	return result, nil
}

func (db *DB) SavePlayersToSeason(seasID SeasonID, plrs []SeasonPlayer) error {
	query := "INSERT OR REPLACE INTO SeasonPlayers (seasonID, playerID) VALUES (?,?);"
	for _, sPlr := range plrs {
		_, err := db.Exec(query, seasID, sPlr.ID)
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
		err = rows.Scan(&deck.ID, &deck.Name, &deck.Link)
		if err != nil {
			err = fmt.Errorf("Scan: %v", err)
			return nil, err
		}
		result = append(result, deck)
	}
	return result, nil
}

func (db *DB) FindDeckByID(id DeckID) (Deck, error) {
	result := Deck{}
	query := "SELECT * FROM Decks WHERE id = ?;"
	rows, err := db.Query(query, id)
	if err != nil {
		err = fmt.Errorf("Query: %v", err)
		return Deck{}, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&result.ID, &result.Name, &result.Link)
		if err != nil {
			err = fmt.Errorf("Scan: %v", err)
			return Deck{}, err
		}
	}
	return result, nil
}

func (db *DB) SaveDeck(deck Deck) error {
	query := "INSERT OR REPLACE INTO Decks (id, name, link) VALUES (?, ?, ?);"
	_, err := db.Exec(query, deck.ID, deck.Name, deck.Link)
	if err != nil {
		err = fmt.Errorf("Exec: %v", err)
		return err
	}
	return nil
}

func (db *DB) SeasonNameAvailable(name string) (bool, error) {
	err := db.doNameQuery("Seasons", name)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *DB) PlayerNameAvailable(name string) (bool, error) {
	err := db.doNameQuery("Players", name)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *DB) DeckNameAvailable(name string) (bool, error) {
	err := db.doNameQuery("Decks", name)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *DB) doNameQuery(table string, name string) error {
	query := fmt.Sprintf("SELECT 0 FROM %s WHERE name = ?", table)

	rows, err := db.Query(query, name)
	if err != nil {
		err = fmt.Errorf("Error checking for name availability: %v", err)
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return fmt.Errorf("Name already exists")
	}
	return nil
}
