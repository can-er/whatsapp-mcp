---
type: sprint
id: sprint-05
title: "Mode lecture seule du serveur MCP (PR upstream)"
status: TODO          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed:
tags: [sprint]
---

# Sprint 05 — Mode lecture seule du serveur MCP (PR upstream)

> **status:: TODO**  ·  **test:: [[TC-05-read-only-mode]]**

## 🎯 Objectif

Avec `WHATSAPP_MCP_READ_ONLY=1`, les outils `send_message`, `send_file` et `send_audio_message` ne sont pas enregistrés ; le README avertit explicitement de l'injection de prompt via les messages reçus. **PR vers lharries/whatsapp-mcp.**

## 📋 Tâches

- [ ] Décorateur conditionnel pour les outils d'envoi
- [ ] README : mode lecture seule + avertissement prompt injection
- [ ] Test Python
- [ ] Branche basée sur `upstream/main`, PR upstream

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-05-read-only-mode]]

## ✅ Definition of Done

- [ ] Tâches terminées
- [ ] Test [[TC-05-read-only-mode]] **PASS**
- [ ] Commit fait, message préfixé `[Sprint 05]`
- [ ] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [ ] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé.

## 🔗 Liens

- ADR : [[ADR-0001-perimetre-upstream]]
