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

type TournamentMatchesCreated struct {
	ID         string       `json:"id"`
	OccurredOn time.Time    `json:"occurred-on"`
	Tournament TournamentID `json:"tournament"`
}

type TournamentGameEnded struct {
	ID         string       `json:"id"`
	OccurredOn time.Time    `json:"occurred-on"`
	Tournament TournamentID `json:"tournament"`
	Match      int          `json:"match"` //index
	Game       int          `json:"game"`  //index
	Winner     PlayerID     `json:"winner"`
	Draw       bool         `json:"draw"`
}

type TournamentGamesToWinChanged struct {
	ID         string       `json:"id"`
	OccurredOn time.Time    `json:"occurred-on"`
	Tournament TournamentID `json:"tournament"`
	GamesToWin int          `json:"gamesToWin"`
}

func NewTournament(s *Server) *Tournament {
	return &Tournament{
		Server:         s,
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
	if trn.GamesToWin == 0 {
		return fmt.Errorf("Games To Win can't be 0")
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
	if p == PhaseRounds {
		err = trn.CreateMatches()
		if err != nil {
			return err
		}
	}
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

func (trn *Tournament) ChangeGamesToWin(n int) error {
	if trn.ID == "" {
		return fmt.Errorf("Tournament does not exist")
	}
	if n <= 0 {
		return fmt.Errorf("Games To Win has to be at least 1")
	}
	trn.Apply(TournamentGamesToWinChanged{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Tournament: trn.ID,
		GamesToWin: n,
	})
	log.Printf("Event: Tournament %v: GamesToWin changed to %d\n", trn.ID, n)
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

func (trn *Tournament) CreateMatches() error {
	if trn.ID == "" {
		return fmt.Errorf("Tournament does not exist")
	}
	if trn.Matches != nil {
		return fmt.Errorf("Tournament already has matches")
	}
	trn.Apply(TournamentMatchesCreated{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Tournament: trn.ID,
	})
	log.Printf("Event: Tournament %v: Matches created\n", trn.ID)
	return nil
}

func (trn *Tournament) EndGame(match int, game int, wnr PlayerID, draw bool) error {
	if trn.ID == "" {
		return fmt.Errorf("Tournament does not exist")
	}
	if match >= len(trn.Matches) {
		return fmt.Errorf("Match index does not exist")
	}
	if game >= len(trn.Matches[match].Games) {
		return fmt.Errorf("Game index does not exist")
	}
	if trn.Matches[match].Games[game].Ended {
		return fmt.Errorf("Game has already ended")
	}
	if wnr == "" && !draw {
		return nil
	}
	if wnr != "" && draw {
		wnr = ""
	}
	trn.Apply(TournamentGameEnded{
		ID:         uuid.MakeV4(),
		OccurredOn: time.Now().UTC(),
		Tournament: trn.ID,
		Match:      match,
		Game:       game,
		Winner:     wnr,
		Draw:       draw,
	})
	log.Printf("Event: Tournament %v: Match %d: Game %d ended... Winner: %v, Draw: %v", trn.ID, match, game, wnr, draw)
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
		trn.Deleted = true
	case TournamentNameChanged:
		trn.Name = e.Name
	case TournamentGamesToWinChanged:
		trn.GamesToWin = e.GamesToWin
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
	case TournamentMatchesCreated:
		trn.MakeMatches()
	case TournamentGameEnded:
		g := &trn.Matches[e.Match].Games[e.Game]
		g.Winner = e.Winner
		g.Draw = e.Draw
		g.Ended = true
		trn.manageGameWins(e.Match, e.Game)
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

func LoadTournament(s *Server, tID TournamentID) (*Tournament, error) {
	codec, err := Codec()
	if err != nil {
		return nil, err
	}
	trn := NewTournament(s)
	streamID := string(tID)
	for rec := range s.es.Load(streamID) {
		e, err := codec.Decode(rec)
		if err != nil {
			return nil, err
		}
		trn.Mutate(e)
	}
	if trn.Deleted {
		return nil, fmt.Errorf("Tournament was deleted")
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
