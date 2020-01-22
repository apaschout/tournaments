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
	ID       SeasonID       `json:"id"`
	Name     string         `json:"name"`
	Start    string         `json:"start,omitempty"`
	End      string         `json:"end,omitempty"`
	Ongoing  bool           `json:"ongoing,omitempty"`
	Finished bool           `json:"finished,omitempty"`
	Format   string         `json:"format,omitempty"`
	Players  []SeasonPlayer `json:"players,omitempty"`
}

type SeasonPlayer struct {
	ID   PlayerID `json:"id"`
	Deck DeckID   `json:"deck,omitempty"`
}

type PatchPayload struct {
	Action string `json:"action"`
}

type SeasonID string

func (s *Server) handleGETSeasons(w http.ResponseWriter, r *http.Request) {
	var err error
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

func (s *Server) handlePATCHSeason(w http.ResponseWriter, r *http.Request) {
	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	pl := PatchPayload{}
	err = json.Unmarshal(jsn, &pl)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	sID := SeasonID(r.Context().Value(":id").(string))
	seas, err := s.db.FindSeasonByID(sID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	switch pl.Action {
	case "start":
		if seas.Start != "" {
			err = fmt.Errorf("Season already started.")
			handleError(w, http.StatusInternalServerError, err)
			return
		}
		seas.Begin()
	case "end":
		if seas.End != "" {
			err = fmt.Errorf("Season already ended.")
			handleError(w, http.StatusInternalServerError, err)
			return
		}
		seas.Finish()
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
	for _, sPlr := range seas.Players {
		plr, err := s.db.FindPlayerByID(sPlr.ID)
		if err != nil {
			return plrs, err
		}
		item := hyper.Item{
			Label: "Player",
			Type:  "player",
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
			Href: resolve("../players/%s", sPlr.ID).String(),
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
	seas.Start = date.Format("2006-01-02T15:04:05")
	seas.Ongoing = true
}

func (seas *Season) Finish() {
	date := time.Now()
	seas.End = date.Format("2006-01-02T15:04:05")
	seas.Ongoing = false
	seas.Finished = true
}
