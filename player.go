package tournaments

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/cognicraft/hyper"
)

type Player struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name"`
}

func (s *Server) handleGETPlayers(w http.ResponseWriter, r *http.Request) {
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
	for _, plr := range s.players {
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
			Href: resolve("./%d", plr.Id).String(),
		}
		item.AddLink(pLink)
		res.AddItem(item)
	}
	res.AddLinks(links)
	hyper.Write(w, http.StatusOK, res)
}

func (s *Server) handleGETPlayer(w http.ResponseWriter, r *http.Request) {
	resolve := hyper.ExternalURLResolver(r)
	pID, err := strconv.Atoi(r.Context().Value(":id").(string))
	if err != nil {
		log.Println(err)
		hyper.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	var plr Player
	for _, p := range s.players {
		if p.Id == pID {
			plr = p
			break
		}
	}
	if plr.Id == 0 {
		err := fmt.Errorf("Couldn't find PlayerID %d", pID)
		hyper.WriteError(w, http.StatusNotFound, err)
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
		Href: resolve("./%d", plr.Id).String(),
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

	query := "INSERT INTO Players (name) VALUES (?)"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		log.Printf("Prepare: %v\n", err)
		hyper.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	_, err = stmt.Exec(plr.Name)
	if err != nil {
		log.Printf("Exec: %v\n", err)
		hyper.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	s.LoadPlayers()

	w.WriteHeader(http.StatusNoContent)
}
