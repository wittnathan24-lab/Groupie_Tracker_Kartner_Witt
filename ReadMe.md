# Groupie Tracker – Projet Web 🎸🎤

## Objectif du projet

L'objectif de ce projet est de créer une application web qui consomme une API publique fournie (par exemple l’API *Groupie Trackers*) pour afficher des informations sur des artistes et groupes de musique. L’application doit permettre à l’utilisateur de parcourir les artistes, d’afficher des détails les concernant, et de naviguer facilement entre les différentes vues.

Ce projet met en pratique :
- la consommation d’une API REST,
- la manipulation de données JSON,
- la création d’un serveur web en Go,
- l’affichage dynamique de données via HTML/CSS/JavaScript.

## Comment lancer le serveur

1. Clonez ce dépôt sur votre machine :
   ```bash
   git clone https://github.com/wittnathan24-lab/Groupie_Tracker_Kartner_Witt.git
   cd Groupie_Tracker_Kartner_Witt
```

2. Assurez-vous d’avoir installé Go (version minimale recommandée : **Go 1.18+**).
3. Dans le dossier racine du projet, lancez :

   ```bash
   go run .
   ```

4. Ouvrez votre navigateur et allez sur :

   ```
   http://localhost:8080
   ```

   (ou le port configuré dans votre projet)

## Routes principales

| Route           | Méthode | Description                                     |
|-----------------| ------- |-------------------------------------------------|
| `/index`        | GET     | Page d’accueil                                  |
| `/Artiste`      | GET     | Récupère tous les artistes                      |
| `/Artiste/{id}` | GET     | Page de détails d’un artiste                    |
| `/Liste`        | GET     | Retourne les données JSON de tous les artistes  |


> Les routes ci-dessus sont des exemples standards — adapte-les à celles que tu as réellement implémentées.

## Fonctionnalités implémentées

### Fonctionnalités obligatoires

* Consommation de l’API externe pour récupérer les données d’artistes.
* Serveur web en Go répondant aux requêtes HTTP.
* Affichage dynamique des artistes via une interface web.
* Page de détails pour chaque artiste (nom, date de début, membres, etc.).
* Navigation entre les pages artistes.

### Bonus (si implémentés)

* Barre de recherche dynamique pour filtrer les artistes par nom.
* Carte interactive des lieux de concerts (géolocalisation).
* Filtrage par date de concert / années d’activité.
* Visualisation des données (graphiques des années, des pays, etc.).
* Tests unitaires pour les handlers ou fonctions Go.

## Technologies utilisées

* **Go** – Backend / serveur HTTP
* **HTML5 / CSS3** – Interface utilisateur
* **JavaScript** – Dynamisation des pages
* (Éventuel **Leaflet / Mapbox**) – Cartographie des lieux de concerts

## Structure du projet

```
├─ cmd/
├─ handlers/
├─ static/
│   ├─ css/
│   └─ js/
├─ templates/
├─ main.go
├─ go.mod
└─ README.md
```

## Remarques & bonnes pratiques

* Le serveur ne doit **jamais** planter, même en cas d’erreur de requête API.
* Gestion d’erreurs claire et renvoi de statuts HTTP appropriés (ex : 404, 500).
* Code structuré avec des responsables clairs pour chaque fonctionnalité.

---