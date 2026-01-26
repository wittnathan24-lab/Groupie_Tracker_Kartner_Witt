package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

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

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/Index", IndexPage)
	http.HandleFunc("/Artiste/", ArtistePage)
	http.HandleFunc("/Liste", ListePage)
	http.HandleFunc("/api/search", SearchAPI)
	http.HandleFunc("/toggle-theme", ToggleThemeHandler)

	fmt.Println("🚀 Serveur : http://localhost:8080/Index")
	http.ListenAndServe(":8080", nil)
}

func getDarkMode(r *http.Request) bool {
	cookie, err := r.Cookie("theme")
	return err == nil && cookie.Value == "dark"
}

func FetchArtists() []Artist {
	resp, _ := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	defer resp.Body.Close()
	var artists []Artist
	json.NewDecoder(resp.Body).Decode(&artists)
	return artists
}

func ToggleThemeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("theme")
	newTheme := "dark"
	if err == nil && cookie.Value == "dark" {
		newTheme = "light"
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "theme",
		Value:  newTheme,
		Path:   "/",
		MaxAge: 60 * 60 * 24 * 30,
	})
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("template/Index.html")
	tmpl.Execute(w, struct{ DarkMode bool }{getDarkMode(r)})
}

func ArtistePage(w http.ResponseWriter, r *http.Request) {
	if len(ListOfArtists) == 0 {
		ListOfArtists = FetchArtists()
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/Artiste/")
	id, _ := strconv.Atoi(idStr)
	var selected Artist
	for _, a := range ListOfArtists {
		if a.ID == id {
			selected = a
			break
		}
	}
	resp, _ := http.Get(selected.Relations)
	defer resp.Body.Close()
	var rel Relation
	json.NewDecoder(resp.Body).Decode(&rel)

	data := struct {
		Artist
		RelationsMap map[string][]string
		DarkMode     bool
	}{Artist: selected, RelationsMap: rel.DatesLocations, DarkMode: getDarkMode(r)}
	tmpl, _ := template.ParseFiles("template/Artiste.html")
	tmpl.Execute(w, data)
}

func ListePage(w http.ResponseWriter, r *http.Request) {
	if len(ListOfArtists) == 0 {
		ListOfArtists = FetchArtists()
	}
	tmpl, _ := template.ParseFiles("template/Liste.html")
	data := struct {
		Artists  []Artist
		DarkMode bool
	}{Artists: ListOfArtists, DarkMode: getDarkMode(r)}
	tmpl.Execute(w, data)
}

func SearchAPI(w http.ResponseWriter, r *http.Request) {
	q := strings.ToLower(r.URL.Query().Get("q"))
	if len(ListOfArtists) == 0 {
		ListOfArtists = FetchArtists()
	}
	type item struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Image string `json:"image"`
	}
	var res []item
	for _, a := range ListOfArtists {
		if strings.Contains(strings.ToLower(a.Name), q) && len(res) < 8 {
			res = append(res, item{a.ID, a.Name, a.Image})
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
