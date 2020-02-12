package tournaments

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/cognicraft/event"
	"github.com/cognicraft/hyper"
	"github.com/cognicraft/uuid"
)

type Tournament struct {
	ID      TournamentID `json:"id"`
	Version uint64       `json:"version"`
	Name    string       `json:"name"`
	Phase   Phase        `json:"phase"`
	Start   string       `json:"start,omitempty"`
	End     string       `json:"end,omitempty"`
	Format  string       `json:"format,omitempty"`
	Players []PlayerID   `json:"players,omitempty"`
	*event.ChangeRecorder
}

type TournamentID string

const (
	ActionDelete         = "delete"
	ActionRegisterPlayer = "register-player"
	ActionDropPlayer     = "drop-player"
	ActionChangeName     = "change-name"
	ActionCreate         = "create"
	ActionChangeFormat   = "change-format"
	ActionEndPhase       = "end-phase"
)

const (
	ArgumentTournamentID = "tid"
	ArgumentPlayerID     = "pid"
	ArgumentName         = "name"
	ArgumentFormat       = "format"
)

func (s *Server) handleGETTournaments(w http.ResponseWriter, r *http.Request) {
	isHtmlReq := strings.Contains(r.Header.Get("Accept"), "text/html")
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
	resolve := hyper.ExternalURLResolver(r)
	tID := TournamentID(r.Context().Value(":id").(string))

	trn, err := LoadTournament(s.es, tID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, isHtmlReq)
		return
	}

	res := trn.MakeDetailedHyperItem(resolve)
	plrs, err := s.MakeTournamentPlayersHyperItem(*trn, resolve)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, isHtmlReq)
		return
	}
	res.AddItem(plrs)
	if isHtmlReq {
		// err = templ.ExecuteTemplate(w, "tournament.html", res)
		err = trn.switchTemplates(w, res)
		if err != nil {
			log.Println(err)
		}
	} else {
		hyper.Write(w, http.StatusOK, res)
	}
}

func (s *Server) handlePOSTTournament(w http.ResponseWriter, r *http.Request) {
	isHtmlReq := strings.Contains(r.Header.Get("Accept"), "text/html")
	cmd := hyper.ExtractCommand(r)
	tID := TournamentID(r.Context().Value(":id").(string))

	trn, err := LoadTournament(s.es, tID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, isHtmlReq)
		return
	}
	switch cmd.Action {
	case ActionRegisterPlayer:
		var ok bool
		pID := PlayerID(cmd.Arguments.String(ArgumentPlayerID))
		ok, err = s.p.PlayerExists(pID)
		if !ok {
			handleError(w, http.StatusInternalServerError, err, isHtmlReq)
			return
		}
		err = trn.RegisterPlayer(pID)
	case ActionDropPlayer:
		var ok bool
		pID := PlayerID(cmd.Arguments.String(ArgumentPlayerID))
		ok, err = s.p.PlayerExists(pID)
		if !ok {
			handleError(w, http.StatusInternalServerError, err, isHtmlReq)
			return
		}
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
	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		http.Redirect(w, r, r.URL.String(), http.StatusSeeOther)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) handlePOSTTournaments(w http.ResponseWriter, r *http.Request) {
	isHtmlReq := strings.Contains(r.Header.Get("Accept"), "text/html")
	cmd := hyper.ExtractCommand(r)
	switch cmd.Action {
	case ActionCreate:
		trn := NewTournament()
		err = trn.Create(TournamentID(uuid.MakeV4()))
		if err != nil {
			handleError(w, http.StatusInternalServerError, err, isHtmlReq)
			return
		}
		err = trn.ChangeName(string(trn.ID))
		if err != nil {
			handleError(w, http.StatusInternalServerError, err, isHtmlReq)
			return
		}
		err = trn.Save(s.es, nil)
	default:
		err = fmt.Errorf("Action not recognized")
	}
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, isHtmlReq)
		return
	}
	if strings.Contains(r.Header.Get("Accept"), "text/html") {
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
		if len(trn.Players) == 0 {
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

func (trn *Tournament) switchTemplates(w http.ResponseWriter, item hyper.Item) error {
	switch trn.Phase {
	case PhaseInitialization:
		err = templ.ExecuteTemplate(w, "tournamentInitialization.html", item)
	case PhaseRegistration:
		err = templ.ExecuteTemplate(w, "tournamentRegistration.html", item)
	case PhaseDraft:
		err = templ.ExecuteTemplate(w, "tournamentDraft.html", item)
	case PhaseRounds:
		err = templ.ExecuteTemplate(w, "tournamentRounds.html", item)
	case PhaseEnded:
		err = templ.ExecuteTemplate(w, "tournamentEnded.html", item)
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) MakeTournamentPlayersHyperItem(trn Tournament, resolve hyper.ResolverFunc) (hyper.Item, error) {
	plrs := hyper.Item{
		Label: "Participating Players",
		Type:  "players",
	}
	for _, pID := range trn.Players {
		plr, err := s.p.FindPlayerByID(pID)
		if err != nil {
			return plrs, err
		}
		item := hyper.Item{
			Label: "Player",
			Type:  "player",
			ID:    string(pID),
			Properties: []hyper.Property{
				{
					Label: "Name",
					Name:  "name",
					Value: plr.Name,
				},
			},
		}
		pLink := hyper.Link{
			Rel:  "details",
			Href: resolve("../players/%s", pID).String(),
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
	tLink := hyper.Link{
		Rel:  "details",
		Href: resolve("./%s", trn.ID).String(),
	}
	item.AddLink(tLink)
	return item
}

func (trn *Tournament) MakeDetailedHyperItem(resolve hyper.ResolverFunc) hyper.Item {
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
			{
				Label: "Version",
				Name:  "version",
				Value: trn.Version,
			},
			{
				Label: "Phase",
				Name:  "phase",
				Value: trn.Phase,
			},
			{
				Label: "Start",
				Name:  "start",
				Value: trn.Start,
			},
			{
				Label: "End",
				Name:  "end",
				Value: trn.End,
			},
			{
				Label: "Format",
				Name:  "format",
				Value: trn.Format,
			},
		},
	}
	link := hyper.Link{
		Rel:  "self",
		Href: resolve("./%s", trn.ID).String(),
	}
	actions := hyper.Actions{
		{
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
		},
		{
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
		},
		{
			Label:  "End Phase",
			Rel:    ActionEndPhase,
			Href:   resolve("./%s", trn.ID).String(),
			Method: "POST",
		},
		{
			Label:  "Delete",
			Rel:    ActionDelete,
			Href:   resolve("./%s", trn.ID).String(),
			Method: "POST",
		},
	}
	item.AddActions(actions)
	item.AddLink(link)
	return item
}
