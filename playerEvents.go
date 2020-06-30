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
	Mail       string    `json:"mail"`
	Password   string    `json:"password"`
	Role       string    `json:"role"`
	Tracker    TrackerID `json:"tracker"`
}

type PlayerNameChanged struct {
	ID         string    `json:"id"`
	OccurredOn time.Time `json:"occurred-on"`
	Player     PlayerID  `json:"player"`
	Name       string    `json:"name"`
}

type PlayerTournamentRegistered struct {
	ID         string       `json:"id"`
	OccurredOn time.Time    `json:"occurred-on"`
	Player     PlayerID     `json:"player"`
	Tournament TournamentID `json:"Tournament"`
}

func NewPlayer(s *Server) *Player {
	return &Player{
		Server:         s,
		ChangeRecorder: event.NewChangeRecorder(),
	}
}

func (plr *Player) Create(id PlayerID, tracker TrackerID, role string, mail string, password string) error {
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
		Mail:       mail,
		Password:   password,
		Role:       role,
		Tracker:    tracker,
	})
	trk := NewTracker()
	err = trk.Create(tracker, id)
	if err != nil {
		return err
	}
	err = trk.Save(plr.Server.es, nil)
	if err != nil {
		return err
	}
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
		plr.Tracker = e.Tracker
	case PlayerNameChanged:
		plr.Name = e.Name
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

func LoadPlayer(s *Server, pID PlayerID) (*Player, error) {
	codec, err := Codec()
	if err != nil {
		return nil, err
	}
	plr := NewPlayer(s)
	streamID := string(pID)
	for rec := range s.es.Load(streamID) {
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

func LoadPlayers(s *Server, pIDs []PlayerID) ([]*Player, error) {
	res := make([]*Player, len(pIDs))
	for i, pID := range pIDs {
		res[i], err = LoadPlayer(s, pID)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
