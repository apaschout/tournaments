package tournaments

import (
	"fmt"
	"net/http"

	"github.com/cognicraft/event"
	"github.com/cognicraft/hyper"
	"github.com/cognicraft/uuid"
)

type Season struct {
	ID       SeasonID   `json:"id"`
	Version  uint64     `json:"version"`
	Name     string     `json:"name"`
	Start    string     `json:"start,omitempty"`
	End      string     `json:"end,omitempty"`
	Ongoing  bool       `json:"ongoing,omitempty"`
	Finished bool       `json:"finished,omitempty"`
	Format   string     `json:"format,omitempty"`
	Players  []PlayerID `json:"players,omitempty"`
	*event.ChangeRecorder
}

type SeasonID string

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

func (s *Server) handleGETSeasons(w http.ResponseWriter, r *http.Request) {
	resolve := hyper.ExternalURLResolver(r)
	res := hyper.Item{
		Label: "Seasons",
		Type:  "seasons",
	}
	link := hyper.Link{
		Rel:  "self",
		Href: resolve(".").String(),
	}
	s.seasons, err = s.p.FindAllSeasons()
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	for _, seas := range s.seasons {
		// seas.Players, err = s.db.FindPlayersInSeason(seas.ID)
		// if err != nil {
		// 	handleError(w, http.StatusInternalServerError, err)
		// 	return
		// }
		item := seas.MakeUndetailedHyperItem(resolve)
		plrs, err := s.MakeSeasonPlayersHyperItem(seas, resolve)
		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}
		item.AddItem(plrs)
		res.AddItem(item)
	}
	res.AddLink(link)
	hyper.Write(w, 200, res)
}

func (s *Server) handleGETSeason(w http.ResponseWriter, r *http.Request) {
	resolve := hyper.ExternalURLResolver(r)
	sID := SeasonID(r.Context().Value(":id").(string))

	// seas, err := s.db.FindSeasonByID(sID)
	seas, err := LoadSeason(s.es, sID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	// seas.Players, err = s.db.FindPlayersInSeason(seas.ID)
	// if err != nil {
	// 	handleError(w, http.StatusInternalServerError, err)
	// 	return
	// }

	res := seas.MakeDetailedHyperItem(resolve)
	plrs, err := s.MakeSeasonPlayersHyperItem(*seas, resolve)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	res.AddItem(plrs)
	hyper.Write(w, http.StatusOK, res)
}

func (s *Server) handlePOSTSeason(w http.ResponseWriter, r *http.Request) {
	cmd := hyper.ExtractCommand(r)
	sID := SeasonID(r.Context().Value(":id").(string))
	// seas, err := s.db.FindSeasonByID(sID)
	seas, err := LoadSeason(s.es, sID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	switch cmd.Action {
	case ActionStart:
		err = seas.Begin()
		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}
	case ActionEnd:
		err = seas.Finish()
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
		err = seas.ChangeName(newName)
		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}
	default:
		err = fmt.Errorf("Action not recognized: %s", cmd.Action)
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	err = seas.Save(s.es, nil)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handlePOSTSeasons(w http.ResponseWriter, r *http.Request) {
	cmd := hyper.ExtractCommand(r)
	switch cmd.Action {
	case ActionCreate:
		seas := NewSeason()
		err = seas.Create(SeasonID(uuid.MakeV4()))
		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}
		err = seas.Save(s.es, nil)
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

func (s *Server) MakeSeasonPlayersHyperItem(seas Season, resolve hyper.ResolverFunc) (hyper.Item, error) {
	plrs := hyper.Item{
		Label: "Participating Players",
		Type:  "players",
	}
	for _, pID := range seas.Players {
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

func (seas *Season) MakeUndetailedHyperItem(resolve hyper.ResolverFunc) hyper.Item {
	item := hyper.Item{
		Label: seas.Name,
		Type:  "season",
		ID:    string(seas.ID),
		Properties: []hyper.Property{
			{
				Label: "Name",
				Name:  "name",
				Value: seas.Name,
			},
		},
	}
	sLink := hyper.Link{
		Rel:  "details",
		Href: resolve("./%s", seas.ID).String(),
	}
	item.AddLink(sLink)
	return item
}

func (seas *Season) MakeDetailedHyperItem(resolve hyper.ResolverFunc) hyper.Item {
	item := hyper.Item{
		Label: seas.Name,
		Type:  "season",
		Properties: []hyper.Property{
			{
				Label: "Name",
				Name:  "name",
				Value: seas.Name,
			},
			{
				Label: "ID",
				Name:  "id",
				Value: seas.ID,
			},
			{
				Label: "Version",
				Name:  "version",
				Value: seas.Version,
			},
			{
				Label: "Start",
				Name:  "start",
				Value: seas.Start,
			},
			{
				Label: "End",
				Name:  "end",
				Value: seas.End,
			},
			{
				Label: "Ongoing",
				Name:  "ongoing",
				Value: seas.Ongoing,
			},
			{
				Label: "Finished",
				Name:  "finished",
				Value: seas.Finished,
			},
			{
				Label: "Format",
				Name:  "format",
				Value: seas.Format,
			},
		},
	}
	link := hyper.Link{
		Rel:  "self",
		Href: resolve("./%s", seas.ID).String(),
	}
	item.AddLink(link)
	return item
}
