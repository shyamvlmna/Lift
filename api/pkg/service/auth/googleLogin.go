package auth

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	authConfig = &oauth2.Config{
		ClientID:     "662778233746-oeblr3vkk1om82nmjce90lqlac1p7fvq.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8080/user/googleCallback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	randomState = "randomstate"
)

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := authConfig.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {

	state := r.URL.Query()["state"][0]
	if state != "randomstate" {
		fmt.Fprintln(w, "states don't match")
		return
	}
	// if r.FormValue("state") != randomState {
	// 	fmt.Println("not a valid state")
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }
	code := r.URL.Query()["code"][0]
	tok, err := authConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Fprintln(w, "code token exange failed")
	}
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tok.AccessToken)
	if err != nil {
		fmt.Fprintln(w, "data fetch failed")
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, "json parsing failed")
	}
	fmt.Fprintln(w, string(content))
}
