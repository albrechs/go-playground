package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"spotify/pkgs/helpers"
	"spotify/pkgs/structs"
	"strconv"
)

func TopTracks(w http.ResponseWriter, r *http.Request) {
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
	var result structs.TopTracksResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Print("Unable to unmarshal JSON response.")
	}
	toptracks := make([]structs.TopTracksNode, 0)
	for _, rec := range result.Items {
		analysis, err := getTrackAnalysis(rec.ID)
		if err != nil {
			log.Printf("could not find track analyisis for" + rec.ID)
			log.Print(err)
		}
		toptracks = append(toptracks, structs.TopTracksNode{
			ID:           rec.ID,
			Popularity:   strconv.Itoa(rec.Popularity),
			Artist:       rec.Artists[0].Name,
			Title:        rec.Name,
			ImageURL:     rec.Album.Images[0].URL,
			Key:          helpers.ParseKey(analysis.Key),
			Tempo:        helpers.ParseTempo(analysis.Tempo),
			Loudness:     helpers.ParseLoudness(analysis.Loudness),
			Happiness:    helpers.ParsePercent(analysis.Valence),
			Energy:       helpers.ParsePercent(analysis.Energy),
			Danceability: helpers.ParsePercent(analysis.Danceability),
			Acousticness: helpers.ParsePercent(analysis.Acousticness),
		})
	}
	data := &structs.TopTracksPage{
		Name:  "Trevor",
		Items: toptracks,
	}

	renderTopTracksTemplate(w, data)
}

func renderTopTracksTemplate(w http.ResponseWriter, p *structs.TopTracksPage) {
	e := templates.ExecuteTemplate(w, "top-tracks.html", p)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}

func getTrackAnalysis(id string) (*structs.TrackAnalysisResponse, error) {
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
	var result structs.TrackAnalysisResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Print("Unable to unmarshal JSON response.")
		log.Print(err)
	}
	return &result, nil
}
