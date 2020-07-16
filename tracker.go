package tournaments

import (
	"net/http"
	"strings"

	"github.com/cognicraft/event"
	"github.com/cognicraft/hyper"
)

type TrackerID string

type Tracker struct {
	ID        TrackerID `json:"id"`
	Player    PlayerID  `json:"player"`
	Version   uint64    `json:"version"`
	Matches   int       `json:"matches"`
	MatchWins int       `json:"matchWins"`
	Games     int       `json:"games"`
	GameWins  int       `json:"gameWins"`
	*event.ChangeRecorder
}

func (s *Server) HandleGETTracker(w http.ResponseWriter, r *http.Request) {
	isHtmlReq := strings.Contains(r.Header.Get("Accept"), "text/html")
	resolve := hyper.ExternalURLResolver(r)
	tID := TrackerID(r.Context().Value(":id").(string))

	trk, err := LoadTracker(s.es, tID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, isHtmlReq)
		return
	}

	res := trk.MakeTrackerHyperItem(resolve)

	if isHtmlReq {
		err = templ.ExecuteTemplate(w, "tracker.html", res)
		if err != nil {
			handleError(w, http.StatusInternalServerError, err, isHtmlReq)
			return
		}
	} else {
		hyper.Write(w, http.StatusOK, res)
	}
}

func (trk *Tracker) MakeTrackerHyperItem(resolve hyper.ResolverFunc) hyper.Item {
	res := hyper.Item{
		Label: "Tracker",
		Type:  "tracker",
		ID:    string(trk.ID),
		Properties: hyper.Properties{
			{
				Label: "Player",
				Name:  "player",
				Value: trk.Player,
			},
			{
				Label: "Matches",
				Name:  "matches",
				Value: trk.Matches,
			},
			{
				Label: "Match Wins",
				Name:  "matchWins",
				Value: trk.MatchWins,
			},
			{
				Label: "Games",
				Name:  "games",
				Value: trk.Games,
			},
			{
				Label: "Game Wins",
				Name:  "gameWins",
				Value: trk.GameWins,
			},
		},
	}
	plrDetails := hyper.Link{
		Rel:  "details",
		Href: resolve("/api/players/%s", trk.Player).String(),
	}
	l := hyper.Link{
		Rel:  hyper.RelSelf,
		Href: resolve("./%s", trk.ID).String(),
	}
	res.AddLink(plrDetails)
	res.AddLink(l)
	return res
}
