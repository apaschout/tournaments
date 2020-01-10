package tournaments

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cognicraft/hyper"
	"github.com/cognicraft/uuid"
)

type Season struct {
	ID       SeasonID   `json:"id"`
	Name     string     `json:"name"`
	Start    string     `json:"start,omitempty"`
	End      string     `json:"end,omitempty"`
	Finished bool       `json:"finished,omitempty"`
	Format   string     `json:"format,omitempty"`
	Players  []PlayerID `json:"players,omitempty"`
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
		log.Println(err)
		hyper.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	for _, seas := range s.seasons {
		seas.Players, err = s.db.FindPlayersInSeason(seas.ID)
		if err != nil {
			log.Println(err)
			hyper.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		item := seas.MakeUndetailedSeasonHyperItem(resolve)
		plrs := seas.MakePlayersHyperItem(resolve)
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
		log.Println(err)
		hyper.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	seas.Players, err = s.db.FindPlayersInSeason(seas.ID)
	if err != nil {
		log.Println(err)
		hyper.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	res := seas.MakeDetailedSeasonHyperItem(resolve)
	plrs := seas.MakePlayersHyperItem(resolve)
	res.AddItem(plrs)
	hyper.Write(w, http.StatusOK, res)
}

func (s *Server) handlePOSTSeasons(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		hyper.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	seas := Season{}
	err = json.Unmarshal(b, &seas)
	if err != nil {
		log.Println(err)
		hyper.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	ok, err := s.db.SeasonNameAvailable(seas.Name)
	if !ok {
		log.Println(err)
		hyper.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = uuid.Parse(string(seas.ID))
	if err != nil {
		seas.ID = SeasonID(uuid.MakeV4())
	}
	err = s.db.SaveSeason(seas)
	if err != nil {
		log.Println(err)
		hyper.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (seas *Season) MakePlayersHyperItem(resolve hyper.ResolverFunc) hyper.Item {
	plrs := hyper.Item{
		Label: "Participating Players",
		Type:  "players",
	}
	for _, plrID := range seas.Players {
		item := hyper.Item{
			Label: "Player",
			Type:  "player",
			Properties: []hyper.Property{
				{
					Label: "ID",
					Name:  "id",
					Value: plrID,
				},
			},
		}
		pLink := hyper.Link{
			Rel:  "details",
			Href: resolve("../players/%s", plrID).String(),
		}
		item.AddLink(pLink)
		plrs.AddItem(item)
	}
	return plrs
}

func (seas *Season) MakeUndetailedSeasonHyperItem(resolve hyper.ResolverFunc) hyper.Item {
	item := hyper.Item{
		Label: seas.Name,
		Type:  "season",
		Properties: []hyper.Property{
			{
				Label: "ID",
				Name:  "id",
				Value: seas.ID,
			},
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

func (seas *Season) MakeDetailedSeasonHyperItem(resolve hyper.ResolverFunc) hyper.Item {
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
