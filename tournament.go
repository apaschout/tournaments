package tournaments

import (
	"fmt"
	"net/http"

	"github.com/cognicraft/event"
	"github.com/cognicraft/hyper"
	"github.com/cognicraft/uuid"
)

type Tournament struct {
	ID       TournamentID `json:"id"`
	Version  uint64       `json:"version"`
	Name     string       `json:"name"`
	Start    string       `json:"start,omitempty"`
	End      string       `json:"end,omitempty"`
	Ongoing  bool         `json:"ongoing,omitempty"`
	Finished bool         `json:"finished,omitempty"`
	Format   string       `json:"format,omitempty"`
	Players  []PlayerID   `json:"players,omitempty"`
	*event.ChangeRecorder
}

type TournamentID string

const (
	ActionStart      = "start"
	ActionEnd        = "end"
	ActionAddPlr     = "addPlayer"
	ActionDelPlr     = "delPlayer"
	ActionChangeName = "changeName"
	ActionCreate     = "create"
)

const (
	ArgumentPlayerID = "pID"
	ArgumentName     = "name"
)

func (s *Server) handleGETTournaments(w http.ResponseWriter, r *http.Request) {
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
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	for _, trn := range s.tournaments {
		item := trn.MakeUndetailedHyperItem(resolve)
		res.AddItem(item)
	}
	res.AddLink(link)
	hyper.Write(w, 200, res)
}

func (s *Server) handleGETTournament(w http.ResponseWriter, r *http.Request) {
	resolve := hyper.ExternalURLResolver(r)
	tID := TournamentID(r.Context().Value(":id").(string))

	trn, err := LoadTournament(s.es, tID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	res := trn.MakeDetailedHyperItem(resolve)
	plrs, err := s.MakeTournamentPlayersHyperItem(*trn, resolve)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	res.AddItem(plrs)
	hyper.Write(w, http.StatusOK, res)
}

func (s *Server) handlePOSTTournament(w http.ResponseWriter, r *http.Request) {
	cmd := hyper.ExtractCommand(r)
	tID := TournamentID(r.Context().Value(":id").(string))

	trn, err := LoadTournament(s.es, tID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	switch cmd.Action {
	case ActionStart:
		err = trn.Begin()
		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}
	case ActionEnd:
		err = trn.Finish()
		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}
	case ActionAddPlr:
		handleError(w, http.StatusNotImplemented, fmt.Errorf("Action not implemented yet"))
		return
	case ActionDelPlr:
		handleError(w, http.StatusNotImplemented, fmt.Errorf("Action not implemented yet"))
		return
	case ActionChangeName:
		newName := cmd.Arguments.String(ArgumentName)
		err = trn.ChangeName(newName)
		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}
	default:
		err = fmt.Errorf("Action not recognized: %s", cmd.Action)
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	err = trn.Save(s.es, nil)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handlePOSTTournaments(w http.ResponseWriter, r *http.Request) {
	cmd := hyper.ExtractCommand(r)
	switch cmd.Action {
	case ActionCreate:
		trn := NewTournament()
		err = trn.Create(TournamentID(uuid.MakeV4()))
		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}
		err = trn.Save(s.es, nil)
		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}
	default:
		handleError(w, http.StatusInternalServerError, fmt.Errorf("Action not recognized"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) MakeTournamentPlayersHyperItem(tourn Tournament, resolve hyper.ResolverFunc) (hyper.Item, error) {
	plrs := hyper.Item{
		Label: "Participating Players",
		Type:  "players",
	}
	for _, pID := range tourn.Players {
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
		Properties: []hyper.Property{
			{
				Label: "Name",
				Name:  "name",
				Value: trn.Name,
			},
			{
				Label: "ID",
				Name:  "id",
				Value: trn.ID,
			},
			{
				Label: "Version",
				Name:  "version",
				Value: trn.Version,
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
				Label: "Ongoing",
				Name:  "ongoing",
				Value: trn.Ongoing,
			},
			{
				Label: "Finished",
				Name:  "finished",
				Value: trn.Finished,
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
	item.AddLink(link)
	return item
}
