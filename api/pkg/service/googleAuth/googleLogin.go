package googleauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/shayamvlmna/cab-booking-app/pkg/models"
	"github.com/shayamvlmna/cab-booking-app/pkg/service/user"
)

var (
	ID         = os.Getenv("OAUTHCLIENTID")
	Key        = os.Getenv("OAUTHCLIENTSECRET")
	authConfig = &oauth2.Config{
		ClientID:     ID,
		ClientSecret: Key,
		RedirectURL:  "http://localhost:8080/user/googleCallback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	randomState = "randomstate"
)

type AuthContent struct {
	ID            string `json:"id"`
	Firstname     string `json:"given_name"`
	Lastname      string `json:"family_name"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := authConfig.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	state := r.URL.Query()["state"][0]
	if state != "randomstate" {

		response := &models.Response{
			ResponseStatus:  "fail",
			ResponseMessage: "states don't match",
			ResponseData:    nil,
		}
		err := json.NewEncoder(w).Encode(&response)
		if err != nil {
			return
		}
		return
	}

	code := r.URL.Query()["code"][0]
	tok, err := authConfig.Exchange(context.Background(), code)
	if err != nil {
		response := &models.Response{
			ResponseStatus:  "fail",
			ResponseMessage: "code token exange failed",
			ResponseData:    nil,
		}
		err := json.NewEncoder(w).Encode(&response)
		if err != nil {
			return
		}
		return
	}
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tok.AccessToken)
	if err != nil {
		response := &models.Response{
			ResponseStatus:  "fail",
			ResponseMessage: "data fetch failed",
			ResponseData:    nil,
		}
		err := json.NewEncoder(w).Encode(&response)
		if err != nil {
			return
		}
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response := &models.Response{
			ResponseStatus:  "fail",
			ResponseMessage: "json parsing failed",
			ResponseData:    nil,
		}
		err := json.NewEncoder(w).Encode(&response)
		if err != nil {
			return
		}
		return
	}
	authContent := &AuthContent{}
	if err = json.Unmarshal(content, &authContent); err != nil {
		response := &models.Response{
			ResponseStatus:  "fail",
			ResponseMessage: "unmarshal failed",
			ResponseData:    nil,
		}
		err := json.NewEncoder(w).Encode(&response)
		if err != nil {
			return
		}
		return
	}

	newUser := &models.User{}

	newUser.Firstname = authContent.Firstname
	newUser.Lastname = authContent.Lastname
	newUser.Email = authContent.Email
	newUser.Picture = authContent.Picture

	user.GoogleAuthUser(newUser)

}
