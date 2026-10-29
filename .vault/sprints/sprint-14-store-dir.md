---
type: sprint
id: sprint-14
title: "Dossier de données configurable (WHATSAPP_STORE_DIR) (fork)"
status: DONE          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed: 2026-10-29
tags: [sprint]
---

# Sprint 14 — Dossier de données configurable (WHATSAPP_STORE_DIR) (fork)

> **status:: DONE**  ·  **test:: [[TC-14-store-dir]]**

## 🎯 Objectif

Le pont et le serveur MCP lisent `WHATSAPP_STORE_DIR` (défaut `store/` relatif, comportement actuel) pour la session, les messages, le jeton, l'audit et les médias, ce qui permet plusieurs instances avec `StateDirectory=` distinct.

## 📋 Tâches

- [x] `storePath()` côté Go, utilisé partout (bases, jeton, audit, médias)
- [x] Python : `MESSAGES_DB_PATH` / `BRIDGE_TOKEN_PATH` dérivés de `WHATSAPP_STORE_DIR`
- [x] README

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-14-store-dir]]

## ✅ Definition of Done

- [x] Tâches terminées
- [x] Test [[TC-14-store-dir]] **PASS**
- [x] Commit fait, message préfixé `[Sprint 14]`
- [x] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [x] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé à la demande de ks-support-agent-6b (unités systemd avec StateDirectory= distinct).
- `2026-10-29` — WHATSAPP_STORE_DIR (Go + Python) ; TC-14 PASS ; mergé dans main.

## 🔗 Liens

- ADR : [[ADR-0006-interface-pont-pipeline]]
