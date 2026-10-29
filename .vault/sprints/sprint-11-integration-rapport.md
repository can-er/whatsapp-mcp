---
type: sprint
id: sprint-11
title: "Intégration sur le fork, PR upstream et rapport"
status: DONE          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed: 2026-10-29
tags: [sprint]
---

# Sprint 11 — Intégration sur le fork, PR upstream et rapport

> **status:: DONE**  ·  **test:: [[TC-11-integration-rapport]]**

## 🎯 Objectif

`main` du fork contient les 7 correctifs et l'interface pipeline KS ; PR upstream ouvertes pour 05 et 06 ; rapport envoyé à la session « Prise de rendez-vous depuis l'agenda (2) ».

## 📋 Tâches

- [x] Merger fix/01→10 dans `main` (résolution de conflits)
- [x] gocheck + tests Python sur `main`
- [x] Push du fork
- [x] PR upstream 05 et 06
- [x] Rapport final

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-11-integration-rapport]]

## ✅ Definition of Done

- [x] Tâches terminées
- [x] Test [[TC-11-integration-rapport]] **PASS**
- [x] Commit fait, message préfixé `[Sprint 11]`
- [x] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [x] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé.
- `2026-10-29` — Intégration : main vert (32 Go + 30 Python), PR #388/#389, rapport envoyé.

## 🔗 Liens

- ADR : [[ADR-0002-organisation-git]]
