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

type Season struct {
	Name    string   `json:"name,omitempty"`
	Id      int      `json:"id"`
	Players []Player `json:"players"`
}

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
	for _, seas := range s.seasons {
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
				{
					Label: "Players",
					Name:  "players",
					Value: seas.Players,
				},
			},
		}
		sLink := hyper.Link{
			Rel:  seas.Name,
			Href: resolve("./%d", seas.Id).String(),
		}
		item.AddLink(sLink)
		res.AddItem(item)
	}
	res.AddLink(link)
	hyper.Write(w, 200, res)
}

func (s *Server) handleGETSeason(w http.ResponseWriter, r *http.Request) {
	resolve := hyper.ExternalURLResolver(r)
	sID, err := strconv.Atoi(r.Context().Value(":id").(string))
	if err != nil {
		log.Println(err)
		hyper.WriteError(w, http.StatusInternalServerError, err)
	}
	var seas Season
	for _, v := range s.seasons {
		if v.Id == sID {
			seas = v
			break
		}
	}
	if seas.Id == 0 {
		err = fmt.Errorf("Couldn't find SeasonID %d", sID)
		hyper.WriteError(w, http.StatusNotFound, err)
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
		Href: resolve("./%d", seas.Id).String(),
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

	query := "INSERT INTO Seasons (name) VALUES (?)"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		log.Printf("Prepare: %v\n", err)
		hyper.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	_, err = stmt.Exec(seas.Name)
	if err != nil {
		log.Printf("Exec: %v\n", err)
		hyper.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	s.LoadSeasons()

	w.WriteHeader(http.StatusNoContent)
}
