package tournaments

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cognicraft/uuid"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")

const adminPw = "admin"

type Credentials struct {
	Mail     string `json:"mail"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Mail string   `json:"mail"`
	ID   PlayerID `json:"id"`
	jwt.StandardClaims
}

func (s *Server) handleGETSignUp(w http.ResponseWriter, r *http.Request) {
	err = templ.ExecuteTemplate(w, "signUp.html", nil)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, true)
	}
}

func (s *Server) handleGETSignIn(w http.ResponseWriter, r *http.Request) {
	err = templ.ExecuteTemplate(w, "signIn.html", nil)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, true)
	}
}

func (s *Server) handleSignUp(w http.ResponseWriter, r *http.Request) {
	creds := &Credentials{
		Mail:     r.FormValue("mail"),
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, false)
		return
	}
	ok, err := s.p.IsMailAvailable(creds.Mail)
	if !ok {
		handleError(w, http.StatusInternalServerError, err, false)
		return
	}
	plr := NewPlayer(s)
	ID := PlayerID(uuid.MakeV4())
	tID := TrackerID(uuid.MakeV4())
	err = plr.Create(ID, tID, creds.Mail, string(hashedPassword), "player")
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, false)
		return
	}
	ok, err = s.p.IsPlayerNameAvailable(creds.Username)
	if !ok {
		handleError(w, http.StatusInternalServerError, err, false)
		return
	}
	err = plr.ChangeName(creds.Username)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, false)
		return
	}
	err = plr.Save(s.es, nil)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, false)
		return
	}
	http.Redirect(w, r, "/api/signin", http.StatusSeeOther)
}

func (s *Server) handleSignIn(w http.ResponseWriter, r *http.Request) {
	creds := &Credentials{
		Mail:     r.FormValue("mail"),
		Password: r.FormValue("password"),
	}
	storedCreds := Credentials{}
	var pID PlayerID
	storedCreds, pID, err = s.p.FindCredentialsByMail(creds.Mail)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, false)
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		handleError(w, http.StatusUnauthorized, err, false)
		return
	}

	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Mail: creds.Mail,
		ID:   pID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, false)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	http.Redirect(w, r, "/api/tournaments/", http.StatusSeeOther)
}

func (s *Server) authenticate(w http.ResponseWriter, r *http.Request) error {
	c, err := r.Cookie("token")
	if err != nil {
		return fmt.Errorf("Unauthorized")
	}

	tknStr := c.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	fmt.Println(tkn.Claims)
	if err != nil {
		return fmt.Errorf("Unauthorized")
	}
	if !tkn.Valid {
		return fmt.Errorf("token invalid")
	}
	return nil
}

func (s *Server) refresh(w http.ResponseWriter, r *http.Request) {
	err = s.authenticate(w, r)
	if err != nil {
		http.Redirect(w, r, "/api/signin", http.StatusSeeOther)
		return
	}
	claims := &Claims{}
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, true)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}
