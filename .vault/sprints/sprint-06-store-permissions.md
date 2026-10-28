---
type: sprint
id: sprint-06
title: "Permissions restreintes sur store/ (PR upstream)"
status: TODO          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed:
tags: [sprint]
---

# Sprint 06 — Permissions restreintes sur store/ (PR upstream)

> **status:: TODO**  ·  **test:: [[TC-06-store-permissions]]**

## 🎯 Objectif

Le dossier `store/` (clés de session + messages) est créé en 0700 (resserré s'il existait), les bases et médias en 0600 ; README : note Windows. **PR vers lharries/whatsapp-mcp.**

## 📋 Tâches

- [ ] `ensurePrivateDir` / `ensurePrivateFile`
- [ ] Appliquer à `store/`, `whatsapp.db`, `messages.db`, dossiers et fichiers médias
- [ ] README (Windows : ACL du profil, icacls)
- [ ] Branche basée sur `upstream/main`, PR upstream

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-06-store-permissions]]

## ✅ Definition of Done

- [ ] Tâches terminées
- [ ] Test [[TC-06-store-permissions]] **PASS**
- [ ] Commit fait, message préfixé `[Sprint 06]`
- [ ] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [ ] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé.

## 🔗 Liens

- ADR : [[ADR-0001-perimetre-upstream]], [[ADR-0003-tests-go-docker]]
