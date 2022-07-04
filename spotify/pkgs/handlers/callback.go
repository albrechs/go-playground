package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"spotify/pkgs/structs"

	"github.com/google/go-querystring/query"
)

func Callback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	code := q.Get("code")
	state := q.Get("state")
	cookie, err := r.Cookie(statekey)
	var storedstate string
	if err != nil {
		storedstate = ""
	} else {
		storedstate = string(cookie.Value)
	}
	if len(state) == 0 || state != storedstate {
		http.Redirect(w, r, "/#?error=state_mismatch", http.StatusPermanentRedirect)
	} else {
		c := &http.Cookie{
			Name:   statekey,
			Value:  "",
			MaxAge: -1,
		}
		r.AddCookie(c)

		client := &http.Client{}
		formdata := url.Values{
			"code":         {code},
			"redirect_uri": {redirecturi},
			"grant_type":   {"authorization_code"},
		}
		req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", bytes.NewBufferString(formdata.Encode()))
		if err != nil {
			log.Printf("error: %s", err)
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		authstring := clientid + ":" + clientsecret
		authencoded := base64.StdEncoding.EncodeToString([]byte(authstring))
		req.Header.Set("Authorization", "Basic "+authencoded)

		res, err := client.Do(req)
		if err != nil {
			log.Printf("got error %s", err)
			return
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		var result structs.TokenResponse
		if err := json.Unmarshal(body, &result); err != nil {
			log.Println("unable to unmarshall tokenresponse")
		}

		/* mereq, meerr := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
		if meerr != nil {
			log.Printf("error: %s", meerr)
			return
		}
		mereq.Header.Set("Authorization", "Bearer "+result.AccessToken)
		meres, meerr := client.Do(mereq)
		if meerr != nil {
			log.Printf("got error %s", err)
			return
		}
		defer meres.Body.Close()
		mebody, meerr := ioutil.ReadAll(meres.Body)
		if meerr != nil {
			log.Printf("got error %s", err)
			return
		} */

		os.Setenv("SPOTIFY_API_TOKEN", result.AccessToken)
		callbackquery := structs.CallbackRedirectQuery{
			AccessToken:  result.AccessToken,
			RefreshToken: result.RefreshToken,
		}
		val, err := query.Values(callbackquery)
		if err != nil {
			log.Println("unable to parse login query values")
			return
		}
		querystring := string(val.Encode())
		http.Redirect(w, r, "/#"+querystring, http.StatusFound)
	}
}
