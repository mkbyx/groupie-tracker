package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type List struct {
	Lists []Artiste
}

type Locations struct {
	LocationsStruct []string `json:"locations"`
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


func fetchLocations(artists []Artiste) map[int][]string {
	locationsMap := make(map[int][]string)

	for _, artist := range artists {
		res, err := http.Get(artist.Locations)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		var locations Locations
		json.Unmarshal(body, &locations)
		for i, loca := range locations.LocationsStruct {
			locations.LocationsStruct[i] = strings.ToLower(loca)
		}
		locationsMap[artist.Id] = locations.LocationsStruct
	}
	return locationsMap
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

	var groupList List
	json.Unmarshal(body, &groupList.Lists)
	locationMap := fetchLocations(groupList.Lists)
	
	queryYear := r.URL.Query().Get("year")
	queryYearAlbum := r.URL.Query().Get("yearAlbum")
	queryMembers := r.URL.Query().Get("members")
	searchQuery := r.URL.Query().Get("query")

	var filteredArtists []Artiste
	filteredArtists = groupList.Lists
	if searchQuery != "" {
		searchQuery = strings.ToLower(searchQuery)
		tempFilteredArtists := []Artiste{}
		
		for _, artist := range filteredArtists {
			if strings.Contains(strings.ToLower(artist.Name), searchQuery) {
				tempFilteredArtists = append(tempFilteredArtists, artist)
			}else if containsInList(artist.Members,searchQuery){
				tempFilteredArtists = append(tempFilteredArtists, artist)
			}else if strings.Contains(strconv.Itoa(artist.CreationDate), searchQuery){
				tempFilteredArtists = append(tempFilteredArtists, artist)
			}else if strings.Contains(artist.FirstAlbum, searchQuery){
				tempFilteredArtists = append(tempFilteredArtists, artist)
			}else if containsInList(locationMap[artist.Id],searchQuery){
				tempFilteredArtists = append(tempFilteredArtists, artist)
			}
		}
		filteredArtists = tempFilteredArtists 
	}

	if queryYear != "" {
		year, err := strconv.Atoi(queryYear)
		if err == nil {
			tempFilteredArtists := []Artiste{}
			for _, artist := range filteredArtists {
				if artist.CreationDate >= year && artist.CreationDate <= year+10 {
					tempFilteredArtists = append(tempFilteredArtists, artist)
				}
			}
			filteredArtists = tempFilteredArtists 
		}
	}

	if queryYearAlbum != "" {
		yearAlbum, err := strconv.Atoi(queryYearAlbum)
		if err == nil {
			tempFilteredArtists := []Artiste{}
			for _, artist := range filteredArtists {
				dateParts := strings.Split(artist.FirstAlbum, "-")
				if len(dateParts) == 3 {
					artistFirstAlbumYear, convErr := strconv.Atoi(dateParts[2])
					if convErr == nil && artistFirstAlbumYear >= yearAlbum && artistFirstAlbumYear <= yearAlbum+10 {
						tempFilteredArtists = append(tempFilteredArtists, artist)
					}
				}
			}
			filteredArtists = tempFilteredArtists
		}
	}

	if queryMembers != "" {
		memberCount, err := strconv.Atoi(queryMembers)
		if err == nil {
			tempFilteredArtists := []Artiste{}
			for _, artist := range filteredArtists {
				if len(artist.Members) == memberCount {
					tempFilteredArtists = append(tempFilteredArtists, artist)
				}
			}
			filteredArtists = tempFilteredArtists
		}
	}

	groupList.Lists = filteredArtists
	t.Execute(w, groupList)
}

func containsInList(list []string, query string) bool {
	for _, item := range list {
		if strings.Contains(strings.ToLower(item), query) {
			return true
		}
	}
	return false
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
