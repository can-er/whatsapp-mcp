---
type: sprint
id: sprint-06
title: "Permissions restreintes sur store/ (PR upstream)"
status: DONE          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed: 2026-10-29
tags: [sprint]
---

# Sprint 06 — Permissions restreintes sur store/ (PR upstream)

> **status:: DONE**  ·  **test:: [[TC-06-store-permissions]]**

## 🎯 Objectif

Le dossier `store/` (clés de session + messages) est créé en 0700 (resserré s'il existait), les bases et médias en 0600 ; README : note Windows. **PR vers lharries/whatsapp-mcp.**

## 📋 Tâches

- [x] `ensurePrivateDir` / `ensurePrivateFile`
- [x] Appliquer à `store/`, `whatsapp.db`, `messages.db`, dossiers et fichiers médias
- [x] README (Windows : ACL du profil, icacls)
- [x] Branche basée sur `upstream/main`, PR upstream

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-06-store-permissions]]

## ✅ Definition of Done

- [x] Tâches terminées
- [x] Test [[TC-06-store-permissions]] **PASS**
- [x] Commit fait, message préfixé `[Sprint 06]`
- [x] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [x] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé.
- `2026-10-29` — Permissions 0700/0600 (store, bases, WAL, médias) ; TC-06 PASS ; mergé dans main ; PR upstream #389.

## 🔗 Liens

- ADR : [[ADR-0001-perimetre-upstream]], [[ADR-0003-tests-go-docker]]
