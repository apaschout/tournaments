package tournaments

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cognicraft/hyper"
	"github.com/cognicraft/uuid"
)

type Season struct {
	ID       SeasonID   `json:"id"`
	Name     string     `json:"name"`
	Start    string     `json:"start,omitempty"`
	End      string     `json:"end,omitempty"`
	Ongoing  bool       `json:"ongoing,omitempty"`
	Finished bool       `json:"finished,omitempty"`
	Format   string     `json:"format,omitempty"`
	Players  []PlayerID `json:"players,omitempty"`
}

type SeasonID string

const (
	ActionStart      = "start"
	ActionEnd        = "end"
	ActionAddPlr     = "addPlayer"
	ActionDelPlr     = "delPlayer"
	ActionChangeName = "changeName"
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
	s.seasons, err = s.db.FindAllSeasons()
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	for _, seas := range s.seasons {
		seas.Players, err = s.db.FindPlayersInSeason(seas.ID)
		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}
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

	seas, err := s.db.FindSeasonByID(sID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	seas.Players, err = s.db.FindPlayersInSeason(seas.ID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	res := seas.MakeDetailedHyperItem(resolve)
	plrs, err := s.MakeSeasonPlayersHyperItem(seas, resolve)
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
	seas, err := s.db.FindSeasonByID(sID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	switch cmd.Action {
	case ActionStart:
		if seas.Start != "" {
			err = fmt.Errorf("Season already started.")
			handleError(w, http.StatusInternalServerError, err)
			return
		}
		seas.Begin()
	case ActionEnd:
		if seas.End != "" {
			err = fmt.Errorf("Season already ended.")
			handleError(w, http.StatusInternalServerError, err)
			return
		} else if seas.Start == "" {
			err = fmt.Errorf("Season has not started yet.")
			handleError(w, http.StatusInternalServerError, err)
			return
		}
		seas.Finish()
	case ActionAddPlr:
		pID := PlayerID(cmd.Arguments.String(ArgumentPlayerID))
		ok, err := s.db.PlayerExists(pID)
		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}
		if ok {
			err = seas.AddPlayer(pID)
			if err != nil {
				handleError(w, http.StatusInternalServerError, err)
				return
			}
		}
	case ActionDelPlr:
		pID := PlayerID(cmd.Arguments.String(ArgumentPlayerID))
		err = seas.RemovePlayer(pID)
		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}
	case ActionChangeName:
		newName := cmd.Arguments.String(ArgumentName)
		ok, err := s.db.SeasonNameAvailable(newName)
		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}
		if ok {
			seas.ChangeName(newName)
		}
	default:
		err = fmt.Errorf("Action not recognized: %s", cmd.Action)
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	err = s.db.SaveSeason(seas)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
	}
}

func (s *Server) handlePOSTSeasons(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	seas := Season{}
	err = json.Unmarshal(b, &seas)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	ok, err := s.db.SeasonNameAvailable(seas.Name)
	if !ok {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = uuid.Parse(string(seas.ID))
	if err != nil {
		seas.ID = SeasonID(uuid.MakeV4())
	}
	err = s.db.SaveSeason(seas)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
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
		plr, err := s.db.FindPlayerByID(pID)
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

func (seas *Season) Begin() {
	date := time.Now()
	seas.Start = date.Format("2006-01-02 15:04:05Z")
	seas.Ongoing = true
}

func (seas *Season) Finish() {
	date := time.Now()
	seas.End = date.Format("2006-01-02 15:04:05Z")
	seas.Ongoing = false
	seas.Finished = true
}

func (seas *Season) AddPlayer(pID PlayerID) error {
	for _, p := range seas.Players {
		if p == pID {
			err = fmt.Errorf("%s already exists in Season: %s", pID, seas.ID)
			return err
		}
	}
	seas.Players = append(seas.Players, pID)

	return nil
}

func (seas *Season) RemovePlayer(pID PlayerID) error {
	for i, p := range seas.Players {
		if p == pID {
			seas.Players = append(seas.Players[:i], seas.Players[i+1:]...)
			return nil
		}
	}
	err = fmt.Errorf("%s not found in Season: %s", pID, seas.ID)
	return err
}

func (seas *Season) ChangeName(newName string) {
	seas.Name = newName
}
