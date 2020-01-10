package tournaments

import (
	"net/http"

	"github.com/cognicraft/hyper"
	"github.com/cognicraft/mux"
)

type Server struct {
	router  *mux.Router
	db      Datastore
	seasons []Season
	players []Player
	decks   []Deck
}

func NewServer(db Datastore) *Server {
	return &Server{
		router: mux.New(),
		db:     db,
	}
}

func (s *Server) Init() {

	chain := mux.NewChain(
		mux.CORS(mux.AccessControlDefaults),
		mux.GZIP,
	)
	s.router.Route("/api/").GET(chain.ThenFunc(s.handleGETAPI))

	s.router.Route("/api/seasons/").GET(chain.ThenFunc(s.handleGETSeasons))
	s.router.Route("/api/seasons/:id").GET(chain.ThenFunc(s.handleGETSeason))
	s.router.Route("/api/seasons/").POST(chain.ThenFunc(s.handlePOSTSeasons))

	s.router.Route("/api/players/").GET(chain.ThenFunc(s.handleGETPlayers))
	s.router.Route("/api/players/:id").GET(chain.ThenFunc(s.handleGETPlayer))
	s.router.Route("/api/players/").POST(chain.ThenFunc(s.handlePOSTPlayers))

	s.router.Route("/api/decks/").GET(chain.ThenFunc(s.handleGETDecks))
	s.router.Route("/api/decks/:id").GET(chain.ThenFunc(s.handleGetDeck))
	s.router.Route("/api/decks/").POST(chain.ThenFunc(s.handlePOSTDecks))
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
			Rel:  "seasons",
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
