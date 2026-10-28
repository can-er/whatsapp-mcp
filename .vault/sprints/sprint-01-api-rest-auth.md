---
type: sprint
id: sprint-01
title: "API REST du pont : loopback + jeton + JSON only"
status: DONE          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed: 2026-10-29
tags: [sprint]
---

# Sprint 01 — API REST du pont : loopback + jeton + JSON only

> **status:: DONE**  ·  **test:: [[TC-01-api-rest-auth]]**

## 🎯 Objectif

Le pont n'écoute plus que sur 127.0.0.1 par défaut et refuse toute requête sans jeton valide, sans `Content-Type: application/json`, avec un en-tête `Origin`, ou dont le corps dépasse 1 MiB.

## 📋 Tâches

- [x] `WHATSAPP_BRIDGE_ADDR` (défaut `127.0.0.1:8080`)
- [x] Jeton : `WHATSAPP_BRIDGE_TOKEN` ou `store/bridge_token` généré au 1er lancement (0600)
- [x] Middleware : Origin → 403, Content-Type → 415, jeton (comparaison à temps constant) → 401, `MaxBytesReader` 1 MiB
- [x] `http.Server` dédié (mux propre, timeouts)
- [x] Client Python : en-tête `Authorization: Bearer` lu depuis l'env ou le fichier
- [x] README : section sécurité

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-01-api-rest-auth]]

## ✅ Definition of Done

- [x] Tâches terminées
- [x] Test [[TC-01-api-rest-auth]] **PASS**
- [x] Commit fait, message préfixé `[Sprint 01]`
- [x] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [x] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé.
- `2026-10-29` — Auth API (loopback, Bearer, JSON only, Origin, 1 MiB) ; TC-01 PASS ; mergé dans main.

## 🔗 Liens

- ADR : [[ADR-0004-auth-pont-jeton]]
