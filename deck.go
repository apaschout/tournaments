package tournaments

import (
	"net/http"

	"github.com/cognicraft/hyper"
)

type Deck struct {
	Id   int
	Name string
	Link string
}

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
	for _, dck := range s.decks {
		item := hyper.Item{
			Label: dck.Name,
			Type:  "deck",
			Properties: []hyper.Property{
				{
					Label: "ID",
					Name:  "id",
					Value: dck.Id,
				},
				{
					Label: "Name",
					Name:  "name",
					Value: dck.Name,
				},
				{
					Label: "Deck-Builder",
					Name:  "link",
					Value: dck.Link,
				},
			},
		}
		dLink := hyper.Link{
			Rel:  dck.Name,
			Href: resolve("./%d", dck.Id).String(),
		}
		item.AddLink(dLink)
		res.AddItem(item)
	}
	res.AddLinks(links)
	hyper.Write(w, http.StatusOK, res)
}