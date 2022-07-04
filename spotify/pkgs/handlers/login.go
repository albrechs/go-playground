package handlers

import (
	"log"
	"net/http"
	"spotify/pkgs/helpers"
	"spotify/pkgs/structs"

	"github.com/google/go-querystring/query"
)

func Login(w http.ResponseWriter, r *http.Request) {
	state := helpers.GenerateRandomString(16)
	http.SetCookie(w, &http.Cookie{
		Name:   statekey,
		Value:  state,
		MaxAge: 3000,
	})
	scope := "user-read-private user-read-email user-top-read"

	opts := structs.LoginOpts{
		ResponseType: "code",
		ClientID:     clientid,
		Scope:        scope,
		RedirectURI:  redirecturi,
		State:        state,
	}
	val, err := query.Values(opts)
	if err != nil {
		log.Println("unable to parse login query values")
		return
	}
	querystring := string(val.Encode())
	http.Redirect(w, r, "https://accounts.spotify.com/authorize?"+querystring, http.StatusPermanentRedirect)
}
