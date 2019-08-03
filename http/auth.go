package http

import (
	"encoding/json"
	"fmt"
	"github.com/dghubble/gologin"
	fbLogin "github.com/dghubble/gologin/facebook"
	"github.com/goincremental/negroni-sessions"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	. "jng.dev/tailor/core/core"
	"jng.dev/tailor/core/sqlx"
	"log"
	"math/rand"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, status int, m interface{}) {
	w.WriteHeader(status)
	resp, _ := json.Marshal(m)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(resp))
	log.Print(m)
}

func UserIndexHandler(w http.ResponseWriter, r *http.Request) {
	db := sqlx.Database()
	defer sqlx.Close(db)
	var u []User
	err := db.Select(&u, "SELECT name, email FROM core.user")
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, ErrorMessage{"Error retrieving users"})
		return
	}
	WriteResponse(w, http.StatusOK, UserIndexResponse{u})
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	s := sessions.GetSession(r)
	u := User{}
	c := s.Get("user-details")
	if c == nil {
		WriteResponse(w, http.StatusNotFound, ErrorMessage{"No user details for registration"})
		return
	}
	err := json.Unmarshal(c.([]byte), &u)
	if err != nil {
		WriteResponse(w, http.StatusBadRequest, ErrorMessage{"Could not read user details for registration"})
		return
	}
	switch r.Method {
	case http.MethodGet:
		GetRegistrationHandler(w, r, u)
		return
	case http.MethodPut:
		RegisterHandler(w, r, &u)
		return
	}
	http.NotFoundHandler().ServeHTTP(w, r)
}

func GetRegistrationHandler(w http.ResponseWriter, r *http.Request, u User) {
	retrievedUser := User{}
	success := sqlx.GetUser(&retrievedUser, u.Email)
	if success {
		u = retrievedUser
	}
	ri := RegistrationInfoResponse{
		User: u,
	}
	WriteResponse(w, http.StatusOK, ri)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request, u *User) {
	// should verify details with provider here
	// should make sure user doesn't exists here
	u.ID = uuid.New().String()
	u.VerificationCode = fmt.Sprintf("%03d %03d", rand.Intn(1000), rand.Intn(1000))
	db := sqlx.Database()
	defer sqlx.Close(db)
	result, err := db.NamedExec(
		"INSERT INTO core.user (id, provider, provider_id, email, name, image_url, verification_code) VALUES (:id, :provider, :provider_id, :email, :name, :image_url, :verification_code)",
		u)
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, ErrorMessage{"Could not store user details for registration"})
		return
	}
	i, err := result.RowsAffected()
	if err != nil || i != 1 {
		WriteResponse(w, http.StatusInternalServerError, ErrorMessage{"Error on store user details for registration"})
		return
	}
	rr := RegistrationResponse{
		VerificationCode: u.VerificationCode,
	}
	WriteResponse(w, http.StatusOK, rr)
}

func LoginFacebook(w http.ResponseWriter, r *http.Request) {
	oauth2Config := &oauth2.Config{
		ClientID:     "537772546605649",
		ClientSecret: "578f34dd12ef97c2dcfc29c702fbbe32",
		RedirectURL:  "http://localhost:8080/login/auth?provider=facebook",
		Endpoint:     facebook.Endpoint,
		Scopes:       []string{"email"},
	}
	fbErrorHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteResponse(w, http.StatusBadGateway, ErrorMessage{"Could not interface with Facebook"})
	})
	fbLogin.StateHandler(gologin.DebugOnlyCookieConfig, fbLogin.LoginHandler(oauth2Config, fbErrorHandler)).ServeHTTP(w, r)
}

func LoginAuthHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Query().Get("provider") {
	case "facebook":
		LoginAuthFacebook(w, r)
		break
	default:
		WriteResponse(w, http.StatusBadRequest, ErrorMessage{"Unknown provider"})
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	u := User{}
	email := r.URL.Query().Get("email")
	success := sqlx.GetUser(&u, email)
	if !success {
		WriteResponse(w, http.StatusNotFound, ErrorMessage{"User not found"})
		return
	}
	WriteResponse(w, http.StatusOK, u)
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetUserHandler(w, r)
		break
	}
}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Query().Get("provider") {
	case "facebook":
		LoginFacebook(w, r)
		return
	default:
		WriteResponse(w, http.StatusBadRequest, ErrorMessage{"Unknown provider"})
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	s := sessions.GetSession(r)
	s.Clear()
	w.WriteHeader(http.StatusNoContent)
}

func LoginAuthFacebook(w http.ResponseWriter, r *http.Request) {
	oauth2Config := &oauth2.Config{
		ClientID:     "537772546605649",
		ClientSecret: "578f34dd12ef97c2dcfc29c702fbbe32",
		RedirectURL:  "http://localhost:8080/login/auth?provider=facebook",
		Endpoint:     facebook.Endpoint,
		Scopes:       []string{"email"},
	}
	fbErrorHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// this should be a page not a json
		WriteResponse(w, http.StatusBadRequest, ErrorMessage{"Could not verify with Facebook"})
	})
	fbSuccessHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		fbUser, err := fbLogin.UserFromContext(ctx)
		if err != nil {
			WriteResponse(w, http.StatusBadGateway, ErrorMessage{"Could not load user with Facebook"})
			return
		}
		u := User{
			Provider:   "facebook",
			ProviderID: fbUser.ID,
			Email:      fbUser.Email,
			Name:       fbUser.Name,
			ImageURL:   GetImageURL("facebook", fbUser.ID),
		}
		s := sessions.GetSession(r)
		udt, _ := json.Marshal(u)
		s.Set("user-details", udt)
	})
	fbLogin.StateHandler(gologin.DebugOnlyCookieConfig, fbLogin.CallbackHandler(oauth2Config, fbSuccessHandler, fbErrorHandler)).ServeHTTP(w, r)
}
