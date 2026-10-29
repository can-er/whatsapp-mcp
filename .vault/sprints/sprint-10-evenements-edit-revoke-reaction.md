---
type: sprint
id: sprint-10
title: "Événements v1.1 : edit, revoke, reaction (fork)"
status: DONE          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed: 2026-10-29
tags: [sprint]
---

# Sprint 10 — Événements v1.1 : edit, revoke, reaction (fork)

> **status:: DONE**  ·  **test:: [[TC-10-evenements-edit-revoke-reaction]]**

## 🎯 Objectif

Les modifications, suppressions et réactions sur les messages des chats autorisés produisent des événements `type: edit|revoke|reaction` avec `target_message_id`. Spécifique au fork (optionnel v1.1).

## 📋 Tâches

- [x] Détection des messages protocolaires (edit / revoke) et des réactions dans `handleMessage`
- [x] Événements inbox + webhook au même format (champ `type`, `target_message_id`, `text` ou `reaction`)
- [x] Mise à jour du contenu en base pour les edits

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-10-evenements-edit-revoke-reaction]]

## ✅ Definition of Done

- [x] Tâches terminées
- [x] Test [[TC-10-evenements-edit-revoke-reaction]] **PASS**
- [x] Commit fait, message préfixé `[Sprint 10]`
- [x] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [x] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé à la demande de ks-support-agent-6b (interface négociée, cf. ADR-0006).
- `2026-10-29` — Événements edit / revoke / reaction + mise à jour des edits en base ; TC-10 PASS ; mergé dans main.

## 🔗 Liens

- ADR : [[ADR-0006-interface-pont-pipeline]]
