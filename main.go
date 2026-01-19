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
	// Dynamic search suggestions
	http.HandleFunc("/api/search", SearchAPI)
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
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Load artists if not loaded
	if len(ListOfArtists) == 0 {
		artists, err := FetchArtists()
		if err != nil {
			http.Error(w, "Erreur chargement: "+err.Error(), http.StatusInternalServerError)
			return
		}
		ListOfArtists = artists
	}

	// 1. Check for specific ID in URL path: /Artiste/1
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) > 1 && pathParts[1] != "" {
		// Handle ID logic (template rendering for single artist)
		// For now, we assume you have logic to show Artiste.html
		tmpl, _ := template.ParseFiles("template/Artiste.html")
		// Find artist by ID and execute tmpl...
		_ = tmpl
		return
	}

	// 2. Check for search query: /Artiste?q=queen
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q != "" {
		qLower := strings.ToLower(q)
		for _, a := range ListOfArtists {
			if strings.ToLower(a.Name) == qLower {
				// Exact match found: redirect to the ID-based URL
				http.Redirect(w, r, fmt.Sprintf("/Artiste/%d", a.ID), http.StatusFound)
				return
			}
		}
	}

	// If no exact match or no query, you might want to show a list or 404
	http.Error(w, "Artiste non trouvé", http.StatusNotFound)
}

func ListePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	if len(ListOfArtists) == 0 {
		artists, err := FetchArtists()
		if err != nil {
			http.Error(w, "Erreur chargement API: "+err.Error(), http.StatusInternalServerError)
			return
		}
		ListOfArtists = artists
	}

	tmpl, err := template.ParseFiles("template/Liste.html")
	if err != nil {
		http.Error(w, "Erreur template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, ListOfArtists); err != nil {
		http.Error(w, "Erreur rendu: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
func SearchAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		// empty query -> empty list (avoid noise)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("[]"))
		return
	}
	if len(ListOfArtists) == 0 {
		artists, err := FetchArtists()
		if err != nil {
			http.Error(w, "Erreur lors du chargement des artistes: "+err.Error(), http.StatusInternalServerError)
			return
		}
		ListOfArtists = artists
	}
	qLower := strings.ToLower(q)
	type item struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Image string `json:"image"`
	}
	seen := make(map[int]bool)
	results := make([]item, 0, 8)

	add := func(a Artist) {
		if seen[a.ID] {
			return
		}
		seen[a.ID] = true
		results = append(results, item{ID: a.ID, Name: a.Name, Image: a.Image})
	}
	for _, a := range ListOfArtists {
		if len(results) >= 8 {
			break
		}
		if strings.HasPrefix(strings.ToLower(a.Name), qLower) {
			add(a)
		}
	}
	for _, a := range ListOfArtists {
		if len(results) >= 8 {
			break
		}
		if strings.Contains(strings.ToLower(a.Name), qLower) {
			add(a)
		}
	}
	for _, a := range ListOfArtists {
		if len(results) >= 8 {
			break
		}
		if containsMember(a.Members, qLower) {
			add(a)
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	if err := enc.Encode(results); err != nil {
		http.Error(w, "Erreur d'encodage JSON", http.StatusInternalServerError)
		return
	}
}
