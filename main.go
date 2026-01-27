package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

// --- VARIABLES GLOBALES ---
var ListOfArtists []Artist
var httpClient = &http.Client{
	Timeout: 15 * time.Second, // Timeout augmenté pour éviter les erreurs sur connexions lentes
}

// --- STRUCTURES ---
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

type FilterData struct {
	MinCreation int
	MaxCreation int
	Members     map[int]bool
}

type ErrorPageData struct {
	Code    int
	Message string
	Details string
}

// --- MAIN ---
func main() {
	// Servir les fichiers statiques (CSS, Images)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes avec protection contre les crashs
	http.HandleFunc("/Index", safeHandler(IndexPage))
	http.HandleFunc("/Artiste/", safeHandler(ArtistePage))
	http.HandleFunc("/Liste", safeHandler(ListePage))
	http.HandleFunc("/api/search", SearchAPI) // API gère ses erreurs différemment (JSON)
	http.HandleFunc("/toggle-theme", ToggleThemeHandler)

	// Gestion de la racine pour rediriger ou 404
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/Index", http.StatusSeeOther)
		} else {
			renderError(w, http.StatusNotFound, "Page introuvable", "La route demandée n'existe pas : "+r.URL.Path)
		}
	})

	fmt.Println("Serveur démarré : http://localhost:8080/Index")
	// Le log.Fatal ici est le SEUL point d'arrêt acceptable (si le port est occupé par exemple)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// --- GESTION DES ERREURS & SECURITE ---
// safeHandler est un wrapper pour les handlers HTTP.
// Si la fonction panic (crash), il récupère la main et affiche une erreur 500 au lieu de tuer le serveur.
// - "Aucun crash serveur n'est accepté"
func safeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("🔥 PANIC RECOVERED: %v\nStack: %s", err, debug.Stack())
				renderError(w, http.StatusInternalServerError, "Erreur Critique du Serveur", "Une erreur inattendue s'est produite. Nos équipes ont été notifiées.")
			}
		}()
		fn(w, r)
	}
}

// renderError affiche le template Error.html avec les détails précis
// - "Pages d'erreur personnalisées"
func renderError(w http.ResponseWriter, code int, message, details string) {
	w.WriteHeader(code) // Définit le code HTTP (ex: 404, 500)
	tmpl, err := template.ParseFiles("template/Error.html")
	if err != nil {
		// Fallback ultime si le template d'erreur est cassé
		http.Error(w, fmt.Sprintf("Erreur critique (%d): %s", code, message), code)
		return
	}
	tmpl.Execute(w, ErrorPageData{
		Code:    code,
		Message: message,
		Details: details,
	})
}

func FetchArtists() ([]Artist, error) {
	resp, err := httpClient.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, fmt.Errorf("impossible de contacter l'API distante")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("l'API a répondu avec le code : %d", resp.StatusCode)
	}

	var artists []Artist
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, fmt.Errorf("les données reçues sont corrompues")
	}
	return artists, nil
}

func getDarkMode(r *http.Request) bool {
	cookie, err := r.Cookie("theme")
	return err == nil && cookie.Value == "dark"
}

// --- HANDLERS ---

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
		MaxAge: 60 * 60 * 24 * 30, // 30 jours
	})
	// Redirection vers la page précédente
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	// Validation stricte de la méthode
	if r.Method != http.MethodGet {
		renderError(w, http.StatusMethodNotAllowed, "Méthode non autorisée", "Seul GET est autorisé ici.")
		return
	}

	tmpl, err := template.ParseFiles("template/Index.html")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Erreur de chargement", "Impossible de charger la page d'accueil.")
		return
	}
	tmpl.Execute(w, struct{ DarkMode bool }{getDarkMode(r)})
}

func ArtistePage(w http.ResponseWriter, r *http.Request) {
	// 1. Chargement des données
	if len(ListOfArtists) == 0 {
		var err error
		ListOfArtists, err = FetchArtists()
		if err != nil {
			renderError(w, http.StatusServiceUnavailable, "Service indisponible", "L'API GroupieTracker ne répond pas : "+err.Error())
			return
		}
	}

	// 2. Extraction et Validation de l'ID
	idStr := strings.TrimPrefix(r.URL.Path, "/Artiste/")
	id, err := strconv.Atoi(idStr)
	// - "erreurs de paramètres"
	if err != nil || id <= 0 {
		renderError(w, http.StatusBadRequest, "Requête invalide", "L'identifiant de l'artiste doit être un nombre positif.")
		return
	}

	// 3. Recherche de l'artiste
	var selected *Artist
	for i := range ListOfArtists {
		if ListOfArtists[i].ID == id {
			selected = &ListOfArtists[i]
			break
		}
	}

	// - "404"
	if selected == nil {
		renderError(w, http.StatusNotFound, "Artiste introuvable", fmt.Sprintf("Aucun artiste trouvé avec l'ID %d.", id))
		return
	}

	// 4. Récupération des relations (Concerts)
	// On ne bloque pas la page si les relations échouent, on log juste l'erreur
	var rel Relation
	resp, err := httpClient.Get(selected.Relations)
	if err == nil {
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&rel)
	} else {
		log.Printf("Warning: Impossible de charger les relations pour l'artiste %d: %v", id, err)
	}

	// 5. Rendu
	tmpl, err := template.ParseFiles("template/Artiste.html")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Erreur d'affichage", "Le template de la page artiste est manquant.")
		return
	}

	data := struct {
		Artist
		RelationsMap map[string][]string
		DarkMode     bool
	}{Artist: *selected, RelationsMap: rel.DatesLocations, DarkMode: getDarkMode(r)}

	tmpl.Execute(w, data)
}

func ListePage(w http.ResponseWriter, r *http.Request) {
	// Chargement sécurisé
	if len(ListOfArtists) == 0 {
		var err error
		ListOfArtists, err = FetchArtists()
		if err != nil {
			renderError(w, http.StatusServiceUnavailable, "Service indisponible", "Impossible de récupérer la liste des artistes.")
			return
		}
	}

	// Gestion des filtres (paramètres URL)
	minDateStr := r.URL.Query().Get("min_creation")
	maxDateStr := r.URL.Query().Get("max_creation")
	membersParams := r.URL.Query()["members"]

	// Conversion sécurisée des paramètres (on ignore les erreurs pour utiliser les défauts)
	minDate, _ := strconv.Atoi(minDateStr)
	maxDate, err := strconv.Atoi(maxDateStr)
	if err != nil || maxDate == 0 {
		maxDate = 2030
	}

	selectedMembers := make(map[int]bool)
	for _, m := range membersParams {
		val, err := strconv.Atoi(m)
		if err == nil && val > 0 {
			selectedMembers[val] = true
		}
	}

	// Filtrage
	var filteredArtists []Artist
	for _, artist := range ListOfArtists {
		if artist.Created < minDate || artist.Created > maxDate {
			continue
		}
		if len(selectedMembers) > 0 && !selectedMembers[len(artist.Members)] {
			continue
		}
		filteredArtists = append(filteredArtists, artist)
	}

	tmpl, err := template.ParseFiles("template/Liste.html")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Erreur interne", "Impossible d'afficher la liste.")
		return
	}

	data := struct {
		Artists  []Artist
		DarkMode bool
		Filters  FilterData
	}{
		Artists:  filteredArtists,
		DarkMode: getDarkMode(r),
		Filters: FilterData{
			MinCreation: minDate,
			MaxCreation: maxDate,
			Members:     selectedMembers,
		},
	}
	tmpl.Execute(w, data)
}

func SearchAPI(w http.ResponseWriter, r *http.Request) {
	// API JSON : On ne renvoie pas de HTML en cas d'erreur
	w.Header().Set("Content-Type", "application/json")

	q := strings.ToLower(r.URL.Query().Get("q"))

	if len(ListOfArtists) == 0 {
		var err error
		ListOfArtists, err = FetchArtists()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{"error": "API GroupieTracker indisponible"})
			return
		}
	}

	type item struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Image string `json:"image"`
	}
	var res []item
	// Limite à 8 résultats pour la performance
	for _, a := range ListOfArtists {
		if strings.Contains(strings.ToLower(a.Name), q) && len(res) < 8 {
			res = append(res, item{a.ID, a.Name, a.Image})
		}
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Erreur encodage JSON: %v", err)
	}
}
