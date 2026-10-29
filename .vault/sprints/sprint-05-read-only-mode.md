---
type: sprint
id: sprint-05
title: "Mode lecture seule du serveur MCP (PR upstream)"
status: DONE          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed: 2026-10-29
tags: [sprint]
---

# Sprint 05 — Mode lecture seule du serveur MCP (PR upstream)

> **status:: DONE**  ·  **test:: [[TC-05-read-only-mode]]**

## 🎯 Objectif

Avec `WHATSAPP_MCP_READ_ONLY=1`, les outils `send_message`, `send_file` et `send_audio_message` ne sont pas enregistrés ; le README avertit explicitement de l'injection de prompt via les messages reçus. **PR vers lharries/whatsapp-mcp.**

## 📋 Tâches

- [x] Décorateur conditionnel pour les outils d'envoi
- [x] README : mode lecture seule + avertissement prompt injection
- [x] Test Python
- [x] Branche basée sur `upstream/main`, PR upstream

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-05-read-only-mode]]

## ✅ Definition of Done

- [x] Tâches terminées
- [x] Test [[TC-05-read-only-mode]] **PASS**
- [x] Commit fait, message préfixé `[Sprint 05]`
- [x] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [x] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé.
- `2026-10-29` — Mode lecture seule ; TC-05 PASS ; mergé dans main ; PR upstream #388.

## 🔗 Liens

- ADR : [[ADR-0001-perimetre-upstream]]
