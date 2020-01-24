package tournaments

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/cognicraft/event"
	"github.com/cognicraft/sqlutil"
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
	CREATE TABLE IF NOT EXISTS seasons 
	(id TEXT PRIMARY KEY, name TEXT, start TEXT, end TEXT, format TEXT);`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS players (id TEXT PRIMARY KEY, name TEXT);`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS decks (id TEXT PRIMARY KEY, name TEXT, link TEXT);`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS metadata (key TEXT, value INTEGER);`)
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

func (db *DB) FindAllPlayers() ([]Player, error) {
	result := []Player{}
	query := "SELECT * FROM players"
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

func (db *DB) IsSeasonNameAvailable(name string) (bool, error) {
	err := db.checkNameDuplicate("Seasons", name)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *DB) IsPlayerNameAvailable(name string) (bool, error) {
	err := db.checkNameDuplicate("Players", name)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *DB) IsDeckNameAvailable(name string) (bool, error) {
	err = db.checkNameDuplicate("Decks", name)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *DB) PlayerExists(ID PlayerID) (bool, error) {
	err = db.checkID("Players", string(ID))
	if err != nil {
		return false, err
	}
	return true, nil
}

//throws error if ID does not exist
func (db *DB) checkID(table string, ID string) error {
	query := fmt.Sprintf("SELECT 0 FROM %s WHERE id = ?;", table)

	rows, err := db.Query(query, ID)
	if err != nil {
		err = fmt.Errorf("Checking for ID existence: %v", err)
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return nil
	}
	return fmt.Errorf("ID does not exist")
}

//throws error if name already exists
func (db *DB) checkNameDuplicate(table string, name string) error {
	query := fmt.Sprintf("SELECT 0 FROM %s WHERE name = ?;", table)

	rows, err := db.Query(query, name)
	if err != nil {
		err = fmt.Errorf("Checking for name availability: %v", err)
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return fmt.Errorf("Name already exists")
	}
	return nil
}

func (db *DB) On(rec event.Record) {
	codec, err := Codec()
	if err != nil {
		return
	}
	e, err := codec.Decode(rec)
	if err != nil {
		return
	}
	switch e := e.(type) {
	case SeasonCreated:
		err = sqlutil.Transact(db.DB, func(t *sql.Tx) error {
			query := "INSERT INTO players (id, name, start, end, format) VALUES (?, ?, ?, ?, ?)"
			_, err = db.Exec(query, e.Season, "", "", "", "")
			if err != nil {
				return err
			}
			return nil
		})
	case SeasonNameChanged:
		err = sqlutil.Transact(db.DB, func(t *sql.Tx) error {
			query := "UPDATE seasons SET name = ? WHERE id = ?;"
			_, err = db.Exec(query, e.Name, e.Season)
			if err != nil {
				return err
			}
			return nil
		})
	case PlayerCreated:
		err = sqlutil.Transact(db.DB, func(t *sql.Tx) error {
			query := "INSERT INTO players (id, name) VALUES (?, ?)"
			_, err = db.Exec(query, e.Player, 0, "")
			if err != nil {
				return err
			}
			return nil
		})
	case PlayerNameChanged:
		err = sqlutil.Transact(db.DB, func(t *sql.Tx) error {
			query := "UPDATE players SET name = ? WHERE id = ?;"
			_, err = db.Exec(query, e.Name, e.Player)
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err != nil {
		log.Println(err)
	}
}
