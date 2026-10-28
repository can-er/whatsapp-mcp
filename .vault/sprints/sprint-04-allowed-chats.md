---
type: sprint
id: sprint-04
title: "Liste blanche de conversations (WHATSAPP_ALLOWED_CHATS)"
status: TODO          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed:
tags: [sprint]
---

# Sprint 04 — Liste blanche de conversations (WHATSAPP_ALLOWED_CHATS)

> **status:: TODO**  ·  **test:: [[TC-04-allowed-chats]]**

## 🎯 Objectif

Si `WHATSAPP_ALLOWED_CHATS` est défini, les conversations hors liste ne sont ni stockées (temps réel + historique), ni exposées par le serveur MCP, ni joignables par `/api/send` et `/api/download`. Sans la variable : comportement actuel.

## 📋 Tâches

- [ ] Go : parsing (JID ou numéro), `Allows(jid)`
- [ ] Go : filtre dans `handleMessage` / `handleHistorySync`, refus dans `/api/send` et `/api/download`
- [ ] Python : filtre SQL dans toutes les requêtes de lecture
- [ ] README

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-04-allowed-chats]]

## ✅ Definition of Done

- [ ] Tâches terminées
- [ ] Test [[TC-04-allowed-chats]] **PASS**
- [ ] Commit fait, message préfixé `[Sprint 04]`
- [ ] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [ ] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé.

## 🔗 Liens

- ADR : —
