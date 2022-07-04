package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"spotify/pkgs/handlers"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8888, "Listener port")
	flag.Parse()

	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/callback", handlers.Callback)
	http.HandleFunc("/refresh_token", handlers.RefreshToken)
	http.HandleFunc("/top-tracks", handlers.TopTracks)
	//http.HandleFunc("/track/", trackHandler)

	var frontend fs.FS = os.DirFS("public")
	httpFS := http.FS(frontend)
	fileServer := http.FileServer(httpFS)
	http.Handle("/", fileServer)

	addr := fmt.Sprintf("localhost:%d", port)
	log.Printf("Serving app at http://%s", addr)
	log.Fatalln(http.ListenAndServe(addr, nil))
}
