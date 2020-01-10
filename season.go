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
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Start    string   `json:"start,omitempty"`
	End      string   `json:"end,omitempty"`
	Finished bool     `json:"finished,omitempty"`
	Format   string   `json:"format,omitempty"`
	Players  []Player `json:"players,omitempty"`
}

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
		seas.Players, err = s.db.FindPlayersInSeason(seas.Id)
		if err != nil {
			log.Println(err)
			hyper.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		item := hyper.Item{
			Label: seas.Name,
			Type:  "season",
			Properties: []hyper.Property{
				{
					Label: "ID",
					Name:  "id",
					Value: seas.Id,
				},
				{
					Label: "Name",
					Name:  "name",
					Value: seas.Name,
				},
			},
		}
		plrs := hyper.Item{
			Label: "Participating Players",
			Type:  "players",
		}
		for _, plr := range seas.Players {
			item := hyper.Item{
				Label: plr.Name,
				Type:  "player",
				Properties: []hyper.Property{
					{
						Label: "ID",
						Name:  "id",
						Value: plr.Id,
					},
					{
						Label: "Name",
						Name:  "name",
						Value: plr.Name,
					},
				},
			}
			pLink := hyper.Link{
				Rel:  "details",
				Href: resolve("../players/%s", plr.Id).String(),
			}
			item.AddLink(pLink)
			plrs.AddItem(item)
		}
		sLink := hyper.Link{
			Rel:  seas.Name,
			Href: resolve("./%s", seas.Id).String(),
		}
		item.AddItem(plrs)
		item.AddLink(sLink)
		res.AddItem(item)
	}
	res.AddLink(link)
	hyper.Write(w, 200, res)
}

func (s *Server) handleGETSeason(w http.ResponseWriter, r *http.Request) {
	resolve := hyper.ExternalURLResolver(r)
	sID := r.Context().Value(":id").(string)

	seas, err := s.db.FindSeasonByID(sID)
	if err != nil {
		log.Println(err)
		hyper.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	res := hyper.Item{
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
				Value: seas.Id,
			},
		},
	}
	link := hyper.Link{
		Rel:  "self",
		Href: resolve("./%s", seas.Id).String(),
	}
	res.AddLink(link)
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

	err = s.db.CheckForDuplicateName("Seasons", seas.Name)
	if err != nil {
		log.Println(err)
		hyper.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	seas.Id = uuid.MakeV4()
	err = s.db.SaveSeason(seas)
	if err != nil {
		log.Println(err)
		hyper.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
