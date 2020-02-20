package tournaments

import (
	"fmt"
	"log"
	"time"

	"github.com/cognicraft/event"
	"github.com/cognicraft/uuid"
)

type TournamentCreated struct {
	ID         string       `json:"id"`
	OccurredOn time.Time    `json:"occured-on"`
	Tournament TournamentID `json:"tournament"`
}

type TournamentDeleted struct {
	ID         string       `json:"id"`
	OccurredOn time.Time    `json:"occurred-on"`
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

type TournamentPlayerRegistered struct {
	ID         string       `json:"id"`
	OccurredOn time.Time    `json:"occurred-on"`
	Tournament TournamentID `json:"tournament"`
	Player     PlayerID     `json:"player"`
}

type TournamentPlayerDropped struct {
	ID         string       `json:"id"`
	OccurredOn time.Time    `json:"occurred-on"`
	Tournament TournamentID `json:"tournament"`
	Player     PlayerID     `json:"player"`
}

type TournamentPhaseChanged struct {
	ID         string       `json:"id"`
	OccurredOn time.Time    `json:"occurred-on"`
	Tournament TournamentID `json:"tournament"`
	Phase      Phase        `json:"phase"`
}

type TournamentFormatChanged struct {
	ID         string       `json:"id"`
	OccurredOn time.Time    `json:"occurred-on"`
	Tournament TournamentID `json:"tournament"`
	Format     string       `json:"format"`
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
	log.Printf("Event: Tournament %v: Created\n", trn.ID)
	return nil
}

func (trn *Tournament) Delete() error {
	if trn.ID == "" {
		return fmt.Errorf("Tournament does not exist")
	}
	trn.Apply(TournamentDeleted{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Tournament: trn.ID,
	})
	log.Printf("Event: Tournament %v: Deleted\n", trn.ID)
	return nil
}

func (trn *Tournament) ChangeName(name string) error {
	if trn.ID == "" {
		return fmt.Errorf("Tournament does not exist")
	}
	if name == "" {
		return fmt.Errorf("A Tournament's name may not be empty")
	}
	if trn.Phase != PhaseInitialization {
		return fmt.Errorf("Changing Name is not allowed in this Phase")
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
	log.Printf("Event: Tournament %v: Name Changed To %s\n", trn.ID, name)
	return nil
}

func (trn *Tournament) ChangePhase(p Phase) error {
	if trn.ID == "" {
		return fmt.Errorf("Tournament does not exist")
	}
	if p == "" {
		return fmt.Errorf("Phase not specified")
	}
	if trn.Phase == p {
		return nil
	}
	trn.Apply(TournamentPhaseChanged{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Tournament: trn.ID,
		Phase:      p,
	})
	log.Printf("Event: Tournament %v: Phase Changed To %v\n", trn.ID, p)
	return nil
}

func (trn *Tournament) ChangeFormat(f string) error {
	if trn.ID == "" {
		return fmt.Errorf("Tournament does not exist")
	}
	if f == "" {
		return fmt.Errorf("Format not specified")
	}
	if trn.Phase != PhaseInitialization {
		return fmt.Errorf("Changing Format is not allowed in this Phase")
	}
	if trn.Format == f {
		return nil
	}
	trn.Apply(TournamentFormatChanged{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Tournament: trn.ID,
		Format:     f,
	})
	log.Printf("Event: Tournament %v: Format Changed To %s\n", trn.ID, f)
	return nil
}

func (trn *Tournament) RegisterPlayer(pID PlayerID) error {
	if trn.ID == "" {
		return fmt.Errorf("Tournament does not exist")
	}
	if pID == "" {
		return fmt.Errorf("No Player specified")
	}
	if trn.isPlayerRegistered(pID) {
		return fmt.Errorf("Player already registered")
	}
	if trn.Phase != PhaseRegistration {
		return fmt.Errorf("Not in registration phase")
	}
	trn.Apply(TournamentPlayerRegistered{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Tournament: trn.ID,
		Player:     pID,
	})
	log.Printf("Event: Tournament %v: Player %v Registered\n", trn.ID, pID)
	return nil
}

func (trn *Tournament) DropPlayer(pID PlayerID) error {
	if trn.ID == "" {
		return fmt.Errorf("Tournament does not exist")
	}
	if pID == "" {
		return fmt.Errorf("No Player specified")
	}
	if !trn.isPlayerRegistered(pID) {
		return fmt.Errorf("Player is not registered")
	}
	trn.Apply(TournamentPlayerDropped{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Tournament: trn.ID,
		Player:     pID,
	})
	log.Printf("Event: Tournament %v: Player %v Dropped\n", trn.ID, pID)
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
	log.Printf("Event: Tournament %v: Started\n", trn.ID)
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
	log.Printf("Event: Tournament %v: Ended\n", trn.ID)
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
		trn.Name = string(e.Tournament)
		trn.Phase = PhaseInitialization
	case TournamentDeleted:
		trn.Name = ""
		trn.Phase = ""
		trn.Start = ""
		trn.End = ""
		trn.Format = ""
		trn.Seats = nil
		trn.Players = nil
	case TournamentNameChanged:
		trn.Name = e.Name
	case TournamentPhaseChanged:
		if e.Phase == PhaseDraft {
			trn.permutatePlayers(e.OccurredOn)
		}
		trn.Phase = e.Phase
	case TournamentFormatChanged:
		trn.Format = e.Format
	case TournamentPlayerRegistered:
		trn.Players = append(trn.Players, Participant{Player: e.Player})
	case TournamentPlayerDropped:
		trn.removePlayer(e.Player)
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

func LoadTournament(es *event.Store, tID TournamentID) (*Tournament, error) {
	codec, err := Codec()
	if err != nil {
		return nil, err
	}
	trn := NewTournament()
	streamID := string(tID)
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

func (trn *Tournament) isPlayerRegistered(pID PlayerID) bool {
	for _, v := range trn.Players {
		if v.Player == pID {
			return true
		}
	}
	return false
}

func (trn *Tournament) removePlayer(pID PlayerID) {
	for i, v := range trn.Players {
		if v.Player == pID {
			trn.Players = append(trn.Players[:i], trn.Players[i+1:]...)
		}
	}
}
