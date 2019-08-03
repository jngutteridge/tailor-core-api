package http

import (
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
)

const appLocation = "public/index.html"
const port = ":8081"
const sessionSecret = "d234fjhkg234ddfgjhk6354fdghjk"

func Server() {
	n := negroni.Classic()

	store := cookiestore.New([]byte(sessionSecret))
	n.Use(sessions.Sessions("tailor-core", store))

	r := mux.NewRouter()
	r.HandleFunc("/", RootHandler)
	r.HandleFunc("/login", LoginHandler)
	r.HandleFunc("/login/auth", LoginAuthHandler)
	r.HandleFunc("/api/login/register", RegistrationHandler)
	r.HandleFunc("/api/logout", LogoutHandler).Methods(http.MethodDelete)
	r.HandleFunc("/api/access", RegistrationHandler)
	r.HandleFunc("/api/access/token", RegistrationHandler)
	r.HandleFunc("/api/text", TextHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/user/index", UserIndexHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/user", UserHandler)
	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	n.UseHandler(r)

	n.Run(port)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, appLocation)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	RootHandler(w, r)
}