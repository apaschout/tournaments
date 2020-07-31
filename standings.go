package tournaments

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/cognicraft/hyper"
)

func (s *Server) handleGETStandings(w http.ResponseWriter, r *http.Request) {
	isHtmlReq := strings.Contains(r.Header.Get("Accept"), "text/html")
	resolve := hyper.ExternalURLResolver(r)
	res := hyper.Item{
		Label: "Tournaments",
		Type:  "tournaments",
	}
	selfLink := hyper.Link{
		Rel:  "self",
		Href: resolve(".").String(),
	}
	s.tournaments, err = s.p.FindAllTournaments()
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, isHtmlReq)
		return
	}
	for _, trn := range s.tournaments {
		item := trn.MakeUndetailedHyperItem(resolve)
		res.AddItem(item)
	}
	res.AddLink(selfLink)
	if isHtmlReq {
		err = templ.ExecuteTemplate(w, "standings.html", res)
		if err != nil {
			log.Println(err)
		}
	} else {
		hyper.Write(w, 200, res)
	}
}

func (s *Server) handleGETStanding(w http.ResponseWriter, r *http.Request) {
	isHtmlReq := strings.Contains(r.Header.Get("Accept"), "text/html")
	resolve := hyper.ExternalURLResolver(r)
	tID := TournamentID(r.Context().Value(":id").(string))
	selfLink := hyper.Link{
		Rel:  "self",
		Href: resolve(".").String(),
	}
	trn, err := LoadTournament(s, tID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, isHtmlReq)
		return
	}
	// s.tournaments, err = s.p.FindAllTournaments()
	// if err != nil {
	// 	handleError(w, http.StatusInternalServerError, err, isHtmlReq)
	// 	return
	// }
	// trnsItem := hyper.Item{
	// 	Label: "Tournaments",
	// 	Type:  "tournaments",
	// }
	// for _, v := range s.tournaments {
	// 	trnItem := v.MakeUndetailedHyperItem(resolve)
	// 	trnsItem.AddItem(trnItem)
	// }
	res := hyper.Item{
		Label: fmt.Sprintf("Standings for %s", trn.Name),
		Rel:   string(trn.ID),
	}
	parts, err := trn.MakeParticipantsHyperItem(resolve)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, isHtmlReq)
		return
	}
	// res.AddItem(trnsItem)
	res.AddItem(parts)
	res.AddLink(selfLink)
	if isHtmlReq {
		err = templ.ExecuteTemplate(w, "standing.html", res)
		if err != nil {
			handleError(w, http.StatusInternalServerError, err, isHtmlReq)
			return
		}
	} else {
		hyper.Write(w, http.StatusOK, res)
	}
}
