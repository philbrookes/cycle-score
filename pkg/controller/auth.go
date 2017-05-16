package controller

import (
	"github.com/gorilla/mux"
	"fmt"
	"net/http"
	"../config"
	"github.com/strava/go.strava"
	"os"
	"encoding/json"
	"time"
	"strconv"
)

func ConfigureAuth(router *mux.Router) {
	router.HandleFunc("/check", checkAuth)
	router.HandleFunc("/url", getOAuthUrl)
	router.HandleFunc("/url/{callback}", getOAuthUrlWithCallback)
	configureOAuth(router)
}

type authUrl struct {
	Url string `json:"url"`
}

type authCheckResponse struct {
	Valid bool `json:"valid"`
}

func checkAuth (w http.ResponseWriter, r *http.Request) {
	res := authCheckResponse{Valid: true}
	_, err := r.Cookie("strava_token")
	if err == http.ErrNoCookie {
		res.Valid = false
	}

	w.Header().Add("Content-Type", "Application/JSON")
	json.NewEncoder(w).Encode(res)
}

func getOAuthUrlWithCallback (w http.ResponseWriter,r *http.Request) {
	params := mux.Vars(r)
	strava.ClientId = config.GetConfig().GetClientId()
	strava.ClientSecret = config.GetConfig().GetClientSecret()

	authenticator := &strava.OAuthAuthenticator{
		CallbackURL:            params["callback"],
		RequestClientGenerator: nil,
	}
	authResponse := authUrl{
		Url: authenticator.AuthorizationURL("", strava.Permissions.Public, false),
	}
	w.Header().Add("Content-Type", "Application/JSON")
	json.NewEncoder(w).Encode(authResponse)
}

func getOAuthUrl (w http.ResponseWriter,r *http.Request) {
	strava.ClientId = config.GetConfig().GetClientId()
	strava.ClientSecret = config.GetConfig().GetClientSecret()

	authenticator := &strava.OAuthAuthenticator{
		CallbackURL:            config.GetConfig().GetOAuthCallbackUrl(),
		RequestClientGenerator: nil,
	}
	authResponse := authUrl{
		Url: authenticator.AuthorizationURL("", strava.Permissions.Public, false),
	}
	w.Header().Add("Content-Type", "Application/JSON")
	json.NewEncoder(w).Encode(authResponse)

}

func configureOAuth(router *mux.Router) {
	strava.ClientId = config.GetConfig().GetClientId()
	strava.ClientSecret = config.GetConfig().GetClientSecret()

	authenticator := &strava.OAuthAuthenticator{
		CallbackURL:            config.GetConfig().GetOAuthCallbackUrl(),
		RequestClientGenerator: nil,
	}
	_, err := authenticator.CallbackPath()
	if err != nil {
		// possibly that the callback url set above is invalid
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	router.HandleFunc("/callback", authenticator.HandlerFunc(oAuthSuccess, oAuthFailure))
}


func oAuthSuccess(auth *strava.AuthorizationResponse, w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "Application/JSON")
	authCookie := &http.Cookie{
		Name: "strava_token",
		Value: auth.AccessToken,
		Expires: time.Now().AddDate(0, 0, config.GetConfig().GetCookieLifetime()),
		Path: "/",
	}
	athleteCookie := &http.Cookie{
		Name: "strava_athlete_id",
		Value: strconv.FormatInt(auth.Athlete.Id, 10),
		Expires: time.Now().AddDate(0, 0, config.GetConfig().GetCookieLifetime()),
		Path: "/",
	}
	http.SetCookie(w, authCookie)
	http.SetCookie(w, athleteCookie)
	w.Header().Add("Location", "http://localhost:8080")
	w.WriteHeader(http.StatusFound)
}


func oAuthFailure(err error, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Authorization Failure:\n")

	// some standard error checking
	if err == strava.OAuthAuthorizationDeniedErr {
		fmt.Fprint(w, "The user clicked the 'Do not Authorize' button on the previous page.\n")
		fmt.Fprint(w, "This is the main error your application should handle.")
	} else if err == strava.OAuthInvalidCredentialsErr {
		fmt.Fprint(w, "You provided an incorrect client_id or client_secret.\nDid you remember to set them at the begininng of this file?")
	} else if err == strava.OAuthInvalidCodeErr {
		fmt.Fprint(w, "The temporary token was not recognized, this shouldn't happen normally")
	} else if err == strava.OAuthServerErr {
		fmt.Fprint(w, "There was some sort of server error, try again to see if the problem continues")
	} else {
		fmt.Fprint(w, err)
	}
}