package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

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
