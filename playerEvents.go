package tournaments

import (
	"fmt"
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
	Player     PlayerID  `json:"season"`
	Name       string    `json:"name"`
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
	case PlayerNameChanged:
		plr.Name = e.Name
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
		return nil, fmt.Errorf("Season not found")
	}
	return plr, nil
}
