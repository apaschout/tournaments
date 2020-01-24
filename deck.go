package tournaments

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cognicraft/hyper"
	"github.com/cognicraft/uuid"
)

type Deck struct {
	ID   DeckID
	Name string
	Link string
}

type DeckID string

func (s *Server) handleGETDecks(w http.ResponseWriter, r *http.Request) {
	resolve := hyper.ExternalURLResolver(r)
	res := hyper.Item{
		Label: "Decks",
		Type:  "decks",
	}
	links := []hyper.Link{
		{
			Rel:  "self",
			Href: resolve(".").String(),
		},
	}
	s.decks, err = s.p.FindAllDecks()
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	for _, dck := range s.decks {
		item := dck.MakeUndetailedHyperItem(resolve)
		res.AddItem(item)
	}
	res.AddLinks(links)
	hyper.Write(w, http.StatusOK, res)
}

func (s *Server) handleGetDeck(w http.ResponseWriter, r *http.Request) {
	resolve := hyper.ExternalURLResolver(r)
	dID := DeckID(r.Context().Value(":id").(string))

	dck, err := s.p.FindDeckByID(dID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	res := dck.MakeDetailedHyperItem(resolve)
	hyper.Write(w, http.StatusOK, res)
}

func (s *Server) handlePOSTDecks(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	dck := Deck{}
	err = json.Unmarshal(b, &dck)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	ok, err := s.p.IsDeckNameAvailable(dck.Name)
	if !ok {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = uuid.Parse(string(dck.ID))
	if err != nil {
		dck.ID = DeckID(uuid.MakeV4())
	}
	dck.ID = DeckID(uuid.MakeV4())
	// err = s.db.SaveDeck(dck)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (dck *Deck) MakeUndetailedHyperItem(resolve hyper.ResolverFunc) hyper.Item {
	item := hyper.Item{
		Label: dck.Name,
		Type:  "deck",
		ID:    string(dck.ID),
		Properties: []hyper.Property{
			{
				Label: "Name",
				Name:  "name",
				Value: dck.Name,
			},
		},
	}
	dLink := hyper.Link{
		Rel:  "details",
		Href: resolve("./%s", dck.ID).String(),
	}
	item.AddLink(dLink)
	return item
}

func (dck *Deck) MakeDetailedHyperItem(resolve hyper.ResolverFunc) hyper.Item {
	item := hyper.Item{
		Label: dck.Name,
		Type:  "deck",
		Properties: []hyper.Property{
			{
				Label: "Name",
				Name:  "name",
				Value: dck.Name,
			},
			{
				Label: "ID",
				Name:  "id",
				Value: dck.ID,
			},
			{
				Label: "Deck-Builder",
				Name:  "link",
				Value: dck.Link,
			},
		},
	}
	link := hyper.Link{
		Rel:  "self",
		Href: resolve("./%s", dck.ID).String(),
	}
	item.AddLink(link)
	return item
}
