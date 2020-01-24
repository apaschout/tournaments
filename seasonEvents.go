package tournaments

import (
	"fmt"
	"time"

	"github.com/cognicraft/event"
	"github.com/cognicraft/uuid"
)

type SeasonCreated struct {
	ID         string    `json:"id"`
	OccurredOn time.Time `json:"occured-on"`
	Season     SeasonID  `json:"season"`
}

type SeasonNameChanged struct {
	ID         string    `json:"id"`
	OccurredOn time.Time `json:"occurred-on"`
	Season     SeasonID  `json:"season"`
	Name       string    `json:"name"`
}

type SeasonStarted struct {
	ID         string    `json:"id"`
	OccurredOn time.Time `json:"occurred-on"`
	Season     SeasonID  `json:"season"`
	Start      time.Time `json:"start"`
}

type SeasonEnded struct {
	ID         string    `json:"id"`
	OccurredOn time.Time `json:"occurred-on"`
	Season     SeasonID  `json:"season"`
	End        time.Time `json:"end"`
}

func NewSeason() *Season {
	return &Season{
		ChangeRecorder: event.NewChangeRecorder(),
	}
}

func (seas *Season) Create(id SeasonID) error {
	if seas.ID != "" {
		return fmt.Errorf("Season already exists")
	}
	if id == "" {
		return fmt.Errorf("A Season's ID may not be empty")
	}
	seas.Apply(SeasonCreated{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Season:     id,
	})
	return nil
}

func (seas *Season) ChangeName(name string) error {
	if seas.ID == "" {
		return fmt.Errorf("Season does not exist")
	}
	if name == "" {
		return fmt.Errorf("A Season's name may not be empty")
	}
	if seas.Name == name {
		return nil
	}
	seas.Apply(SeasonNameChanged{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Season:     seas.ID,
		Name:       name,
	})
	return nil
}

func (seas *Season) Begin() error {
	if seas.ID == "" {
		return fmt.Errorf("Season does not exist")
	}
	if seas.Start != "" {
		return fmt.Errorf("Season already started")
	}
	seas.Apply(SeasonStarted{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Season:     seas.ID,
		Start:      time.Now().UTC(),
	})
	return nil
}

func (seas *Season) Finish() error {
	if seas.ID == "" {
		return fmt.Errorf("Season does not exist")
	}
	if seas.Start == "" {
		return fmt.Errorf("Season has not started yet")
	}
	if seas.End != "" {
		return fmt.Errorf("Season has already ended")
	}
	seas.Apply(SeasonEnded{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Season:     seas.ID,
		End:        time.Now().UTC(),
	})
	return nil
}

func (seas *Season) Apply(e event.Event) {
	seas.Record(e)
	seas.Mutate(e)
}

func (seas *Season) Mutate(e event.Event) {
	seas.Version++
	switch e := e.(type) {
	case SeasonCreated:
		seas.ID = e.Season
	case SeasonNameChanged:
		seas.Name = e.Name
	case SeasonStarted:
		seas.Start = e.Start.String()
	case SeasonEnded:
		seas.End = e.End.String()
	}
}

func (seas *Season) Save(es *event.Store, metadata interface{}) error {
	if len(seas.Changes()) == 0 {
		return nil
	}
	streamID := string(seas.ID)
	exp := seas.Version - uint64(len(seas.Changes()))
	codec, err := Codec()
	if err != nil {
		return err
	}
	recs, err := codec.EncodeAll(seas.Changes(), event.WithMetadata(metadata))
	if err != nil {
		return err
	}
	err = es.Append(streamID, exp, recs)
	if err != nil {
		return err
	}
	seas.ClearChanges()
	return nil
}

func LoadSeason(es *event.Store, sID SeasonID) (*Season, error) {
	codec, err := Codec()
	if err != nil {
		return nil, err
	}
	seas := NewSeason()
	streamID := string(sID)
	for rec := range es.Load(streamID) {
		e, err := codec.Decode(rec)
		if err != nil {
			return nil, err
		}
		seas.Mutate(e)
	}
	if seas.ID == "" {
		return nil, fmt.Errorf("Season not found")
	}
	return seas, nil
}
