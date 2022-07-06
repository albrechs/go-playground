package structs

type TopTracksPage struct {
	Name  string
	Items []TopTracksNode
}

type TopTracksNode struct {
	ID           string
	Popularity   string
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

type LoginOpts struct {
	ResponseType string `url:"response_type"`
	ClientID     string `url:"client_id"`
	Scope        string `url:"scope"`
	RedirectURI  string `url:"redirect_uri"`
	State        string `url:"state"`
}

type CallbackRedirectQuery struct {
	AccessToken  string `url:"access_token"`
	RefreshToken string `url:"refresh_token"`
}
