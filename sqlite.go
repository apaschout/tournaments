package tournaments

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/cognicraft/event"
	"github.com/cognicraft/sqlutil"
	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	dsn string
	db  *sql.DB
}

func NewStore(dsn string) (*Store, error) {
	s := &Store{dsn: dsn}
	return s, s.init()
}

func (s *Store) init() error {
	db, err := sql.Open("sqlite3", s.dsn)
	if err != nil {
		return err
	}
	s.db = db
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tournaments (id TEXT, name TEXT, start TEXT, end TEXT, format TEXT, PRIMARY KEY (id));`)
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
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS registered_players (tournament TEXT, player TEXT, PRIMARY KEY (tournament, player));`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS metadata (key TEXT PRIMARY KEY, value INTEGER);`)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) FindAllTournaments() ([]Tournament, error) {
	result := []Tournament{}
	query := "SELECT id, name FROM tournaments"
	rows, err := s.db.Query(query)
	if err != nil {
		err = fmt.Errorf("Query: %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		trn := Tournament{}
		err = rows.Scan(&trn.ID, &trn.Name)
		if err != nil {
			err = fmt.Errorf("Scan: %v", err)
			return nil, err
		}
		result = append(result, trn)
	}
	return result, nil
}

func (s *Store) FindTournamentByID(id TournamentID) (Tournament, error) {
	result := Tournament{}
	query := "SELECT * FROM tournaments WHERE id = ?"
	rows, err := s.db.Query(query, id)
	if err != nil {
		err = fmt.Errorf("Query: %v", err)
		return Tournament{}, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&result.ID, &result.Name, &result.Start, &result.End, &result.Format)
		if err != nil {
			err = fmt.Errorf("Scan: %v", err)
			return Tournament{}, err
		}
	}
	return result, nil
}

func (s *Store) FindAllPlayers() ([]Player, error) {
	result := []Player{}
	query := "SELECT id, name FROM players"
	rows, err := s.db.Query(query)
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

func (s *Store) FindPlayerByID(id PlayerID) (Player, error) {
	result := Player{}
	query := "SELECT * FROM Players WHERE id = ?"
	rows, err := s.db.Query(query, id)
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

func (s *Store) FindAllDecks() ([]Deck, error) {
	result := []Deck{}
	query := "SELECT * FROM Decks ;"
	rows, err := s.db.Query(query)
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

func (s *Store) FindDeckByID(id DeckID) (Deck, error) {
	result := Deck{}
	query := "SELECT * FROM Decks WHERE id = ?;"
	rows, err := s.db.Query(query, id)
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

func (s *Store) IsTournamentNameAvailable(name string) (bool, error) {
	err := s.checkNameDuplicate("tournaments", name)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Store) IsPlayerNameAvailable(name string) (bool, error) {
	err := s.checkNameDuplicate("players", name)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Store) IsDeckNameAvailable(name string) (bool, error) {
	err = s.checkNameDuplicate("decks", name)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Store) PlayerExists(ID PlayerID) (bool, error) {
	err = s.checkID("players", string(ID))
	if err != nil {
		return false, err
	}
	return true, nil
}

//throws error if ID does not exist
func (s *Store) checkID(table string, ID string) error {
	query := fmt.Sprintf("SELECT 0 FROM %s WHERE id = ?;", table)

	rows, err := s.db.Query(query, ID)
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
func (s *Store) checkNameDuplicate(table string, name string) error {
	query := fmt.Sprintf("SELECT 0 FROM %s WHERE name = ?;", table)

	rows, err := s.db.Query(query, name)
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

func (s *Store) GetVersion() uint64 {
	var res uint64
	query := `SELECT value FROM metadata WHERE key = "version"`
	rows, err := s.db.Query(query)
	if err != nil {
		log.Println(err)
		return 0
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&res)
		if err != nil {
			log.Println(err)
			return 0
		}
	}
	return res
}

func (s *Store) On(rec event.Record) {
	codec, err := Codec()
	if err != nil {
		log.Println(err)
		return
	}
	e, err := codec.Decode(rec)
	if err != nil {
		log.Println(err)
		return
	}
	switch e := e.(type) {
	case TournamentCreated:
		err = sqlutil.Transact(s.db, func(t *sql.Tx) error {
			query := "INSERT INTO tournaments (id, name) VALUES (?, ?);"
			_, err = t.Exec(query, e.Tournament, e.Tournament)
			if err != nil {
				log.Printf("%v", err)
				return err
			}
			log.Println("Projection: TournamentCreated")
			return nil
		})
	case TournamentDeleted:
		err = sqlutil.Transact(s.db, func(t *sql.Tx) error {
			query := "DELETE FROM tournaments WHERE id = ?;"
			_, err = t.Exec(query, e.Tournament)
			if err != nil {
				log.Printf("%v", err)
				return err
			}
			log.Println("Projection: TournamentDeleted")
			return nil
		})
	case TournamentNameChanged:
		err = sqlutil.Transact(s.db, func(t *sql.Tx) error {
			query := "UPDATE tournaments SET name = ? WHERE id = ?;"
			_, err = t.Exec(query, e.Name, e.Tournament)
			if err != nil {
				return err
			}
			log.Println("Projection: TournamentNameChanged")
			return nil
		})
	case TournamentPlayerRegistered:
		err = sqlutil.Transact(s.db, func(t *sql.Tx) error {
			query := "INSERT INTO registered_players (tournament, player) VALUES (?, ?);"
			_, err = t.Exec(query, e.Tournament, e.Player)
			if err != nil {
				return err
			}
			log.Println("Projection: TournamentPlayerRegistered")
			return nil
		})
	case TournamentPlayerDropped:
		err = sqlutil.Transact(s.db, func(t *sql.Tx) error {
			query := "DELETE FROM registered_players WHERE tournament = ? AND player = ?;"
			_, err = t.Exec(query, e.Tournament, e.Player)
			if err != nil {
				return err
			}
			log.Println("Projection: TournamentPlayerDropped")
			return nil
		})
	case PlayerCreated:
		err = sqlutil.Transact(s.db, func(t *sql.Tx) error {
			query := "INSERT INTO players (id, name) VALUES (?, ?);"
			_, err = t.Exec(query, e.Player, e.Player)
			if err != nil {
				return err
			}
			log.Println("Projection: PlayerCreated")
			return nil
		})
	case PlayerNameChanged:
		err = sqlutil.Transact(s.db, func(t *sql.Tx) error {
			query := "UPDATE players SET name = ? WHERE id = ?;"
			_, err = t.Exec(query, e.Name, e.Player)
			if err != nil {
				return err
			}
			log.Println("Projection: PlayerNameChanged")
			return nil
		})
	}
	err = sqlutil.Transact(s.db, func(t *sql.Tx) error {
		query := `INSERT OR REPLACE INTO metadata (key, value) VALUES ("version", ?)`
		_, err := t.Exec(query, s.GetVersion()+1)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
}
