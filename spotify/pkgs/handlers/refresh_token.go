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

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	refreshtoken := q.Get("refresh_token")

	client := &http.Client{}
	formdata := url.Values{
		"refresh_token": {refreshtoken},
		"grant_type":    {"refresh_token"},
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
	refreshbody, refresherr := ioutil.ReadAll(res.Body)
	if refresherr != nil {
		log.Printf("got error %s", err)
		return
	}
	var result structs.TokenResponse
	if err := json.Unmarshal(refreshbody, &result); err != nil {
		log.Println("unable to unmarshall tokenresponse")
	}
	log.Println(result.AccessToken)
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
