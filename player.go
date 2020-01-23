package tournaments

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cognicraft/hyper"
	"github.com/cognicraft/uuid"
)

type Player struct {
	ID   PlayerID `json:"id,omitempty"`
	Name string   `json:"name"`
}

type PlayerID string

func (s *Server) handleGETPlayers(w http.ResponseWriter, r *http.Request) {
	resolve := hyper.ExternalURLResolver(r)
	res := hyper.Item{
		Label: "Players",
		Type:  "players",
	}
	link := hyper.Link{
		Rel:  "self",
		Href: resolve(".").String(),
	}
	s.players, err = s.db.FindAllPlayers()
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	for _, plr := range s.players {
		item := plr.MakeUndetailedHyperItem(resolve)
		res.AddItem(item)
	}
	res.AddLink(link)
	hyper.Write(w, http.StatusOK, res)
}

func (s *Server) handleGETPlayer(w http.ResponseWriter, r *http.Request) {
	resolve := hyper.ExternalURLResolver(r)
	plr := Player{}
	pID := PlayerID(r.Context().Value(":id").(string))
	plr, err = s.db.FindPlayerByID(pID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	res := plr.MakeDetailedHyperItem(resolve)
	hyper.Write(w, http.StatusOK, res)
}

func (s *Server) handlePOSTPlayer(w http.ResponseWriter, r *http.Request) {
	cmd := hyper.ExtractCommand(r)
	pID := PlayerID(r.Context().Value(":id").(string))
	plr, err := s.db.FindPlayerByID(pID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	switch cmd.Action {
	case ActionChangeName:
		newName := cmd.Arguments.String(ArgumentName)
		ok, err := s.db.PlayerNameAvailable(newName)
		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}
		if ok {
			plr.ChangeName(newName)
		}
	default:
		err = fmt.Errorf("Action not recognized: %s", cmd.Action)
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	err = s.db.SavePlayer(plr)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
	}
}

func (s *Server) handlePOSTPlayers(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	plr := Player{}
	err = json.Unmarshal(b, &plr)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	ok, err := s.db.PlayerNameAvailable(plr.Name)
	if !ok {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = uuid.Parse(string(plr.ID))
	if err != nil {
		plr.ID = PlayerID(uuid.MakeV4())
	}
	plr.ID = PlayerID(uuid.MakeV4())
	err = s.db.SavePlayer(plr)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
		Properties: []hyper.Property{
			{
				Label: "Name",
				Name:  "name",
				Value: plr.Name,
			},
			{
				Label: "ID",
				Name:  "id",
				Value: plr.ID,
			},
		},
	}
	link := hyper.Link{
		Rel:  "self",
		Href: resolve("./%s", plr.ID).String(),
	}
	item.AddLink(link)
	return item
}

func (plr *Player) ChangeName(newName string) {
	plr.Name = newName
}
