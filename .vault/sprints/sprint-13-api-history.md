---
type: sprint
id: sprint-13
title: "GET /api/history paginé au format des événements v1 (fork)"
status: DONE          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed: 2026-10-29
tags: [sprint]
---

# Sprint 13 — GET /api/history paginé au format des événements v1 (fork)

> **status:: DONE**  ·  **test:: [[TC-13-api-history]]**

## 🎯 Objectif

`GET /api/history?chat_jid=…&before=…&limit=…` (jeton du pont, liste blanche appliquée) renvoie l'historique stocké au format des événements v1 (média `not_downloaded` tant qu'il n'est pas téléchargé), paginé par curseur.

## 📋 Tâches

- [x] Requête SQL : chat (numéro + LID résolu), tri décroissant timestamp + id, curseur `next_before`
- [x] Conversion ligne → événement v1 (expéditeur reconstruit via LID store, média téléchargé ou `not_downloaded`)
- [x] 403 hors liste blanche, 400 paramètres invalides, `limit` 50 / max 500
- [x] README

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-13-api-history]]

## ✅ Definition of Done

- [x] Tâches terminées
- [x] Test [[TC-13-api-history]] **PASS**
- [x] Commit fait, message préfixé `[Sprint 13]`
- [x] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [x] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé à la demande de ks-support-agent-6b (spécification reçue, cf. ADR-0006).
- `2026-10-29` — GET /api/history paginé au format v1 ; TC-13 PASS ; mergé dans main.

## 🔗 Liens

- ADR : [[ADR-0006-interface-pont-pipeline]]
