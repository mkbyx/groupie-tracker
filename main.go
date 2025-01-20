package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

type List struct {
	Lists []Artiste
}

type Artiste struct {
	Id           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Locations    string   `json:"locations"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationdate"`
	FirstAlbum   string   `json:"firstalbum"`
	ConcertDates string   `json:"concertdates"`
	Relations    string   `json:"relations"`
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./static/home.html")
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var GroupList List
	json.Unmarshal(body, &GroupList.Lists)
	t.Execute(w, GroupList)
}

func artistPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./static/artiste.html")
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var GroupList List
	json.Unmarshal(body, &GroupList.Lists)
	artistID := r.URL.Query().Get("id")
	if artistID == "" {
		http.Error(w, "ID de l'artiste non spécifié", http.StatusBadRequest)
		return
	}
	var selectedArtist *Artiste
	for _, artist := range GroupList.Lists {
		if fmt.Sprintf("%d", artist.Id) == artistID {
			selectedArtist = &artist
			break
		}
	}
	if selectedArtist == nil {
		http.Error(w, "Artiste non trouvé", http.StatusNotFound)
		return
	}
	t.Execute(w, selectedArtist)
}

func main() {
	staticFiles := http.FileServer(http.Dir("./static"))
	http.Handle("/styles/", http.StripPrefix("/styles/", staticFiles))
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/artiste", artistPage)
	fmt.Println("Serveur démarré sur : http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

