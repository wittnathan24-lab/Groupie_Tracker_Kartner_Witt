# Groupie Tracker – Projet Web 🎸🎤

## Objectif du projet

L'objectif de ce projet est de créer une application web qui consomme une API publique fournie pour afficher des informations sur des artistes et groupes de musique. L’application doit permettre à l’utilisateur de parcourir les artistes, d’afficher des détails les concernant, et de naviguer facilement entre les différentes vues.

Ce projet met en pratique :
- la consommation d’une API REST,
- la manipulation de données JSON,
- la création d’un serveur web en Go,
- l’affichage dynamique de données via HTML/CSS.

## Comment lancer le serveur

1. Clonez ce dépôt sur votre machine :
   ```bash
   git clone https://github.com/wittnathan24-lab/Groupie_Tracker_Kartner_Witt.git
   cd Groupie_Tracker_Kartner_Witt

2. Assurez-vous d’avoir installé Go (version minimale recommandée : **Go 1.18+**).
3. Dans le dossier racine du projet, lancez :

   ```bash
   go run .


4. Ouvrez votre navigateur et allez sur :

   ```
   http://localhost:8080/
   ```


## Routes principales

| Route           | Méthode | Description                                     |
|-----------------| ------- |-------------------------------------------------|
| `/index`        | GET     | Page d’accueil                                  |
| `/Artiste`      | GET     | Récupère tous les artistes                      |
| `/Artiste/{id}` | GET     | Page de détails d’un artiste                    |
| `/Liste`        | GET     | Retourne les données JSON de tous les artistes  |


## Fonctionnalités implémentées

### Fonctionnalités obligatoires

* Consommation de l’API externe pour récupérer les données d’artistes.
* Serveur web en Go répondant aux requêtes HTTP.
* Affichage dynamique des artistes via une interface web.
* Page de détails pour chaque artiste (nom, date de début, membres, etc.).

### Bonus

* Barre de recherche dynamique pour filtrer les artistes par nom.
* Filtrage par date de concert / années d’activité.

## Technologies utilisées

* **Go** – Backend / serveur HTTP
* **HTML5 / CSS3** – Interface utilisateur

## Structure du projet

```

├─ static/
│   ├─ css/
│         ├─global.css
│         ├─Liste.css
│         └─Artiste.css
├─ templates/
│   ├─Index.html
│   ├─Artiste.html
│   ├─Liste.html
│   └─Error.html
├─ main.go
├─ go.mod
└─ README.md
```

## Remarques & bonnes pratiques

* Gestion d’erreurs claire et renvoi de statuts HTTP appropriés (ex : 404, 500).
* Code structuré avec des responsables clairs pour chaque fonctionnalité.

## Contributions
Réalisé par WITT Nathan et KARTNER Allan
---
