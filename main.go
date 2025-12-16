package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/Index", IndexPage)
	http.HandleFunc("/Artiste", ArtistePage)
	http.HandleFunc("/Liste", ListePage)
	fmt.Println("Serveur démarré sur http://localhost:8080")
	fmt.Println("Accédez à http://localhost:8080/login pour commencer.")
	http.ListenAndServe(":8080", nil)
}

var ListOfArtists []Artist

type Artist struct {
	ID         int
	Name       string
	Image      string
	Members    []string
	Created    int
	FirstAlbum string
	Locations  string
	Concerts   string
	Relations  string
}

func FetchArtists() ([]Artist, error) {
	const url = "https://groupietrackers.herokuapp.com/api/artists"
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	var artists []Artist
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, err
	}
	return artists, nil
}

func IndexPage(w http.ResponseWriter, r *http.Request) {

}

func ArtistePage(w http.ResponseWriter, r *http.Request) {

}

func ListePage(w http.ResponseWriter, r *http.Request) {

}
