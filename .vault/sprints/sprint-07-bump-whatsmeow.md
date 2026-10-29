---
type: sprint
id: sprint-07
title: "Mise à jour de whatsmeow"
status: DONE          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed: 2026-10-29
tags: [sprint]
---

# Sprint 07 — Mise à jour de whatsmeow

> **status:: DONE**  ·  **test:: [[TC-07-bump-whatsmeow]]**

## 🎯 Objectif

Le pont compile et passe `go vet` / `go test` avec une version récente de `go.mau.fi/whatsmeow` (API avec `context.Context`).

## 📋 Tâches

- [x] `go get go.mau.fi/whatsmeow@latest` + `go mod tidy`
- [x] Adapter les appels (sqlstore.New, GetFirstDevice, Download, GetGroupInfo, Contacts…)
- [x] Test d'intégration CGO : ouverture du sqlstore en mémoire + création d'un device
- [x] README (version de Go requise)

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-07-bump-whatsmeow]]

## ✅ Definition of Done

- [x] Tâches terminées
- [x] Test [[TC-07-bump-whatsmeow]] **PASS**
- [x] Commit fait, message préfixé `[Sprint 07]`
- [x] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [x] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé.
- `2026-10-29` — whatsmeow à jour (API context), liste blanche LID-aware ; TC-07 PASS ; mergé dans main.

## 🔗 Liens

- ADR : —
