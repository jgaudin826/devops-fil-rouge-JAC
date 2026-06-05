# Mediadex-API

**Équipe :**
bonjour

[Jérémy Gaudin](https://github.com/jgaudin826), 

[Aurélien Dugast](https://github.com/LaNuggets), 

[Cassian Joly](https://github.com/Cassian-J)

**Groupe / promo :** JAC / B3 Dev Ynov


**Dépôt :** https://github.com/jgaudin826/devops-fil-rouge-JAC.git

---

## Description du sujet


Mediadex-API est une API de gestion de médias permettant le stockage et l'organisation de médias, de collections, de tags personnalisés et d'utilisateurs.

Elle est destinée aux personnes qui consomment différents types de médias (séries, livres, films, jeux vidéo, etc.) et qui souhaitent suivre leur progression, conserver un historique de leurs contenus terminés, leur attribuer des notes ou des commentaires personnalisés, ainsi que gérer une liste de contenus à découvrir plus tard (livres à lire, films à regarder, jeux à essayer, etc.).

---

## Stack technique prévue

| Composant | Choix | Justification |
| --------- | ----- | -------------------------- |
| Backend / API | Go, gorm, chi | Bonne maîtrise de ces technologies |
| Base de données | PostgreSQL | Pour le challenge, découvrir de nouvelles choses. C'est une technologie ancienne, donc bien documentée et robuste. |
| Orchestration cible | Compose puis K8s | Simplicité d'éxecution, versioning. |
| Image Dockerfile | golang:1.26.4-alpine | Stable, version présice, légère |

---

## Rôles dans l'équipe

| Membre | Rôle | Responsabilité principale |
| ------ | ---- | ------------------------- |
| Cassian Joly | Lead Dev | Création du Dockerfile service |
| Jérémy Gaudin | Lead Ops | Docker compose, K8s, doc de déploiement |
| Aurélien Dugast | Lead Qualité / CI | Pipelines, test, revue de sécurité |
| Aurélien Dugast | Lead Doc / Produit | Documentation, note-archi et analyse post-mortem |

---

## Objectifs du fil rouge

1. Avoir une API conteneurisée avec healthcheck d'ici S3.
2. Pipeline CI qui build et push l'image sur chaque merge `main`.
3. Déployer sur cluster kind avec 2 replicas d'ici S4.

---

## Jalons — état d'avancement

| Séance | Livrable | Statut (à cocher) |
| ------ | -------- | ----------------- |
| S1 | README cadrage | X |
| S2 | Dockerfile(s) + DB en container | X |
| S3 | docker-compose | X |
| S5 | Pipelines CI | X |
| S6 | Soutenance prête | X |

---

## Démarrage local

```bash
# Cloner le dépôt
git clone https://github.com/jgaudin826/devops-fil-rouge-JAC.git
cd devops-fil-rouge-JAC

cp .env.example .env
# Renseigner PORT, JWT_SECRET_KEY, POSTGRES_USER, etc.
go mod tidy
go get
go run main.go
```

## Démarrage avec Docker compose
```
docker compose up -d
```

---

## Communication d'équipe

Canal utilisé : Discord
