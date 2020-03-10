package tournaments

import (
	"fmt"
	"log"
	"time"

	"github.com/cognicraft/event"
	"github.com/cognicraft/uuid"
)

type PlayerCreated struct {
	ID         string    `json:"id"`
	OccurredOn time.Time `json:"occured-on"`
	Player     PlayerID  `json:"player"`
}

type PlayerNameChanged struct {
	ID         string    `json:"id"`
	OccurredOn time.Time `json:"occurred-on"`
	Player     PlayerID  `json:"player"`
	Name       string    `json:"name"`
}

type PlayerMatchPlayed struct {
	ID         string    `json:"id"`
	OccurredOn time.Time `json:"occurred-on"`
	Player     PlayerID  `json:"player"`
}

type PlayerMatchWon struct {
	ID         string    `json:"id"`
	OccurredOn time.Time `json:"occurred-on"`
	Player     PlayerID  `json:"player"`
}

type PlayerGamePlayed struct {
	ID         string    `json:"id"`
	OccurredOn time.Time `json:"occurred-on"`
	Player     PlayerID  `json:"player"`
}

type PlayerGameWon struct {
	ID         string    `json:"id"`
	OccurredOn time.Time `json:"occurred-on"`
	Player     PlayerID  `json:"player"`
}

type PlayerTournamentRegistered struct {
	ID         string       `json:"id"`
	OccurredOn time.Time    `json:"occurred-on"`
	Player     PlayerID     `json:"player"`
	Tournament TournamentID `json:"Tournament"`
}

func NewPlayer() *Player {
	return &Player{
		ChangeRecorder: event.NewChangeRecorder(),
	}
}

func (plr *Player) Create(id PlayerID) error {
	if plr.ID != "" {
		return fmt.Errorf("Player already exists")
	}
	if id == "" {
		return fmt.Errorf("A Player's ID may not be empty")
	}
	plr.Apply(PlayerCreated{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Player:     id,
	})
	log.Printf("Event: Player %s: Created\n", plr.ID)
	return nil
}

func (plr *Player) ChangeName(name string) error {
	if plr.ID == "" {
		return fmt.Errorf("Player does not exist")
	}
	if name == "" {
		return fmt.Errorf("A Player's name may not be empty")
	}
	if plr.Name == name {
		return nil
	}
	plr.Apply(PlayerNameChanged{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Player:     plr.ID,
		Name:       name,
	})
	log.Printf("Event: Player %s: Name Changed To %s\n", plr.ID, name)
	return nil
}

func (plr *Player) IncrementMatches() error {
	if plr.ID == "" {
		return fmt.Errorf("Player does not exist")
	}
	plr.Apply(PlayerMatchPlayed{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Player:     plr.ID,
	})
	log.Printf("Event: Player %s: Total Matches Incremented\n", plr.ID)
	return nil
}

func (plr *Player) IncrementMatchesWon() error {
	if plr.ID == "" {
		return fmt.Errorf("Player does not exist")
	}
	plr.Apply(PlayerMatchWon{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Player:     plr.ID,
	})
	log.Printf("Event: Player %s: Matches Won Incremented\n", plr.ID)
	return nil
}

func (plr *Player) IncrementGames() error {
	if plr.ID == "" {
		return fmt.Errorf("Player does not exist")
	}
	plr.Apply(PlayerGamePlayed{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Player:     plr.ID,
	})
	log.Printf("Event: Player %s: Total Games incremented\n", plr.ID)
	return nil
}

func (plr *Player) IncrementGamesWon() error {
	if plr.ID == "" {
		return fmt.Errorf("Player does not exist")
	}
	plr.Apply(PlayerGameWon{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Player:     plr.ID,
	})
	log.Printf("Event: Player %s: Games Won incremented\n", plr.ID)
	return nil
}

func (plr *Player) RegisterTournament(tID TournamentID) error {
	if plr.ID == "" {
		return fmt.Errorf("Player does not exist")
	}
	if tID == "" {
		return fmt.Errorf("No Tournament specified")
	}
	plr.Apply(PlayerTournamentRegistered{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Player:     plr.ID,
		Tournament: tID,
	})
	log.Printf("Event: Player %s: Tournament %s registered", plr.ID, tID)
	return nil
}

func (plr *Player) Apply(e event.Event) {
	plr.Record(e)
	plr.Mutate(e)
}

func (plr *Player) Mutate(e event.Event) {
	plr.Version++
	switch e := e.(type) {
	case PlayerCreated:
		plr.ID = e.Player
		plr.Name = string(e.Player)
	case PlayerNameChanged:
		plr.Name = e.Name
	case PlayerMatchPlayed:
		plr.MatchesPlayed++
	case PlayerMatchWon:
		plr.MatchesWon++
	case PlayerGamePlayed:
		plr.GamesPlayed++
	case PlayerGameWon:
		plr.GamesWon++
	case PlayerTournamentRegistered:
		plr.Tournaments = append(plr.Tournaments, e.Tournament)
	}
}

func (plr *Player) Save(es *event.Store, metadata interface{}) error {
	if len(plr.Changes()) == 0 {
		return nil
	}
	streamID := string(plr.ID)
	exp := plr.Version - uint64(len(plr.Changes()))
	codec, err := Codec()
	if err != nil {
		return err
	}
	recs, err := codec.EncodeAll(plr.Changes(), event.WithMetadata(metadata))
	if err != nil {
		return err
	}
	err = es.Append(streamID, exp, recs)
	if err != nil {
		return err
	}
	plr.ClearChanges()
	return nil
}

func LoadPlayer(es *event.Store, pID PlayerID) (*Player, error) {
	codec, err := Codec()
	if err != nil {
		return nil, err
	}
	plr := NewPlayer()
	streamID := string(pID)
	for rec := range es.Load(streamID) {
		e, err := codec.Decode(rec)
		if err != nil {
			return nil, err
		}
		plr.Mutate(e)
	}
	if plr.ID == "" {
		return nil, fmt.Errorf("Player not found")
	}
	return plr, nil
}
