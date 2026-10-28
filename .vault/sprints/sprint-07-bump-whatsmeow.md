---
type: sprint
id: sprint-07
title: "Mise à jour de whatsmeow"
status: TODO          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed:
tags: [sprint]
---

# Sprint 07 — Mise à jour de whatsmeow

> **status:: TODO**  ·  **test:: [[TC-07-bump-whatsmeow]]**

## 🎯 Objectif

Le pont compile et passe `go vet` / `go test` avec une version récente de `go.mau.fi/whatsmeow` (API avec `context.Context`).

## 📋 Tâches

- [ ] `go get go.mau.fi/whatsmeow@latest` + `go mod tidy`
- [ ] Adapter les appels (sqlstore.New, GetFirstDevice, Download, GetGroupInfo, Contacts…)
- [ ] Test d'intégration CGO : ouverture du sqlstore en mémoire + création d'un device
- [ ] README (version de Go requise)

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-07-bump-whatsmeow]]

## ✅ Definition of Done

- [ ] Tâches terminées
- [ ] Test [[TC-07-bump-whatsmeow]] **PASS**
- [ ] Commit fait, message préfixé `[Sprint 07]`
- [ ] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [ ] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé.

## 🔗 Liens

- ADR : —
