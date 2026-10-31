---
type: sprint
id: sprint-18
title: "POST /api/react et réponse citée (reply_to_message_id) dans /api/send (fork)"
status: DONE          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-31
closed: 2026-10-31
tags: [sprint]
---

# Sprint 18 — POST /api/react et réponse citée (reply_to_message_id) dans /api/send (fork)

> **status:: DONE**  ·  **test:: [[TC-18-reaction-reponse-citee]]**

## 🎯 Objectif

Le pipeline peut réagir à un message (`POST /api/react`) et répondre en citant un message (`reply_to_message_id` dans `/api/send`). Le message visé est retrouvé dans le store du pont (même discussion, expéditeur pour les groupes) ; jeton d'envoi, liste blanche, idempotence `client_msg_id` et journal d'audit s'appliquent ; aucun contenu journalisé.

## 📋 Tâches

- [x] `MessageStore.FindMessage` : message par ID parmi les JID d'une discussion (numéro + LID)
- [x] `POST /api/react` : jeton d'envoi, liste blanche (403), message inconnu (404), emoji non vide et court, idempotence + audit (`kind: reaction`)
- [x] `/api/send` + `reply_to_message_id` : message cité d'une autre discussion ou inconnu → 400, média + citation → 400 ; ExtendedTextMessage + ContextInfo (StanzaID, Participant, QuotedMessage)
- [x] `requireAPIAuth` : `/api/react` exige le jeton d'envoi comme `/api/send`
- [x] README

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-18-reaction-reponse-citee]]

## ✅ Definition of Done

- [x] Tâches terminées
- [x] Test [[TC-18-reaction-reponse-citee]] **PASS**
- [x] Commit fait, message préfixé `[Sprint 18]`
- [x] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [x] Bip Windows émis

## 🧾 Journal

- `2026-10-31` — Sprint créé à la demande de ks-support-agent-47, pour Caner (sprint 13 du pipeline, groupe de l'équipe) : réagir à un message et répondre en le citant.
- `2026-10-31` — Livré : commit afdbd6f (fix/18-reaction-reponse-citee), mergé dans main. TC-18 PASS.
- `2026-10-31` — Validé en réel en discussion directe ; activé côté pipeline (KS_GROUP_FEEDBACK=1). Cas groupe à confirmer.

## 🔗 Liens

- ADR : [[ADR-0006-interface-pont-pipeline]]
