package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cognicraft/hyper"
	"github.com/cognicraft/mux"
)

var (
	err   error
	templ *template.Template
)

type TemplateData struct {
	TournamentsData hyper.Item
	PlayersData     hyper.Item
	Data            hyper.Item
}

func main() {
	templ, err = templ.ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}
	router := mux.New()
	chain := mux.NewChain()

	router.Route("/").GET(chain.ThenFunc(index))
	router.Route("/tournament").GET(chain.ThenFunc(handleGETTournament))
	router.Route("/tournament").POST(chain.ThenFunc(handlePOSTTournament))
	router.Route("/player").GET(chain.ThenFunc(handleGETPlayer))
	router.Route("/players").GET(chain.ThenFunc(handleGETPlayers))
	router.Route("/tournaments").GET(chain.ThenFunc(handleGETTournaments))
	router.Route("/tournaments").POST(chain.ThenFunc(handlePOSTTournaments))

	router.Route("/js/:file").GET(http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))
	router.Route("/css/:file").GET(http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))

	err = http.ListenAndServe(":5010", router)
	if err != nil {
		fmt.Println(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	tData := TemplateData{}
	err = tData.getTournamentsData()
	if err != nil {
		log.Println(err)
		return
	}
	err = tData.getPlayersData()
	if err != nil {
		log.Println(err)
		return
	}
	err = templ.ExecuteTemplate(w, "index.html", tData)
	if err != nil {
		fmt.Println(err)
	}
}

func handleGETTournament(w http.ResponseWriter, r *http.Request) {
	tData := TemplateData{}
	id := r.FormValue("id")
	err := tData.getTournamentsData()
	if err != nil {
		log.Println(err)
		return
	}
	err = tData.getPlayersData()
	if err != nil {
		log.Println(err)
		return
	}
	err = tData.getTrnData(id)
	if err != nil {
		log.Println(err)
		return
	}
	err = templ.ExecuteTemplate(w, "tournament.html", tData)
	if err != nil {
		log.Println(err)
	}
}

func handleGETPlayer(w http.ResponseWriter, r *http.Request) {
	tData := TemplateData{}
	id := r.FormValue("id")
	err = tData.getTournamentsData()
	if err != nil {
		log.Println(err)
		return
	}
	err = tData.getPlayersData()
	if err != nil {
		log.Println(err)
		return
	}
	err = tData.getPlrData(id)
	if err != nil {
		log.Println(err)
		return
	}
	err = templ.ExecuteTemplate(w, "player.html", tData)
	if err != nil {
		log.Println(err)
	}
}

func handleGETPlayers(w http.ResponseWriter, r *http.Request) {
	tData := TemplateData{}
	err = tData.getTournamentsData()
	if err != nil {
		log.Println(err)
		return
	}
	err = tData.getPlayersData()
	if err != nil {
		log.Println(err)
		return
	}
	err = templ.ExecuteTemplate(w, "players.html", tData)
	if err != nil {
		log.Println(err)
	}
}

func handleGETTournaments(w http.ResponseWriter, r *http.Request) {
	tData := TemplateData{}
	err = tData.getTournamentsData()
	if err != nil {
		log.Println(err)
		return
	}
	err = tData.getPlayersData()
	if err != nil {
		log.Println(err)
		return
	}
	err = templ.ExecuteTemplate(w, "tournaments.html", tData)
	if err != nil {
		log.Println(err)
	}
}

func handlePOSTTournament(w http.ResponseWriter, r *http.Request) {
	action := r.FormValue("action")
	fmt.Println(action)
	id := r.URL.Query()["id"]
	name := r.FormValue("new-name")
	switch action {
	case "change-name":
		cmd := map[string]string{"@action": action, "name": name}
		bs, err := json.Marshal(cmd)
		if err != nil {
			log.Println(err)
			return
		}
		body := bytes.NewBuffer(bs)
		resp, err := http.Post(fmt.Sprintf("http://127.0.0.1:5000/api/tournaments/%s", id[0]), "", body)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(resp.Status)
		http.Redirect(w, r, fmt.Sprintf("/tournament?id=%s", id[0]), http.StatusSeeOther)
	}
}

func handlePOSTTournaments(w http.ResponseWriter, r *http.Request) {
	action := map[string]string{"@action": "create"}
	bs, err := json.Marshal(action)
	if err != nil {
		log.Println(err)
		return
	}
	body := bytes.NewBuffer(bs)
	resp, err := http.Post("http://127.0.0.1:5000/api/tournaments/", "", body)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(resp.Status)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// func postActionChangeName() error {
// 	action := map[string]string{"@action": "change-name"}
// 	bs, err := json.Marshal(action)
// 	if err != nil {
// 		return err
// 	}
// 	body := bytes.NewBuffer(bs)
// 	resp, err := http.Post("http://127.0.0.1:5000/api/tournament/")
// }

func (td *TemplateData) getTournamentsData() error {
	resp, err := http.Get("http://127.0.0.1:5000/api/tournaments/")
	if err != nil {
		return err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	err = json.Unmarshal(bs, &td.TournamentsData)
	if err != nil {
		return err
	}
	return nil
}

func (td *TemplateData) getPlayersData() error {
	resp, err := http.Get("http://127.0.0.1:5000/api/players/")
	if err != nil {
		return err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	err = json.Unmarshal(bs, &td.PlayersData)
	if err != nil {
		return err
	}
	return nil
}

func (td *TemplateData) getTrnData(id string) error {
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:5000/api/tournaments/%s", id))
	if err != nil {
		return err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	err = json.Unmarshal(bs, &td.Data)
	if err != nil {
		return err
	}
	return nil
}

func (td *TemplateData) getPlrData(id string) error {
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:5000/api/players/%s", id))
	if err != nil {
		return err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	err = json.Unmarshal(bs, &td.Data)
	if err != nil {
		return err
	}
	return nil
}
