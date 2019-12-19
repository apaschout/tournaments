package tournaments

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/cognicraft/hyper"
	"github.com/cognicraft/mux"

	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	router  *mux.Router
	db      *sql.DB
	seasons []Season
	players []Player
}

var (
	err error
)

func NewServer() *Server {
	return &Server{
		router: mux.New(),
	}
}

func (s *Server) Init() {
	s.InitDB()

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
}

func (s *Server) InitDB() {
	s.db, err = sql.Open("sqlite3", "../../db/tournaments.db")
	if err != nil {
		log.Fatal(err)
	}

	s.LoadSeasons()
	s.LoadPlayers()
}

func (s *Server) LoadSeasons() {
	sQuery := "SELECT * FROM Seasons"
	rows, err := s.db.Query(sQuery)
	if err != nil {
		log.Fatalf("Query: %v\n", err)
	}
	s.seasons = []Season{}
	for rows.Next() {
		seas := Season{}
		err = rows.Scan(&seas.Id, &seas.Name)
		if err != nil {
			log.Fatalf("Scan: %v\n", err)
		}
		s.seasons = append(s.seasons, seas)
	}
}

func (s *Server) LoadPlayers() {
	pQuery := "SELECT * FROM Players"
	rows, err := s.db.Query(pQuery)
	if err != nil {
		log.Fatalf("Query: %v\n", err)
	}
	s.players = []Player{}
	for rows.Next() {
		plr := Player{}
		err = rows.Scan(&plr.Id, &plr.Name)
		if err != nil {
			log.Fatalf("Scan: %v\n", err)
		}
		s.players = append(s.players, plr)
	}
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
