package tournaments

import (
	"fmt"
	"net/http"
	"strings"
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
	err = plr.Create(ID, tID, "player", creds.Mail, string(hashedPassword))
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
		fmt.Println(storedCreds.Password)
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

func (s *Server) authentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isHTMLReq := strings.Contains(r.Header.Get("Accept"), "text/html")
		c, err := r.Cookie("token")
		if err != nil {
			handleError(w, http.StatusUnauthorized, err, isHTMLReq)
			return
		}

		tknStr := c.Value
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			handleError(w, http.StatusUnauthorized, err, isHTMLReq)
			return
		}
		if !tkn.Valid {
			handleError(w, http.StatusUnauthorized, fmt.Errorf("Token invalid"), isHTMLReq)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func (s *Server) refreshToken(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			h.ServeHTTP(w, r)
			return
		}
		tknStr := c.Value
		claims := &Claims{}
		_, err = jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			h.ServeHTTP(w, r)
			return
		}
		fmt.Println(time.Until(time.Unix(claims.ExpiresAt, 0)))
		//if token expires in >5min, skip refreshing
		if time.Until(time.Unix(claims.ExpiresAt, 0)) > 5*time.Minute {
			h.ServeHTTP(w, r)
			return
		}
		expirationTime := time.Now().Add(15 * time.Minute)
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
			Path:    "/api",
			Expires: expirationTime,
		})
		fmt.Println("refreshed jwt")
		h.ServeHTTP(w, r)
	})
}
func (s *Server) authenticate(r *http.Request) (PlayerID, error) {
	c, err := r.Cookie("token")
	if err != nil {
		return "", fmt.Errorf("Unauthorized")
	}

	tknStr := c.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("Unauthorized")
	}
	if !tkn.Valid {
		return "", fmt.Errorf("token invalid")
	}
	return claims.ID, nil
}

func (s *Server) checkAdminPermissions(accID PlayerID, errMsg string) error {
	acc, err := s.p.FindPlayerByID(accID)
	if err != nil {
		return err
	}
	if acc.Role != "admin" {
		return fmt.Errorf(errMsg)
	}
	return nil
}

func (s *Server) checkOrganizerPermissions(accID PlayerID, errMsg string) error {
	acc, err := s.p.FindPlayerByID(accID)
	if err != nil {
		return err
	}
	if acc.Role != "admin" && acc.Role != "organizer" {
		return fmt.Errorf(errMsg)
	}
	return nil
}

func (s *Server) checkPlayerPermissions(accID, pID PlayerID) error {
	acc, err := s.p.FindPlayerByID(accID)
	if err != nil {
		return err
	}
	if acc.Role != "admin" && acc.ID != pID {
		return fmt.Errorf("Unable to edit Player: Insufficient Permissions")
	}
	return nil
}

func (s *Server) refresh(w http.ResponseWriter, r *http.Request) {
	pID, err := s.authenticate(r)
	fmt.Println(pID)
	if err != nil {
		http.Redirect(w, r, "/api/signin", http.StatusSeeOther)
		return
	}
	claims := &Claims{}
	expirationTime := time.Now().Add(15 * time.Minute)
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
