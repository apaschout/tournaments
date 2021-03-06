package tournaments

import (
	"html/template"
	"log"
	"net/http"

	"github.com/cognicraft/event"
	"github.com/cognicraft/hyper"
	"github.com/cognicraft/mux"
	"github.com/cognicraft/uuid"
	"golang.org/x/crypto/bcrypt"
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
		"itemByType":          itemByType,
		"wins":                wins,
		"getParticipants":     getParticipants,
		"sortParticipants":    sortParticipants,
	}
	templ = template.Must(template.New("server").Funcs(funcMap).ParseGlob("assets/templates/*.html"))

	sub := s.es.SubscribeToStreamFrom(event.All, s.p.GetVersion())
	sub.On(s.p.On)

	//create admin player if server is newly set up
	if s.es.Version(event.All) == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPw), 8)
		if err != nil {
			log.Fatal(err)
			return
		}
		a := NewPlayer(s)
		ID := PlayerID(uuid.MakeV4())
		tID := TrackerID(uuid.MakeV4())
		err = a.Create(ID, tID, "admin", "admin", string(hashedPassword))
		if err != nil {
			log.Fatal(err)
			return
		}
		err = a.ChangeName("admin")
		if err != nil {
			log.Fatal(err)
			return
		}
		err = a.Save(s.es, nil)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	chain := mux.NewChain(
		mux.CORS(mux.AccessControlDefaults),
		mux.GZIP,
		s.refreshToken,
	)
	aChain := mux.NewChain(
		mux.CORS(mux.AccessControlDefaults),
		mux.GZIP,
		s.authentication,
		s.refreshToken,
	)
	s.router.Route("/").GET(chain.ThenFunc(s.handleGETIndex))

	s.router.Route("/api/").GET(chain.ThenFunc(s.handleGETAPI))

	s.router.Route("/api/signup").POST(chain.ThenFunc(s.handleSignUp))
	s.router.Route("/api/signup").GET(chain.ThenFunc(s.handleGETSignUp))
	s.router.Route("/api/signin").POST(chain.ThenFunc(s.handleSignIn))
	s.router.Route("/api/signin").GET(chain.ThenFunc(s.handleGETSignIn))

	s.router.Route("/api/tournaments/").GET(aChain.ThenFunc(s.handleGETTournaments))
	s.router.Route("/api/tournaments/:id").GET(aChain.ThenFunc(s.handleGETTournament))
	s.router.Route("/api/tournaments/:id").POST(aChain.ThenFunc(s.handlePOSTTournament))
	s.router.Route("/api/tournaments/").POST(aChain.ThenFunc(s.handlePOSTTournaments))

	s.router.Route("/api/standings/").GET(aChain.ThenFunc(s.handleGETStandings))
	s.router.Route("/api/standings/:id").GET(aChain.ThenFunc(s.handleGETStanding))

	s.router.Route("/api/players/").GET(aChain.ThenFunc(s.handleGETPlayers))
	s.router.Route("/api/players/:id").GET(aChain.ThenFunc(s.handleGETPlayer))
	s.router.Route("/api/players/:id").POST(aChain.ThenFunc(s.handlePOSTPlayer))

	s.router.Route("/api/trackers/:id").GET(aChain.ThenFunc(s.HandleGETTracker))

	s.router.Route("/api/decks/").GET(aChain.ThenFunc(s.handleGETDecks))
	s.router.Route("/api/decks/:id").GET(aChain.ThenFunc(s.handleGetDeck))
	s.router.Route("/api/decks/").POST(aChain.ThenFunc(s.handlePOSTDecks))

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
