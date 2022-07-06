package handlers

import (
	"html/template"
	"os"
)

const clientid = "58f7888547c24f058ece41fed973bf37" // Your client id
const statekey = "spotify_auth_state"

var clientsecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
var redirecturi = "http://localhost:8888/callback" // Your redirect uri
var templates = template.Must(template.ParseFiles("tmpl/top-tracks.html"))
