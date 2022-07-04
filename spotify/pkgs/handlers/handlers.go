package handlers

import (
	"html/template"
)

const clientid = "58f7888547c24f058ece41fed973bf37"     // Your client id
const clientsecret = "8405db97ae294a9885f064855b5fae79" // Your secret
const statekey = "spotify_auth_state"

var redirecturi = "http://localhost:8888/callback" // Your redirect uri
var templates = template.Must(template.ParseFiles("tmpl/top-tracks.html"))
