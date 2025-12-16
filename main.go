package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/Index", IndexPage)
	// Search endpoint (expects ?q=...)
	http.HandleFunc("/Artiste", ArtistePage)
	// Detail endpoint: /Artiste/{id}
	http.HandleFunc("/Artiste/", ArtistePage)
	http.HandleFunc("/Liste", ListePage)
	fmt.Println("Serveur démarré sur http://localhost:8080")
	fmt.Println("Accédez à http://localhost:8080/Index pour commencer.")
	http.ListenAndServe(":8080", nil)
}

var ListOfArtists []Artist

type Artist struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Image      string   `json:"image"`
	Members    []string `json:"members"`
	Created    int      `json:"creationDate"`
	FirstAlbum string   `json:"firstAlbum"`
	Locations  string   `json:"locations"`
	Concerts   string   `json:"concertDates"`
	Relations  string   `json:"relations"`
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
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("template/Index.html")
	if err != nil {
		http.Error(w, "Erreur template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Erreur rendu: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func containsMember(members []string, qLower string) bool {
	for _, m := range members {
		if strings.Contains(strings.ToLower(m), qLower) {
			return true
		}
	}
	return false
}

func ArtistePage(w http.ResponseWriter, r *http.Request) {
}

func ListePage(w http.ResponseWriter, r *http.Request) {

}
