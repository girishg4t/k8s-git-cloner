package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

const sessionStoreKey = "sess"

const (
	defaultConfigFile = "config.json"

	githubAuthorizeUrl = "https://github.com/login/oauth/authorize"
	githubTokenUrl     = "https://github.com/login/oauth/access_token"
	redirectUrl        = ""
)

var (
	oauthCfg *oauth2.Config
	scopes   = []string{"public_repo"}
	store    *sessions.CookieStore
)

const defaultAddr = ":8081"

var giturl = ""
var directory = ""

func main() {
	err := godotenv.Load()
	oauthCfg = &oauth2.Config{
		ClientID:     os.Getenv("client_id"),
		ClientSecret: os.Getenv("client_secret"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  githubAuthorizeUrl,
			TokenURL: githubTokenUrl,
		},
		RedirectURL: redirectUrl,
		Scopes:      scopes,
	}
	store = sessions.NewCookieStore([]byte(os.Getenv("serverSecret")))

	if err != nil {
		Info("Not able to read the configuration %s", err)
		return
	}
	addr := defaultAddr
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("server starting to listen on %s", addr)
	http.HandleFunc("/clone", cloneHandler)
	http.HandleFunc("/clear", clearSession)
	http.HandleFunc("/callback", callback)
	http.HandleFunc("/", view)
	if err := http.ListenAndServe(addr, nil); err != nil {
		Info("server listen error: %+v", err)
	}
}

func view(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, form, "", "", "", "", "", "", "")
}

func init() {
	gob.Register(&oauth2.Token{})
}

func setUIValues(r *http.Request) {
	giturl = r.FormValue("giturl")
	directory = r.FormValue("directory")
}
