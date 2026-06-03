# Mediadex-API

**Équipe :** [Jérémy Gaudin](https://github.com/jgaudin826), [Aurélien Dugast](https://github.com/LaNuggets), [Cassian Joly](https://github.com/Cassian-J)
**Groupe / promo :** JAC B3 Dev Ynov
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
| Orchestration cible | Compose puis K8s |  |

---

## Rôles dans l'équipe

| Membre | Rôle | Responsabilité principale |
| ------ | ---- | ------------------------- |
| Cassian Joly | Lead Dev | Création du Dockerfile service |
| Jérémy Gaudin | Lead Ops | Docker compose, K8s, doc de déploiement |
|  | Lead Qualité / CI | Pipelines, test, revue de sécurité |
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
| S3 | docker-compose + CI vert | ☐ |
| S4 | Manifests K8s appliqués | ☐ |
| S5 | Monitoring + post-mortem | ☐ |
| S6 | Soutenance prête | ☐ |

---

## Démarrage local (à compléter au fil des séances)

```bash
git clone https://github.com/jgaudin826/devops-fil-rouge-JAC.git
```

---

## Communication d'équipe

Canal utilisé : Discord
