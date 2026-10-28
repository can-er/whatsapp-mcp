---
type: sprint
id: sprint-00
title: "Initialisation (fork, vault, outillage)"
status: DONE          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed: 2026-10-29
tags: [sprint]
---

# Sprint 00 — Initialisation (fork, vault, outillage)

> **status:: DONE**  ·  **test:: [[TC-00-initialisation]]**

## 🎯 Objectif

Disposer du fork, du vault de pilotage et d'une chaîne de test reproductible (Go via Docker/CGO, Python via uv).

## 📋 Tâches

- [x] Vérifier SECURITY.md / private vulnerability reporting de l'upstream
- [x] Recenser les issues/PR existantes sur les failles
- [x] Fork `can-er/whatsapp-mcp`, remotes `origin` (fork) / `upstream`
- [x] Installer Go, script `gocheck.sh` (Docker, CGO)
- [x] Créer `.vault/`, `AGENTS.md`, stub `CLAUDE.md`

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-00-initialisation]]

## ✅ Definition of Done

- [x] Tâches terminées
- [x] Test [[TC-00-initialisation]] **PASS**
- [x] Commit fait, message préfixé `[Sprint 00]`
- [x] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [x] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé.
- `2026-10-29` — Fork, vault, AGENTS.md, gocheck (Docker/CGO) en place ; TC-00 PASS.

## 🔗 Liens

- ADR : [[ADR-0001-perimetre-upstream]], [[ADR-0002-organisation-git]], [[ADR-0003-tests-go-docker]]
