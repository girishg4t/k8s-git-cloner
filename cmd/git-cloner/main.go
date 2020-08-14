package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	git "github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/joho/godotenv"
)

const defaultAddr = ":8081"

var giturl = ""
var directory = ""

func main() {
	err := godotenv.Load()
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

	values := map[string]string{
		"client_id":     os.Getenv("client_id"),
		"client_secret": os.Getenv("client_secret"),
		"code":          header["code"][0],
		"redirect_uri":  "http://" + r.Host + "/token",
	}

	jsonValue, _ := json.Marshal(values)

	result, _ := http.Post("https://github.com/login/oauth/access_token",
		"application/json", bytes.NewBuffer(jsonValue))
	defer result.Body.Close()
	body, _ := ioutil.ReadAll(result.Body)

	tokenURL := "http://" + r.Host + "/token?" + string(body)
	http.Redirect(w, r, tokenURL, http.StatusTemporaryRedirect)
}

func handler(w http.ResponseWriter, r *http.Request) {
	var oauth2Endpoint = "https://github.com/login/oauth/authorize"
	req, _ := http.NewRequest("GET", oauth2Endpoint, nil)
	req.Header.Add("Accept", "application/json")
	q := req.URL.Query()
	q.Add("client_id", os.Getenv("client_id"))
	q.Add("redirect_uri", "http://"+r.Host+"/callback")
	q.Add("scope", "public_repo")
	q.Add("include_granted_scopes", "true")

	req.URL.RawQuery = q.Encode()
	giturl = r.FormValue("giturl")
	directory = r.FormValue("directory")
	http.Redirect(w, r, req.URL.String(), http.StatusTemporaryRedirect)
}
