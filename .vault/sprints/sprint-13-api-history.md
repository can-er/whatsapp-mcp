---
type: sprint
id: sprint-13
title: "GET /api/history paginé au format des événements v1 (fork)"
status: TODO          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed:
tags: [sprint]
---

# Sprint 13 — GET /api/history paginé au format des événements v1 (fork)

> **status:: TODO**  ·  **test:: [[TC-13-api-history]]**

## 🎯 Objectif

`GET /api/history?chat_jid=…&before=…&limit=…` (jeton du pont, liste blanche appliquée) renvoie l'historique stocké au format des événements v1 (média `not_downloaded` tant qu'il n'est pas téléchargé), paginé par curseur.

## 📋 Tâches

- [ ] Requête SQL : chat (numéro + LID résolu), tri décroissant timestamp + id, curseur `next_before`
- [ ] Conversion ligne → événement v1 (expéditeur reconstruit via LID store, média téléchargé ou `not_downloaded`)
- [ ] 403 hors liste blanche, 400 paramètres invalides, `limit` 50 / max 500
- [ ] README

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-13-api-history]]

## ✅ Definition of Done

- [ ] Tâches terminées
- [ ] Test [[TC-13-api-history]] **PASS**
- [ ] Commit fait, message préfixé `[Sprint 13]`
- [ ] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [ ] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé à la demande de ks-support-agent-6b (spécification reçue, cf. ADR-0006).

## 🔗 Liens

- ADR : [[ADR-0006-interface-pont-pipeline]]
