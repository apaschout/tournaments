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

type Player struct {
	ID          PlayerID       `json:"id,omitempty"`
	Version     uint64         `json:"version"`
	Name        string         `json:"name"`
	Tournaments []TournamentID `json:"tournaments"`
	Tracker     TrackerID      `json:"tracker"`
	*event.ChangeRecorder
	Server *Server
}

type PlayerID string

func (s *Server) handleGETPlayers(w http.ResponseWriter, r *http.Request) {
	isHtmlReq := strings.Contains(r.Header.Get("Accept"), "text/html")
	resolve := hyper.ExternalURLResolver(r)
	res := hyper.Item{
		Label: "Players",
		Type:  "players",
	}
	link := hyper.Link{
		Rel:  "self",
		Href: resolve(".").String(),
	}
	s.players, err = s.p.FindAllPlayers()
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, isHtmlReq)
		return
	}
	for _, plr := range s.players {
		item := plr.MakeUndetailedHyperItem(resolve)
		res.AddItem(item)
	}
	res.AddLink(link)
	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		err = templ.ExecuteTemplate(w, "players.html", res)
		if err != nil {
			log.Println(err)
		}
	} else {
		hyper.Write(w, http.StatusOK, res)
	}
}

func (s *Server) handleGETPlayer(w http.ResponseWriter, r *http.Request) {
	isHtmlReq := strings.Contains(r.Header.Get("Accept"), "text/html")
	resolve := hyper.ExternalURLResolver(r)
	pID := PlayerID(r.Context().Value(":id").(string))
	plr, err := LoadPlayer(s, pID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, isHtmlReq)
		return
	}

	res := plr.MakeDetailedHyperItem(resolve)
	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		err = templ.ExecuteTemplate(w, "player.html", res)
		if err != nil {
			log.Println(err)
		}
	} else {
		hyper.Write(w, http.StatusOK, res)
	}
}

func (s *Server) handlePOSTPlayer(w http.ResponseWriter, r *http.Request) {
	isHtmlReq := strings.Contains(r.Header.Get("Accept"), "text/html")
	cmd := hyper.ExtractCommand(r)
	pID := PlayerID(r.Context().Value(":id").(string))

	plr, err := LoadPlayer(s, pID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, isHtmlReq)
		return
	}
	switch cmd.Action {
	case ActionChangeName:
		newName := cmd.Arguments.String(ArgumentName)
		ok, err := s.p.IsPlayerNameAvailable(newName)
		if !ok {
			handleError(w, http.StatusInternalServerError, err, isHtmlReq)
			return
		}
		err = plr.ChangeName(newName)
		if err != nil {
			handleError(w, http.StatusInternalServerError, err, isHtmlReq)
			return
		}
	default:
		err = fmt.Errorf("Action not recognized: %s", cmd.Action)
		handleError(w, http.StatusInternalServerError, err, isHtmlReq)
		return
	}
	err = plr.Save(s.es, nil)
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

func (s *Server) handlePOSTPlayers(w http.ResponseWriter, r *http.Request) {
	isHtmlReq := strings.Contains(r.Header.Get("Accept"), "text/html")
	cmd := hyper.ExtractCommand(r)
	switch cmd.Action {
	case ActionCreate:
		plr := NewPlayer(s)
		err = plr.Create(PlayerID(uuid.MakeV4()), TrackerID(uuid.MakeV4()))
		if err != nil {
			handleError(w, http.StatusInternalServerError, err, isHtmlReq)
			return
		}
		err = plr.Save(s.es, nil)
		if err != nil {
			handleError(w, http.StatusInternalServerError, err, isHtmlReq)
			return
		}
	default:
		handleError(w, http.StatusInternalServerError, fmt.Errorf("Action not recognized"), isHtmlReq)
		return
	}
	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		http.Redirect(w, r, r.URL.String(), http.StatusSeeOther)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (plr *Player) MakeUndetailedHyperItem(resolve hyper.ResolverFunc) hyper.Item {
	item := hyper.Item{
		Label: plr.Name,
		Type:  "player",
		ID:    string(plr.ID),
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
		Href: resolve("./%s", plr.ID).String(),
	}
	item.AddLink(pLink)
	return item
}

func (plr *Player) MakeDetailedHyperItem(resolve hyper.ResolverFunc) hyper.Item {
	item := hyper.Item{
		Label: plr.Name,
		Type:  "player",
		ID:    string(plr.ID),
		Properties: []hyper.Property{
			{
				Label: "Name",
				Name:  "name",
				Value: plr.Name,
			},
			{
				Label: "Tournaments",
				Name:  "tournaments",
				Value: plr.Tournaments,
			},
			{
				Label: "Tracker",
				Name:  "tracker",
				Value: plr.Tracker,
			},
		},
	}
	actions := hyper.Actions{
		{
			Label:  "Change Name",
			Rel:    ActionChangeName,
			Href:   resolve("./%s", plr.ID).String(),
			Method: "POST",
			Parameters: hyper.Parameters{
				{
					Name:        ArgumentName,
					Placeholder: "New Name...",
				},
			},
		},
	}
	link := hyper.Link{
		Rel:  "self",
		Href: resolve("./%s", plr.ID).String(),
	}
	item.AddActions(actions)
	item.AddLink(link)
	return item
}
