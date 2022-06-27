package auth

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/shayamvlmna/cab-booking-app/pkg/models"
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

type AuthContent struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
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
			Token:           "",
		}
		json.NewEncoder(w).Encode(&response)
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
		response := &models.Response{
			ResponseStatus:  "fail",
			ResponseMessage: "code token exange failed",
			ResponseData:    nil,
		}
		json.NewEncoder(w).Encode(&response)
		return
	}
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tok.AccessToken)
	if err != nil {
		response := &models.Response{
			ResponseStatus:  "fail",
			ResponseMessage: "data fetch failed",
			ResponseData:    nil,
		}
		json.NewEncoder(w).Encode(&response)
		return
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response := &models.Response{
			ResponseStatus:  "fail",
			ResponseMessage: "json parsing failed",
			ResponseData:    nil,
		}
		json.NewEncoder(w).Encode(&response)
		return
	}
	authContent := &AuthContent{}
	if err = json.Unmarshal(content, &authContent); err != nil {
		response := &models.Response{
			ResponseStatus:  "fail",
			ResponseMessage: "unmarshal failed",
			ResponseData:    nil,
		}
		json.NewEncoder(w).Encode(&response)
		return
	}

	json.NewEncoder(w).Encode(&authContent)
	// fmt.Fprintln(w, string(content))
	// {
	// 	"id": "109429758760150763543",
	// 	"email": "shyamvlmna@gmail.com",
	// 	"verified_email": true,
	// 	"name": "Shyamjith P Vilamana",
	// 	"given_name": "Shyamjith",
	// 	"family_name": "P Vilamana",
	// 	"picture": "https://lh3.googleusercontent.com/a-/AOh14Gj4L240leqI64MfmshtoQsqLv_vm0RTPoZ4Z9yCHg=s96-c",
	// 	"locale": "en"
	//   }
	// 	authContent := &AuthContent{}

	// 	json.Unmarshal(content, &authContent)

	// 	log.Println(authContent)
	// 	json.NewEncoder(w).Encode(&authContent)

	// fmt.Println(authContent.Email)
}
