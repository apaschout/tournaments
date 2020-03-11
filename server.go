package tournaments

import (
	"html/template"
	"log"
	"net/http"

	"github.com/cognicraft/event"
	"github.com/cognicraft/hyper"
	"github.com/cognicraft/mux"
)

var (
	err   error
	templ *template.Template
)

type Server struct {
	router      *mux.Router
	p           Projection
	es          *event.Store
	tournaments []Tournament
	players     []Player
	decks       []Deck
}

func NewServer(p Projection, es *event.Store) *Server {
	s := Server{
		router: mux.New(),
		p:      p,
		es:     es,
	}
	s.init()
	return &s
}

func (s *Server) init() {
	funcMap := template.FuncMap{
		"add":                 func(a, b int) int { return a + b },
		"details":             getDetails,
		"seat":                createSeatForIndex,
		"action":              actionByRel,
		"participantNameByID": participantNameByID,
		"propertyByName":      propertyByName,
		"wins":                wins,
		"getParticipants":     getParticipants,
		"sortParticipants":    sortParticipants,
	}
	templ = template.Must(template.New("server").Funcs(funcMap).ParseGlob("assets/templates/*.html"))

	sub := s.es.SubscribeToStreamFrom(event.All, s.p.GetVersion())
	sub.On(s.p.On)

	chain := mux.NewChain(
		mux.CORS(mux.AccessControlDefaults),
		mux.GZIP,
	)
	s.router.Route("/").GET(chain.ThenFunc(s.handleGETIndex))

	s.router.Route("/api/").GET(chain.ThenFunc(s.handleGETAPI))

	s.router.Route("/api/tournaments/").GET(chain.ThenFunc(s.handleGETTournaments))
	s.router.Route("/api/tournaments/:id").GET(chain.ThenFunc(s.handleGETTournament))
	s.router.Route("/api/tournaments/:id").POST(chain.ThenFunc(s.handlePOSTTournament))
	s.router.Route("/api/tournaments/").POST(chain.ThenFunc(s.handlePOSTTournaments))

	s.router.Route("/api/players/").GET(chain.ThenFunc(s.handleGETPlayers))
	s.router.Route("/api/players/:id").GET(chain.ThenFunc(s.handleGETPlayer))
	s.router.Route("/api/players/:id").POST(chain.ThenFunc(s.handlePOSTPlayer))
	s.router.Route("/api/players/").POST(chain.ThenFunc(s.handlePOSTPlayers))

	s.router.Route("/api/decks/").GET(chain.ThenFunc(s.handleGETDecks))
	s.router.Route("/api/decks/:id").GET(chain.ThenFunc(s.handleGetDeck))
	s.router.Route("/api/decks/").POST(chain.ThenFunc(s.handlePOSTDecks))

	s.router.Route("/js/:file").GET(http.StripPrefix("/js/", http.FileServer(http.Dir("./assets/js"))))
	s.router.Route("/css/:file").GET(http.StripPrefix("/css/", http.FileServer(http.Dir("./assets/css"))))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) handleGETAPI(w http.ResponseWriter, r *http.Request) {
	resolve := hyper.ExternalURLResolver(r)
	res := hyper.Item{
		Label: "API",
	}
	links := []hyper.Link{
		{
			Rel:  "self",
			Href: resolve(".").String(),
		},
		{
			Rel:  "tournaments",
			Href: resolve("./seasons/").String(),
		},
		{
			Rel:  "players",
			Href: resolve("./players/").String(),
		},
	}
	res.AddLinks(links)
	hyper.Write(w, http.StatusOK, res)
}

func (s *Server) handleGETIndex(w http.ResponseWriter, r *http.Request) {

}

func handleError(w http.ResponseWriter, status int, err error, isHtmlReq bool) {
	log.Println(err)
	if isHtmlReq {
		err = templ.ExecuteTemplate(w, "error.html", err)
		if err != nil {
			writeError(w, status, err)
		}
	} else {
		writeError(w, status, err)
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	type errorCoder interface {
		Code() string
	}

	e := hyper.Error{}
	e.Message = err.Error()
	if errC, ok := err.(errorCoder); ok {
		e.Code = errC.Code()
	}

	hyper.Write(w, status, hyper.Item{Errors: hyper.Errors{e}})
}
