package tournaments

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cognicraft/hyper"
)

type Player struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name"`
}

func (s *Server) handleGETPlayers(w http.ResponseWriter, r *http.Request) {
	var err error
	resolve := hyper.ExternalURLResolver(r)
	res := hyper.Item{
		Label: "Players",
		Type:  "players",
	}
	links := []hyper.Link{
		{
			Rel:  "self",
			Href: resolve(".").String(),
		},
	}
	s.players, err = s.db.FindAllPlayers()
	if err != nil {
		log.Println(err)
		hyper.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	for _, plr := range *s.players {
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
			Rel:  plr.Name,
			Href: resolve("./%s", plr.Id).String(),
		}
		item.AddLink(pLink)
		res.AddItem(item)
	}
	res.AddLinks(links)
	hyper.Write(w, http.StatusOK, res)
}

func (s *Server) handleGETPlayer(w http.ResponseWriter, r *http.Request) {
	resolve := hyper.ExternalURLResolver(r)
	pID := r.Context().Value(":id").(string)
	plr, err := s.db.FindPlayerByID(pID)
	if err != nil {
		log.Println(err)
		hyper.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	res := hyper.Item{
		Label: plr.Name,
		Type:  "player",
		Properties: []hyper.Property{
			{
				Label: "Name",
				Name:  "name",
				Value: plr.Name,
			},
			{
				Label: "ID",
				Name:  "id",
				Value: plr.Id,
			},
		},
	}
	link := hyper.Link{
		Rel:  "self",
		Href: resolve("./%s", plr.Id).String(),
	}
	res.AddLink(link)
	hyper.Write(w, http.StatusOK, res)
}

func (s *Server) handlePOSTPlayers(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		hyper.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	plr := Player{}
	err = json.Unmarshal(b, &plr)
	if err != nil {
		log.Println(err)
		hyper.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	err = s.db.CreatePlayer(&plr)
	if err != nil {
		log.Println(err)
		hyper.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
