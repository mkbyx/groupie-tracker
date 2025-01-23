package groupietracker

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
