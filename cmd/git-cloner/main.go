package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
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
	scopes   = []string{"repo"}
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
		Info("Not able to read the configuration")
		return
	}
	addr := defaultAddr
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("server starting to listen on %s", addr)
	http.HandleFunc("/clone", handler)
	http.HandleFunc("/callback", callback)
	http.HandleFunc("/token", handleTokenCallback)
	http.HandleFunc("/", view)
	if err := http.ListenAndServe(addr, nil); err != nil {
		Info("server listen error: %+v", err)
	}
}

func handleTokenCallback(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(r.URL.String())
	fmt.Println(u.RawQuery)
	header, _ := url.ParseQuery(u.RawQuery)
	fmt.Println(header["access_token"][0])

	Info("git clone %s %s", giturl, directory)

	_, err := git.PlainClone(directory, false, &git.CloneOptions{
		Auth: &githttp.BasicAuth{
			Username: "girishg4t",
			Password: header["access_token"][0],
		},
		URL:      giturl,
		Progress: os.Stdout,
	})
	CheckIfError(err)
	fmt.Fprintf(w, form, "", "", "", "", "", "", "Done")
}
func view(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, form, "", "", "", "", "", "", "")
}

func callback(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)
	u, _ := url.Parse(r.URL.String())
	fmt.Println(u.RawQuery)
	header, _ := url.ParseQuery(u.RawQuery)
	fmt.Println(header["code"][0])

	// values := map[string]string{
	// 	"client_id":     os.Getenv("client_id"),
	// 	"client_secret": os.Getenv("client_secret"),
	// 	"code":          header["code"][0],
	// 	"redirect_uri":  "http://" + r.Host + "/token",
	// }

	// jsonValue, _ := json.Marshal(values)

	token, err := oauthCfg.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
	if err != nil {
		fmt.Fprintln(w, "there was an issue getting your token")
		return
	}

	if !token.Valid() {
		fmt.Fprintln(w, "retreived invalid token")
		return
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	// var oauth2Endpoint = "https://github.com/login/oauth/authorize"
	// req, _ := http.NewRequest("GET", oauth2Endpoint, nil)
	// req.Header.Add("Accept", "application/json")
	// q := req.URL.Query()
	// q.Add("client_id", os.Getenv("client_id"))
	// q.Add("redirect_uri", "http://"+r.Host+"/callback")
	// q.Add("scope", "public_repo")
	// q.Add("include_granted_scopes", "true")

	// req.URL.RawQuery = q.Encode()
	giturl = r.FormValue("giturl")
	directory = r.FormValue("directory")
	// http.Redirect(w, r, req.URL.String(), http.StatusFound)

	b := make([]byte, 16)
	rand.Read(b)

	state := base64.URLEncoding.EncodeToString(b)

	session, _ := store.Get(r, sessionStoreKey)
	session.Values["state"] = state
	session.Save(r, w)

	url := oauthCfg.AuthCodeURL("sdgasdg")
	http.Redirect(w, r, url, 302)
}
