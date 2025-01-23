package main

import (
	"fmt"
	"net/http"
	gt "groupietracker/back"
)

func main() {
	staticFiles := http.FileServer(http.Dir("./static"))
	http.Handle("/styles/", http.StripPrefix("/styles/", staticFiles))
	http.HandleFunc("/", gt.MainPage)
	http.HandleFunc("/artiste", gt.ArtistPage)
	fmt.Println("Serveur démarré sur : http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
