package tournaments

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/cognicraft/event"
	"github.com/cognicraft/hyper"
	"github.com/cognicraft/uuid"
)

type Tournament struct {
	ID           TournamentID  `json:"id"`
	Version      uint64        `json:"version"`
	Name         string        `json:"name"`
	Phase        Phase         `json:"phase"`
	Start        string        `json:"start,omitempty"`
	End          string        `json:"end,omitempty"`
	Format       string        `json:"format,omitempty"`
	Seats        []Seat        `json:"seats"`
	Matches      []Match       `json:"matches"`
	GamesToWin   int           `json:"gamesToWin"`
	Participants []Participant `json:"players,omitempty"`
	Deleted      bool          `json:"deleted"`
	*event.ChangeRecorder
	Server *Server
}

type Participant struct {
	Player    PlayerID `json:"player"`
	SeatIndex int      `json:"seatIndex"`
	Deck      DeckID   `json:"deck"`
	Matches   int      `json:"matches"`
	Games     int      `json:"games"`
	MatchWins int      `json:"matchWins"`
	GameWins  int      `json:"gameWins"`
}

type Seat struct {
	Index  int      `json:"index"`
	Player PlayerID `json:"player"`
}

type TournamentID string

const (
	ActionDelete           = "delete"
	ActionRegisterPlayer   = "register-player"
	ActionDropPlayer       = "drop-player"
	ActionCreate           = "create"
	ActionChangeFormat     = "change-format"
	ActionEndPhase         = "end-phase"
	ActionEndGame          = "end-game"
	ActionChangeGamesToWin = "change-gamestowin"
)

const (
	ArgumentTournamentID = "tid"
	ArgumentPlayerID     = "pid"
	ArgumentName         = "name"
	ArgumentRole         = "role"
	ArgumentFormat       = "format"
	ArgumentMatch        = "match"
	ArgumentGame         = "game"
	ArgumentGamesToWin   = "gamestowin"
	ArgumentDraw         = "draw"
)

func (s *Server) handleGETTournaments(w http.ResponseWriter, r *http.Request) {
	isHtmlReq := strings.Contains(r.Header.Get("Accept"), "text/html")
	aID, err := s.authenticate(r)
	fmt.Println(aID)
	if err != nil {
		http.Redirect(w, r, "/api/signin", http.StatusSeeOther)
		return
	}
	resolve := hyper.ExternalURLResolver(r)
	res := hyper.Item{
		Label: "Tournaments",
		Type:  "tournaments",
	}
	link := hyper.Link{
		Rel:  "self",
		Href: resolve(".").String(),
	}
	s.tournaments, err = s.p.FindAllTournaments()
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, isHtmlReq)
		return
	}
	for _, trn := range s.tournaments {
		item := trn.MakeUndetailedHyperItem(resolve)
		res.AddItem(item)
	}
	res.AddLink(link)
	if isHtmlReq {
		err = templ.ExecuteTemplate(w, "tournaments.html", res)
		if err != nil {
			log.Println(err)
		}
	} else {
		hyper.Write(w, 200, res)
	}
}

func (s *Server) handleGETTournament(w http.ResponseWriter, r *http.Request) {
	isHtmlReq := strings.Contains(r.Header.Get("Accept"), "text/html")
	aID, err := s.authenticate(r)
	fmt.Println(aID)
	if err != nil {
		http.Redirect(w, r, "/api/signin", http.StatusSeeOther)
		return
	}
	resolve := hyper.ExternalURLResolver(r)
	tID := TournamentID(r.Context().Value(":id").(string))

	trn, err := LoadTournament(s, tID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, isHtmlReq)
		return
	}

	res := trn.MakeDetailedHyperItem(resolve)
	plrs, err := trn.MakeParticipantsHyperItem(resolve)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, isHtmlReq)
		return
	}
	res.AddItem(plrs)
	if isHtmlReq {
		err = switchPhase(trn.Phase, w, r, res)
		if err != nil {
			log.Println(err)
		}
	} else {
		hyper.Write(w, http.StatusOK, res)
	}
}

func (s *Server) handlePOSTTournament(w http.ResponseWriter, r *http.Request) {
	isHtmlReq := strings.Contains(r.Header.Get("Accept"), "text/html")
	accID, err := s.authenticate(r)
	if err != nil {
		http.Redirect(w, r, "/api/signin", http.StatusSeeOther)
		return
	}
	err = s.checkOrganizerPermissions(accID, "Unable to edit Tournament: Insufficient Permissions")
	if err != nil {
		handleError(w, http.StatusForbidden, err, isHtmlReq)
		return
	}
	cmd := hyper.ExtractCommand(r)
	tID := TournamentID(r.Context().Value(":id").(string))

	trn, err := LoadTournament(s, tID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, isHtmlReq)
		return
	}
	switch cmd.Action {
	case ActionRegisterPlayer:
		pID := PlayerID(cmd.Arguments.String(ArgumentPlayerID))
		err = trn.RegisterPlayer(pID)
	case ActionDropPlayer:
		pID := PlayerID(cmd.Arguments.String(ArgumentPlayerID))
		err = trn.DropPlayer(pID)
	case ActionChangeName:
		var ok bool
		newName := cmd.Arguments.String(ArgumentName)
		ok, err = s.p.IsTournamentNameAvailable(newName)
		if !ok {
			handleError(w, http.StatusInternalServerError, err, isHtmlReq)
			return
		}
		err = trn.ChangeName(newName)
	case ActionChangeGamesToWin:
		n := cmd.Arguments.Int(ArgumentGamesToWin)
		err = trn.ChangeGamesToWin(n)
	case ActionChangeFormat:
		f := cmd.Arguments.String(ArgumentFormat)
		err = trn.ChangeFormat(f)
	case ActionEndPhase:
		switch trn.Format {
		case "cube":
			err = trn.handleEndPhaseCube()
		case "":
			err = fmt.Errorf("Can't proceed to next Phase: Format not set")
		default:
			err = fmt.Errorf("Format not recognized: %s", trn.Format)
		}
	case ActionDelete:
		err = trn.Delete()
	case ActionEndGame:
		m := cmd.Arguments.Int(ArgumentMatch)
		g := cmd.Arguments.Int(ArgumentGame)
		wnr := cmd.Arguments.String(ArgumentPlayerID)
		draw := cmd.Arguments.Bool(ArgumentDraw)
		err = trn.EndGame(m, g, PlayerID(wnr), draw)
		if err != nil {
			handleError(w, http.StatusInternalServerError, err, isHtmlReq)
			return
		}
	default:
		err = fmt.Errorf("Action not recognized: %s", cmd.Action)
	}
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, isHtmlReq)
		return
	}
	err = trn.Save(s.es, nil)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, isHtmlReq)
		return
	}
	if isHtmlReq {
		http.Redirect(w, r, r.URL.String(), http.StatusSeeOther)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) handlePOSTTournaments(w http.ResponseWriter, r *http.Request) {
	isHtmlReq := strings.Contains(r.Header.Get("Accept"), "text/html")
	accID, err := s.authenticate(r)
	if err != nil {
		http.Redirect(w, r, "/api/signin", http.StatusSeeOther)
		return
	}
	err = s.checkOrganizerPermissions(accID, "Unable to create Tournament: Insufficient Permissions")
	if err != nil {
		handleError(w, http.StatusForbidden, err, isHtmlReq)
		return
	}
	cmd := hyper.ExtractCommand(r)
	switch cmd.Action {
	case ActionCreate:
		trn := NewTournament(s)
		err = trn.Create(TournamentID(uuid.MakeV4()))
		if err != nil {
			handleError(w, http.StatusInternalServerError, err, isHtmlReq)
			return
		}
		err = trn.Save(s.es, nil)
	default:
		err = fmt.Errorf("Action not recognized: %s", cmd.Action)
	}
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, isHtmlReq)
		return
	}
	if isHtmlReq {
		http.Redirect(w, r, r.URL.String(), http.StatusSeeOther)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (trn *Tournament) handleEndPhaseCube() error {
	switch trn.Phase {
	case PhaseInitialization:
		if trn.Name == "" {
			return fmt.Errorf("Can't proceed to next Phase: Name not set")
		}
		err = trn.ChangePhase(PhaseRegistration)
	case PhaseRegistration:
		if len(trn.Participants) == 0 {
			return fmt.Errorf("Can't proceed to next Phase: No Players registered")
		}
		err = trn.ChangePhase(PhaseDraft)
		if err != nil {
			return err
		}
		err = trn.Begin()
	case PhaseDraft:
		err = trn.ChangePhase(PhaseRounds)
	case PhaseRounds:
		for _, mtc := range trn.Matches {
			if !mtc.Ended {
				return fmt.Errorf("Not all Matches have ended")
			}
		}
		err = trn.ChangePhase(PhaseEnded)
		if err != nil {
			return err
		}
		err = trn.Finish()
	case PhaseEnded:
		err = fmt.Errorf("Tournament has already ended")
	default:
		err = fmt.Errorf("Phase not recognized")
	}
	if err != nil {
		return err
	}
	return nil
}

func switchPhase(p Phase, w http.ResponseWriter, r *http.Request, data interface{}) error {
	switch p {
	case PhaseInitialization:
		err = templ.ExecuteTemplate(w, "tournamentInitialization.html", data)
	case PhaseRegistration:
		err = templ.ExecuteTemplate(w, "tournamentRegistration.html", data)
	case PhaseDraft:
		err = templ.ExecuteTemplate(w, "tournamentDraft.html", data)
	case PhaseRounds:
		err = templ.ExecuteTemplate(w, "tournamentRoundRobin.html", data)
	case PhaseEnded:
		err = templ.ExecuteTemplate(w, "tournamentEnded.html", data)
	default:
		http.Redirect(w, r, "/api/tournaments/", http.StatusSeeOther)
	}
	if err != nil {
		return err
	}
	return nil
}

func (trn *Tournament) permutatePlayers(eventTime time.Time) {
	r := rand.New(rand.NewSource(eventTime.Unix()))
	perm := r.Perm(len(trn.Participants))
	for i, randIndex := range perm {
		trn.Participants[i].SeatIndex = randIndex
	}
}

func (trn *Tournament) getParticipantByID(ID PlayerID) *Participant {
	for i := 0; i < len(trn.Participants); i++ {
		if trn.Participants[i].Player == ID {
			return &trn.Participants[i]
		}
	}
	return nil
}

func (trn *Tournament) MakeParticipantsHyperItem(resolve hyper.ResolverFunc) (hyper.Item, error) {
	plrs := hyper.Item{
		Label: "Participating Players",
		Type:  "participants",
	}
	for i, par := range trn.Participants {
		plr, err := trn.Server.p.FindPlayerByID(par.Player)
		if err != nil {
			return plrs, err
		}
		item := hyper.Item{
			Label: plr.Name,
			Type:  "participant",
			ID:    string(par.Player),
			Properties: []hyper.Property{
				{
					Label: "Name",
					Name:  "name",
					Value: plr.Name,
				},
				{
					Label: "Seat Index",
					Name:  "seatIndex",
					Value: trn.Participants[i].SeatIndex,
				},
				{
					Label: "Deck",
					Name:  "deck",
					Value: trn.Participants[i].Deck,
				},
				{
					Label: "Matches",
					Name:  "matches",
					Value: trn.Participants[i].Matches,
				},
				{
					Label: "Match Wins",
					Name:  "matchWins",
					Value: trn.Participants[i].MatchWins,
				},
				{
					Label: "Games",
					Name:  "games",
					Value: trn.Participants[i].Games,
				},
				{
					Label: "Game Wins",
					Name:  "gameWins",
					Value: trn.Participants[i].GameWins,
				},
			},
		}
		pLink := hyper.Link{
			Rel:  "details",
			Href: resolve("../players/%s", par.Player).String(),
		}
		item.AddLink(pLink)
		plrs.AddItem(item)
	}
	return plrs, nil
}

func (trn *Tournament) MakeUndetailedHyperItem(resolve hyper.ResolverFunc) hyper.Item {
	item := hyper.Item{
		Label: trn.Name,
		Type:  "tournament",
		ID:    string(trn.ID),
		Properties: []hyper.Property{
			{
				Label: "Name",
				Name:  "name",
				Value: trn.Name,
			},
		},
	}
	actions := hyper.Actions{
		{
			Label:  "Delete",
			Rel:    ActionDelete,
			Href:   resolve("./%s", trn.ID).String(),
			Method: "POST",
		},
	}
	tLink := hyper.Link{
		Rel:  "details",
		Href: resolve("./%s", trn.ID).String(),
	}
	item.AddActions(actions)
	item.AddLink(tLink)
	return item
}

// func MakeDetailedTrnHyperItem(trn Tournament, resolve hyper.ResolverFunc) hyper.Item {
// 	item := hyper.Item{
// 		Label: trn.Name,
// 		Type:  "tournament",
// 		ID:    string(trn.ID),
// 		Properties: []hyper.Property{
// 			{
// 				Label: "Name",
// 				Name:  "name",
// 				Value: trn.Name,
// 			},
// 			{
// 				Label: "Version",
// 				Name:  "version",
// 				Value: trn.Version,
// 			},
// 			{
// 				Label: "Phase",
// 				Name:  "phase",
// 				Value: trn.Phase,
// 			},
// 			{
// 				Label: "Matches",
// 				Name:  "matches",
// 				Value: trn.Matches,
// 			},
// 			{
// 				Label: "Start",
// 				Name:  "start",
// 				Value: trn.Start,
// 			},
// 			{
// 				Label: "End",
// 				Name:  "end",
// 				Value: trn.End,
// 			},
// 			{
// 				Label: "Format",
// 				Name:  "format",
// 				Value: trn.Format,
// 			},
// 			{
// 				Label: "Games To Win",
// 				Name:  "gamesToWin",
// 				Value: trn.GamesToWin,
// 			},
// 		},
// 	}
// 	link := hyper.Link{
// 		Rel:  "self",
// 		Href: resolve("./%s", trn.ID).String(),
// 	}
// 	actions := hyper.Actions{
// 		{
// 			Label:  "Change Name",
// 			Rel:    ActionChangeName,
// 			Href:   resolve("./%s", trn.ID).String(),
// 			Method: "POST",
// 			Parameters: hyper.Parameters{
// 				{
// 					Name:        ArgumentName,
// 					Placeholder: "New Name...",
// 				},
// 			},
// 		},
// 		{
// 			Label:  "Change Games To Win",
// 			Rel:    ActionChangeGamesToWin,
// 			Href:   resolve("./%s", trn.ID).String(),
// 			Method: "POST",
// 			Parameters: hyper.Parameters{
// 				{
// 					Name:        ArgumentGamesToWin,
// 					Placeholder: "Games To Win...",
// 				},
// 			},
// 		},
// 		{
// 			Label:  "Change Format",
// 			Rel:    ActionChangeFormat,
// 			Href:   resolve("./%s", trn.ID).String(),
// 			Method: "POST",
// 			Parameters: hyper.Parameters{
// 				{
// 					Name:        ArgumentFormat,
// 					Placeholder: "New Format...",
// 				},
// 			},
// 		},
// 		{
// 			Label:  "End Phase",
// 			Rel:    ActionEndPhase,
// 			Href:   resolve("./%s", trn.ID).String(),
// 			Method: "POST",
// 		},
// 		{
// 			Label:  "Delete",
// 			Rel:    ActionDelete,
// 			Href:   resolve("./%s", trn.ID).String(),
// 			Method: "POST",
// 		},
// 		{
// 			Label:  "Register Player",
// 			Rel:    ActionRegisterPlayer,
// 			Href:   resolve("./%s", trn.ID).String(),
// 			Method: "POST",
// 			Parameters: hyper.Parameters{
// 				{
// 					Name: ArgumentPlayerID,
// 				},
// 			},
// 		},
// 		{
// 			Label:  "Drop Player",
// 			Rel:    ActionDropPlayer,
// 			Href:   resolve("./%s", trn.ID).String(),
// 			Method: "POST",
// 			Parameters: hyper.Parameters{
// 				{
// 					Name: ArgumentPlayerID,
// 				},
// 			},
// 		},
// 		{
// 			Label:  "End Game",
// 			Rel:    ActionEndGame,
// 			Href:   resolve("./%s", trn.ID).String(),
// 			Method: "POST",
// 			Parameters: hyper.Parameters{
// 				{
// 					Name: ArgumentMatch,
// 				},
// 				{
// 					Name: ArgumentGame,
// 				},
// 				{
// 					Name: ArgumentPlayerID,
// 				},
// 				{
// 					Name: ArgumentDraw,
// 				},
// 			},
// 		},
// 	}
// 	item.AddActions(actions)
// 	item.AddLink(link)
// 	return item
// }

func (trn *Tournament) MakeDetailedHyperItem(resolve hyper.ResolverFunc) hyper.Item {
	res := hyper.Item{
		Label: trn.Name,
		Type:  "tournament",
		ID:    string(trn.ID),
	}
	selfLink := hyper.Link{
		Rel:  hyper.RelSelf,
		Href: resolve("./%s", trn.ID).String(),
	}
	//Properties
	nameProp := hyper.Property{
		Label: "Name",
		Name:  "name",
		Value: trn.Name,
	}
	phaseProp := hyper.Property{
		Label: "Phase",
		Name:  "phase",
		Value: trn.Phase,
	}
	matchesProp := hyper.Property{
		Label: "Matches",
		Name:  "matches",
		Value: trn.Matches,
	}
	startProp := hyper.Property{
		Label: "Start",
		Name:  "start",
		Value: trn.Start,
	}
	endProp := hyper.Property{
		Label: "End",
		Name:  "end",
		Value: trn.End,
	}
	formatProp := hyper.Property{
		Label: "Format",
		Name:  "format",
		Value: trn.Format,
	}
	g2wProp := hyper.Property{
		Label: "Games To Win",
		Name:  "gamesToWin",
		Value: trn.GamesToWin,
	}
	//Actions
	nameAct := hyper.Action{
		Label:  "Change Name",
		Rel:    ActionChangeName,
		Href:   resolve("./%s", trn.ID).String(),
		Method: "POST",
		Parameters: hyper.Parameters{
			{
				Name:        ArgumentName,
				Placeholder: "New Name...",
			},
		},
	}
	g2wAct := hyper.Action{
		Label:  "Change Games To Win",
		Rel:    ActionChangeGamesToWin,
		Href:   resolve("./%s", trn.ID).String(),
		Method: "POST",
		Parameters: hyper.Parameters{
			{
				Name:        ArgumentGamesToWin,
				Placeholder: "Games To Win...",
			},
		},
	}
	formatAct := hyper.Action{
		Label:  "Change Format",
		Rel:    ActionChangeFormat,
		Href:   resolve("./%s", trn.ID).String(),
		Method: "POST",
		Parameters: hyper.Parameters{
			{
				Name:        ArgumentFormat,
				Placeholder: "New Format...",
			},
		},
	}
	phaseAct := hyper.Action{
		Label:  "End Phase",
		Rel:    ActionEndPhase,
		Href:   resolve("./%s", trn.ID).String(),
		Method: "POST",
	}
	deleteAct := hyper.Action{
		Label:  "Delete",
		Rel:    ActionDelete,
		Href:   resolve("./%s", trn.ID).String(),
		Method: "POST",
	}
	registerAct := hyper.Action{
		Label:  "Register Player",
		Rel:    ActionRegisterPlayer,
		Href:   resolve("./%s", trn.ID).String(),
		Method: "POST",
		Parameters: hyper.Parameters{
			{
				Name: ArgumentPlayerID,
			},
		},
	}
	dropAct := hyper.Action{
		Label:  "Drop Player",
		Rel:    ActionDropPlayer,
		Href:   resolve("./%s", trn.ID).String(),
		Method: "POST",
		Parameters: hyper.Parameters{
			{
				Name: ArgumentPlayerID,
			},
		},
	}
	endGameAct := hyper.Action{
		Label:  "End Game",
		Rel:    ActionEndGame,
		Href:   resolve("./%s", trn.ID).String(),
		Method: "POST",
		Parameters: hyper.Parameters{
			{
				Name: ArgumentMatch,
			},
			{
				Name: ArgumentGame,
			},
			{
				Name: ArgumentPlayerID,
			},
			{
				Name: ArgumentDraw,
			},
		},
	}

	res.AddLink(selfLink)
	res.AddProperty(nameProp)
	res.AddProperty(phaseProp)
	res.AddAction(deleteAct)
	switch trn.Phase {
	case PhaseInitialization:
		res.AddProperty(formatProp)
		res.AddProperty(g2wProp)

		res.AddAction(formatAct)
		res.AddAction(nameAct)
		res.AddAction(g2wAct)
		res.AddAction(phaseAct)
	case PhaseRegistration:
		res.AddProperty(formatProp)
		res.AddAction(registerAct)
		res.AddAction(dropAct)
		res.AddAction(phaseAct)
	case PhaseDraft:
		res.AddProperty(formatProp)
		res.AddProperty(startProp)
		res.AddAction(phaseAct)
	case PhaseRounds:
		res.AddProperty(matchesProp)

		res.AddAction(endGameAct)
		res.AddAction(phaseAct)
	case PhaseEnded:
		res.AddProperty(formatProp)
		res.AddProperty(g2wProp)
		res.AddProperty(matchesProp)
		res.AddProperty(startProp)
		res.AddProperty(endProp)
	}

	return res
}
