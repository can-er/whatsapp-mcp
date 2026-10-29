---
type: sprint
id: sprint-04
title: "Liste blanche de conversations (WHATSAPP_ALLOWED_CHATS)"
status: DONE          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed: 2026-10-29
tags: [sprint]
---

# Sprint 04 — Liste blanche de conversations (WHATSAPP_ALLOWED_CHATS)

> **status:: DONE**  ·  **test:: [[TC-04-allowed-chats]]**

## 🎯 Objectif

Si `WHATSAPP_ALLOWED_CHATS` est défini, les conversations hors liste ne sont ni stockées (temps réel + historique), ni exposées par le serveur MCP, ni joignables par `/api/send` et `/api/download`. Sans la variable : comportement actuel.

## 📋 Tâches

- [x] Go : parsing (JID ou numéro), `Allows(jid)`
- [x] Go : filtre dans `handleMessage` / `handleHistorySync`, refus dans `/api/send` et `/api/download`
- [x] Python : filtre SQL dans toutes les requêtes de lecture
- [x] README

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-04-allowed-chats]]

## ✅ Definition of Done

- [x] Tâches terminées
- [x] Test [[TC-04-allowed-chats]] **PASS**
- [x] Commit fait, message préfixé `[Sprint 04]`
- [x] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [x] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé.
- `2026-10-29` — Liste blanche Go (stockage + API) et Python (vues temporaires) ; TC-04 PASS ; mergé dans main.
- `2026-10-29` — Incident : numéros réels dans les tests → historique réécrit et force-push, garde-fou pre-commit ([[ADR-0007-donnees-de-test-fictives]]).

## 🔗 Liens

- ADR : [[ADR-0007-donnees-de-test-fictives]]
