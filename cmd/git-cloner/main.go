package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"

	git "github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
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

func cloneRepo(token string) {

	Info("git clone %s %s", giturl, directory)

	_, err := git.PlainClone(directory, false, &git.CloneOptions{
		Auth: &githttp.BasicAuth{
			Username: "girishg4t",
			Password: token,
		},
		URL:      giturl,
		Progress: os.Stdout,
	})
	CheckIfError(err)
}

func view(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, form, "", "", "", "", "", "", "")
}

func clearSession(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionStoreKey)
	if err != nil {
		fmt.Fprintln(w, "aborted")
		return
	}

	session.Options.MaxAge = -1

	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}

func callback(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionStoreKey)
	if err != nil {
		fmt.Fprintln(w, "aborted")
		return
	}

	if r.URL.Query().Get("state") != session.Values["state"] {
		fmt.Fprintln(w, "no state match; possible csrf OR cookies not enabled")
		return
	}

	token, err := oauthCfg.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
	if err != nil {
		fmt.Fprintln(w, "there was an issue getting your token")
		return
	}

	if !token.Valid() {
		fmt.Fprintln(w, "retreived invalid token")
		return
	}
	session.Values["githubAccessToken"] = token
	session.Save(r, w)

	cloneRepo(token.AccessToken)
	fmt.Fprintf(w, form, "", "", "", "", "", "", "Done")

}

func cloneHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionStoreKey)
	setUIValues(r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	if accessToken, ok := session.Values["githubAccessToken"].(*oauth2.Token); ok {
		cloneRepo(accessToken.AccessToken)
		fmt.Fprintf(w, form, "", "", "", "", "", "", "Done")
	} else {

		b := make([]byte, 16)
		rand.Read(b)

		state := base64.URLEncoding.EncodeToString(b)

		session, _ := store.Get(r, sessionStoreKey)
		session.Values["state"] = state
		session.Save(r, w)

		url := oauthCfg.AuthCodeURL(state)
		http.Redirect(w, r, url, 302)
	}
}

func init() {
	gob.Register(&oauth2.Token{})
}

func setUIValues(r *http.Request) {
	giturl = r.FormValue("giturl")
	directory = r.FormValue("directory")
}
