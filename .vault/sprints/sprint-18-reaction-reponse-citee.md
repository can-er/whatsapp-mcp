---
type: sprint
id: sprint-18
title: "POST /api/react et réponse citée (reply_to_message_id) dans /api/send (fork)"
status: ONGOING          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-31
closed:
tags: [sprint]
---

# Sprint 18 — POST /api/react et réponse citée (reply_to_message_id) dans /api/send (fork)

> **status:: ONGOING**  ·  **test:: [[TC-18-reaction-reponse-citee]]**

## 🎯 Objectif

Le pipeline peut réagir à un message (`POST /api/react`) et répondre en citant un message (`reply_to_message_id` dans `/api/send`). Le message visé est retrouvé dans le store du pont (même discussion, expéditeur pour les groupes) ; jeton d'envoi, liste blanche, idempotence `client_msg_id` et journal d'audit s'appliquent ; aucun contenu journalisé.

## 📋 Tâches

- [ ] `MessageStore.FindMessage` : message par ID parmi les JID d'une discussion (numéro + LID)
- [ ] `POST /api/react` : jeton d'envoi, liste blanche (403), message inconnu (404), emoji non vide et court, idempotence + audit (`kind: reaction`)
- [ ] `/api/send` + `reply_to_message_id` : message cité d'une autre discussion ou inconnu → 400, média + citation → 400 ; ExtendedTextMessage + ContextInfo (StanzaID, Participant, QuotedMessage)
- [ ] `requireAPIAuth` : `/api/react` exige le jeton d'envoi comme `/api/send`
- [ ] README

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-18-reaction-reponse-citee]]

## ✅ Definition of Done

- [ ] Tâches terminées
- [ ] Test [[TC-18-reaction-reponse-citee]] **PASS**
- [ ] Commit fait, message préfixé `[Sprint 18]`
- [ ] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [ ] Bip Windows émis

## 🧾 Journal

- `2026-10-31` — Sprint créé à la demande de ks-support-agent-47, pour Caner (sprint 13 du pipeline, groupe de l'équipe) : réagir à un message et répondre en le citant.

## 🔗 Liens

- ADR : [[ADR-0006-interface-pont-pipeline]]
