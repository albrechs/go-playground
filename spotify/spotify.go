package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"

	"github.com/google/go-querystring/query"
)

var clientid = "58f7888547c24f058ece41fed973bf37"     // Your client id
var clientsecret = "8405db97ae294a9885f064855b5fae79" // Your secret
var redirecturi = "http://localhost:8888/callback"    // Your redirect uri
var statekey = "spotify_auth_state"
var templates = template.Must(template.ParseFiles("tmpl/top-tracks.html"))

type TopTracksPage struct {
	Name  string
	Items []TopTracksNode
}

func renderTopTracksTemplate(w http.ResponseWriter, p *TopTracksPage) {
	e := templates.ExecuteTemplate(w, "top-tracks.html", p)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}

func parsePercent(p float64) string {
	return fmt.Sprintf("%.0f", (p*100)) + "%"
}

func parseKey(k int) string {
	switch k {
	case 0:
		return "C"
	case 1:
		return "C♯"
	case 2:
		return "D"
	case 3:
		return "D♯"
	case 4:
		return "E"
	case 5:
		return "F"
	case 6:
		return "F♯"
	case 7:
		return "G"
	case 8:
		return "G♯"
	case 9:
		return "A"
	case 10:
		return "A♯"
	case 11:
		return "B"
	default:
		return "??"
	}
}

func parseTempo(t float64) string {
	return fmt.Sprintf("%.0f BPM", t)
}

func parseLoudness(l float64) string {
	return fmt.Sprintf("%.0f db", l)
}

type TopTracksResponse struct {
	Items []struct {
		Album struct {
			AlbumType string `json:"album_type"`
			Artists   []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				URI  string `json:"uri"`
			} `json:"artists"`
			AvailableMarkets []string `json:"available_markets"`
			ExternalUrls     struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href   string `json:"href"`
			ID     string `json:"id"`
			Images []struct {
				Height int    `json:"height"`
				URL    string `json:"url"`
				Width  int    `json:"width"`
			} `json:"images"`
			Name                 string `json:"name"`
			ReleaseDate          string `json:"release_date"`
			ReleaseDatePrecision string `json:"release_date_precision"`
			TotalTracks          int    `json:"total_tracks"`
			Type                 string `json:"type"`
			URI                  string `json:"uri"`
		} `json:"album"`
		Artists []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href string `json:"href"`
			ID   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
			URI  string `json:"uri"`
		} `json:"artists"`
		AvailableMarkets []string `json:"available_markets"`
		DiscNumber       int      `json:"disc_number"`
		DurationMs       int      `json:"duration_ms"`
		Explicit         bool     `json:"explicit"`
		ExternalIds      struct {
			Isrc string `json:"isrc"`
		} `json:"external_ids"`
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href        string `json:"href"`
		ID          string `json:"id"`
		IsLocal     bool   `json:"is_local"`
		Name        string `json:"name"`
		Popularity  int    `json:"popularity"`
		PreviewURL  string `json:"preview_url"`
		TrackNumber int    `json:"track_number"`
		Type        string `json:"type"`
		URI         string `json:"uri"`
	} `json:"items"`
	Total    int         `json:"total"`
	Limit    int         `json:"limit"`
	Offset   int         `json:"offset"`
	Href     string      `json:"href"`
	Previous interface{} `json:"previous"`
	Next     string      `json:"next"`
}

type TopTracksNode struct {
	ID           string
	Artist       string
	Title        string
	ImageURL     string
	Key          string
	Tempo        string
	Loudness     string
	Happiness    string
	Energy       string
	Danceability string
	Acousticness string
}

func topTracksHandler(w http.ResponseWriter, r *http.Request) {
	authheader := fmt.Sprintf("Bearer %s", os.Getenv("SPOTIFY_API_TOKEN"))

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/top/tracks?limit=50", nil)
	if err != nil {
		log.Printf("Got error %s", err)
		return
	}
	req.Header.Set("Authorization", authheader)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Got error %s", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Got error %s", err)
		return
	}
	var result TopTracksResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Print("Unable to unmarshal JSON response.")
	}
	toptracks := make([]TopTracksNode, 0)
	for _, rec := range result.Items {
		analysis, err := getTrackAnalysis(rec.ID)
		if err != nil {
			log.Printf("could not find track analyisis for" + rec.ID)
			log.Print(err)
		}
		toptracks = append(toptracks, TopTracksNode{
			ID:           rec.ID,
			Artist:       rec.Artists[0].Name,
			Title:        rec.Name,
			ImageURL:     rec.Album.Images[0].URL,
			Key:          parseKey(analysis.Key),
			Tempo:        parseTempo(analysis.Tempo),
			Loudness:     parseLoudness(analysis.Loudness),
			Happiness:    parsePercent(analysis.Valence),
			Energy:       parsePercent(analysis.Energy),
			Danceability: parsePercent(analysis.Danceability),
			Acousticness: parsePercent(analysis.Acousticness),
		})
	}
	log.Print(toptracks)
	data := &TopTracksPage{
		Name:  "Trevor",
		Items: toptracks,
	}

	renderTopTracksTemplate(w, data)
}

func generateRandomString(n int) string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

type LoginOpts struct {
	ResponseType string `url:"response_type"`
	ClientID     string `url:"client_id"`
	Scope        string `url:"scope"`
	RedirectURI  string `url:"redirect_uri"`
	State        string `url:"state"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	state := generateRandomString(16)
	http.SetCookie(w, &http.Cookie{
		Name:   statekey,
		Value:  state,
		MaxAge: 3000,
	})
	scope := "user-read-private user-read-email user-top-read"

	opts := LoginOpts{
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

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type MeResponse struct {
	Country         string `json:"country"`
	DisplayName     string `json:"display_name"`
	Email           string `json:"email"`
	ExplicitContent struct {
		FilterEnabled bool `json:"filter_enabled"`
		FilterLocked  bool `json:"filter_locked"`
	} `json:"explicit_content"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Followers struct {
		Href  interface{} `json:"href"`
		Total int         `json:"total"`
	} `json:"followers"`
	Href   string `json:"href"`
	ID     string `json:"id"`
	Images []struct {
		Height interface{} `json:"height"`
		URL    string      `json:"url"`
		Width  interface{} `json:"width"`
	} `json:"images"`
	Product string `json:"product"`
	Type    string `json:"type"`
	URI     string `json:"uri"`
}

type CallbackRedirectQuery struct {
	AccessToken  string `url:"access_token"`
	RefreshToken string `url:"refresh_token"`
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	code := q.Get("code")
	state := q.Get("state")
	cookie, err := r.Cookie(statekey)
	log.Print(cookie)
	var storedstate string
	if err != nil {
		storedstate = ""
	} else {
		storedstate = string(cookie.Value)
	}
	log.Println("state: " + state)
	log.Println("storedstate: " + storedstate)
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
		var result TokenResponse
		if err := json.Unmarshal(body, &result); err != nil {
			log.Println("unable to unmarshall tokenresponse")
		}
		log.Println(result.AccessToken)

		mereq, meerr := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
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
		}
		log.Print(string(mebody))

		os.Setenv("SPOTIFY_API_TOKEN", result.AccessToken)
		callbackquery := CallbackRedirectQuery{
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

type TokenRefreshResponse struct {
	AccessToken string `json:"access_token"`
}

func refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
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
	var result TokenResponse
	if err := json.Unmarshal(refreshbody, &result); err != nil {
		log.Println("unable to unmarshall tokenresponse")
	}
	log.Println(result.AccessToken)
	os.Setenv("SPOTIFY_API_TOKEN", result.AccessToken)
	callbackquery := CallbackRedirectQuery{
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

type TrackAnalysisResponse struct {
	Danceability     float64 `json:"danceability"`
	Energy           float64 `json:"energy"`
	Key              int     `json:"key"`
	Loudness         float64 `json:"loudness"`
	Mode             int     `json:"mode"`
	Speechiness      float64 `json:"speechiness"`
	Acousticness     float64 `json:"acousticness"`
	Instrumentalness float64 `json:"instrumentalness"`
	Liveness         float64 `json:"liveness"`
	Valence          float64 `json:"valence"`
	Tempo            float64 `json:"tempo"`
	Type             string  `json:"type"`
	ID               string  `json:"id"`
	URI              string  `json:"uri"`
	TrackHref        string  `json:"track_href"`
	AnalysisURL      string  `json:"analysis_url"`
	DurationMs       int     `json:"duration_ms"`
	TimeSignature    int     `json:"time_signature"`
}

func getTrackAnalysis(id string) (*TrackAnalysisResponse, error) {
	authheader := fmt.Sprintf("Bearer %s", os.Getenv("SPOTIFY_API_TOKEN"))

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/audio-features/"+id, nil)
	if err != nil {
		log.Printf("Got error %s", err)
		return nil, err
	}
	req.Header.Set("Authorization", authheader)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Got error %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Got error %s", err)
		return nil, err
	}
	//log.Print(string(body))
	var result TrackAnalysisResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Print("Unable to unmarshal JSON response.")
		log.Print(err)
	}
	return &result, nil
}

func main() {
	var port int
	flag.IntVar(&port, "port", 8888, "Listener port")
	flag.Parse()

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/callback", callbackHandler)
	http.HandleFunc("/refresh_token", refreshTokenHandler)
	http.HandleFunc("/top-tracks", topTracksHandler)
	//http.HandleFunc("/track/", trackHandler)

	var frontend fs.FS = os.DirFS("public")
	httpFS := http.FS(frontend)
	fileServer := http.FileServer(httpFS)
	http.Handle("/", fileServer)

	addr := fmt.Sprintf("localhost:%d", port)
	log.Printf("Serving app at http://%s", addr)
	log.Fatalln(http.ListenAndServe(addr, nil))
}
