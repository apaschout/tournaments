package tournaments

import (
	"fmt"
	"time"

	"github.com/cognicraft/event"
	"github.com/cognicraft/uuid"
)

type TournamentCreated struct {
	ID         string       `json:"id"`
	OccurredOn time.Time    `json:"occured-on"`
	Tournament TournamentID `json:"tournament"`
}

type TournamentNameChanged struct {
	ID         string       `json:"id"`
	OccurredOn time.Time    `json:"occurred-on"`
	Tournament TournamentID `json:"tournament"`
	Name       string       `json:"name"`
}

type TournamentStarted struct {
	ID         string       `json:"id"`
	OccurredOn time.Time    `json:"occurred-on"`
	Tournament TournamentID `json:"tournament"`
	Start      time.Time    `json:"start"`
}

type TournamentEnded struct {
	ID         string       `json:"id"`
	OccurredOn time.Time    `json:"occurred-on"`
	Tournament TournamentID `json:"tournament"`
	End        time.Time    `json:"end"`
}

func NewTournament() *Tournament {
	return &Tournament{
		ChangeRecorder: event.NewChangeRecorder(),
	}
}

func (trn *Tournament) Create(id TournamentID) error {
	if trn.ID != "" {
		return fmt.Errorf("Tournament already exists")
	}
	if id == "" {
		return fmt.Errorf("A Tournament's ID may not be empty")
	}
	trn.Apply(TournamentCreated{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Tournament: id,
	})
	return nil
}

func (trn *Tournament) ChangeName(name string) error {
	if trn.ID == "" {
		return fmt.Errorf("Tournament does not exist")
	}
	if name == "" {
		return fmt.Errorf("A Tournament's name may not be empty")
	}
	if trn.Name == name {
		return nil
	}
	trn.Apply(TournamentNameChanged{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Tournament: trn.ID,
		Name:       name,
	})
	return nil
}

func (trn *Tournament) Begin() error {
	if trn.ID == "" {
		return fmt.Errorf("Tournament does not exist")
	}
	if trn.Start != "" {
		return fmt.Errorf("Tournament already started")
	}
	trn.Apply(TournamentStarted{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Tournament: trn.ID,
		Start:      time.Now().UTC(),
	})
	return nil
}

func (trn *Tournament) Finish() error {
	if trn.ID == "" {
		return fmt.Errorf("Tournament does not exist")
	}
	if trn.Start == "" {
		return fmt.Errorf("Tournament has not started yet")
	}
	if trn.End != "" {
		return fmt.Errorf("Tournament has already ended")
	}
	trn.Apply(TournamentEnded{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Tournament: trn.ID,
		End:        time.Now().UTC(),
	})
	return nil
}

func (trn *Tournament) Apply(e event.Event) {
	trn.Record(e)
	trn.Mutate(e)
}

func (trn *Tournament) Mutate(e event.Event) {
	trn.Version++
	switch e := e.(type) {
	case TournamentCreated:
		trn.ID = e.Tournament
	case TournamentNameChanged:
		trn.Name = e.Name
	case TournamentStarted:
		trn.Start = e.Start.String()
	case TournamentEnded:
		trn.End = e.End.String()
	}
}

func (trn *Tournament) Save(es *event.Store, metadata interface{}) error {
	if len(trn.Changes()) == 0 {
		return nil
	}
	streamID := string(trn.ID)
	exp := trn.Version - uint64(len(trn.Changes()))
	codec, err := Codec()
	if err != nil {
		return err
	}
	recs, err := codec.EncodeAll(trn.Changes(), event.WithMetadata(metadata))
	if err != nil {
		return err
	}
	err = es.Append(streamID, exp, recs)
	if err != nil {
		return err
	}
	trn.ClearChanges()
	return nil
}

func LoadTournament(es *event.Store, sID TournamentID) (*Tournament, error) {
	codec, err := Codec()
	if err != nil {
		return nil, err
	}
	trn := NewTournament()
	streamID := string(sID)
	for rec := range es.Load(streamID) {
		e, err := codec.Decode(rec)
		if err != nil {
			return nil, err
		}
		trn.Mutate(e)
	}
	if trn.ID == "" {
		return nil, fmt.Errorf("Tournament not found")
	}
	return trn, nil
}
