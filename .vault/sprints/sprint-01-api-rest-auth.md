---
type: sprint
id: sprint-01
title: "API REST du pont : loopback + jeton + JSON only"
status: TODO          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed:
tags: [sprint]
---

# Sprint 01 — API REST du pont : loopback + jeton + JSON only

> **status:: TODO**  ·  **test:: [[TC-01-api-rest-auth]]**

## 🎯 Objectif

Le pont n'écoute plus que sur 127.0.0.1 par défaut et refuse toute requête sans jeton valide, sans `Content-Type: application/json`, avec un en-tête `Origin`, ou dont le corps dépasse 1 MiB.

## 📋 Tâches

- [ ] `WHATSAPP_BRIDGE_ADDR` (défaut `127.0.0.1:8080`)
- [ ] Jeton : `WHATSAPP_BRIDGE_TOKEN` ou `store/bridge_token` généré au 1er lancement (0600)
- [ ] Middleware : Origin → 403, Content-Type → 415, jeton (comparaison à temps constant) → 401, `MaxBytesReader` 1 MiB
- [ ] `http.Server` dédié (mux propre, timeouts)
- [ ] Client Python : en-tête `Authorization: Bearer` lu depuis l'env ou le fichier
- [ ] README : section sécurité

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-01-api-rest-auth]]

## ✅ Definition of Done

- [ ] Tâches terminées
- [ ] Test [[TC-01-api-rest-auth]] **PASS**
- [ ] Commit fait, message préfixé `[Sprint 01]`
- [ ] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [ ] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé.

## 🔗 Liens

- ADR : [[ADR-0004-auth-pont-jeton]]
