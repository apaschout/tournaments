package tournaments

import (
	"fmt"
	"log"
	"time"

	"github.com/cognicraft/event"
	"github.com/cognicraft/uuid"
)

type TrackerCreated struct {
	ID         string    `json:"id"`
	OccurredOn time.Time `json:"occurred-on"`
	Tracker    TrackerID `json:"tracker"`
	Player     PlayerID  `json:"player"`
}

type TrackerMatchPlayed struct {
	ID         string    `json:"id"`
	OccurredOn time.Time `json:"occurred-on"`
	Tracker    TrackerID `json:"tracker"`
}

type TrackerMatchWon struct {
	ID         string    `json:"id"`
	OccurredOn time.Time `json:"occurred-on"`
	Tracker    TrackerID `json:"tracker"`
}

type TrackerGamePlayed struct {
	ID         string    `json:"id"`
	OccurredOn time.Time `json:"occurred-on"`
	Tracker    TrackerID `json:"tracker"`
}

type TrackerGameWon struct {
	ID         string    `json:"id"`
	OccurredOn time.Time `json:"occurred-on"`
	Tracker    TrackerID `json:"tracker"`
}

func NewTracker() *Tracker {
	return &Tracker{
		ChangeRecorder: event.NewChangeRecorder(),
	}
}

func (trk *Tracker) Create(id TrackerID, plr PlayerID) error {
	if trk.ID != "" {
		return fmt.Errorf("Tracker already exists")
	}
	if id == "" {
		return fmt.Errorf("A Tracker's ID may not be empty")
	}
	if plr == "" {
		return fmt.Errorf("No Player specified")
	}
	trk.Apply(TrackerCreated{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Tracker:    id,
		Player:     plr,
	})
	log.Printf("Event: Tracker %v: Created\n", trk.ID)
	return nil
}

func (trk *Tracker) IncrementMatches() error {
	if trk.ID == "" {
		return fmt.Errorf("Tracker does not exist")
	}
	trk.Apply(TrackerMatchPlayed{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Tracker:    trk.ID,
	})
	log.Printf("Event: Tracker %s: Matches for Player %s incremented", trk.ID, trk.Player)
	return nil
}

func (trk *Tracker) IncrementMatchesWon() error {
	if trk.ID == "" {
		return fmt.Errorf("Tracker does not exist")
	}
	trk.Apply(TrackerMatchWon{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Tracker:    trk.ID,
	})
	log.Printf("Event: Tracker %s: MatchWins for Player %s incremented", trk.ID, trk.Player)
	return nil
}

func (trk *Tracker) IncrementGames() error {
	if trk.ID == "" {
		return fmt.Errorf("Tracker does not exist")
	}
	trk.Apply(TrackerGamePlayed{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Tracker:    trk.ID,
	})
	log.Printf("Event: Tracker %s: Games for Player %s incremented", trk.ID, trk.Player)
	return nil
}

func (trk *Tracker) IncrementGamesWon() error {
	if trk.ID == "" {
		return fmt.Errorf("Tracker does not exist")
	}
	trk.Apply(TrackerGameWon{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Tracker:    trk.ID,
	})
	log.Printf("Event: Tracker %s: GameWins for Player %s incremented", trk.ID, trk.Player)
	return nil
}

func (trk *Tracker) Mutate(e event.Event) {
	trk.Version++
	switch e := e.(type) {
	case TrackerCreated:
		trk.ID = e.Tracker
		trk.Player = e.Player
	case TrackerGamePlayed:
		trk.Games++
	case TrackerGameWon:
		trk.GameWins++
	case TrackerMatchPlayed:
		trk.Matches++
	case TrackerMatchWon:
		trk.MatchWins++
	}
}

func (trk *Tracker) Apply(e event.Event) {
	trk.Record(e)
	trk.Mutate(e)
}

func (trk *Tracker) Save(es *event.Store, metadata interface{}) error {
	if len(trk.Changes()) == 0 {
		return nil
	}
	streamID := string(trk.ID)
	exp := trk.Version - uint64(len(trk.Changes()))
	codec, err := Codec()
	if err != nil {
		return err
	}
	recs, err := codec.EncodeAll(trk.Changes(), event.WithMetadata(metadata))
	if err != nil {
		return err
	}
	err = es.Append(streamID, exp, recs)
	if err != nil {
		return err
	}
	trk.ClearChanges()
	return nil
}

func LoadTracker(es *event.Store, tID TrackerID) (*Tracker, error) {
	codec, err := Codec()
	if err != nil {
		return nil, err
	}
	trk := NewTracker()
	streamID := string(tID)
	for rec := range es.Load(streamID) {
		e, err := codec.Decode(rec)
		if err != nil {
			return nil, err
		}
		trk.Mutate(e)
	}
	if trk.ID == "" {
		return nil, fmt.Errorf("Tracker not found")
	}
	return trk, nil
}
