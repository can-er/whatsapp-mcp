---
type: sprint
id: sprint-09
title: "API d'envoi sécurisée (jeton dédié, idempotence, audit) + /api/health (fork)"
status: TODO          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed:
tags: [sprint]
---

# Sprint 09 — API d'envoi sécurisée (jeton dédié, idempotence, audit) + /api/health (fork)

> **status:: TODO**  ·  **test:: [[TC-09-api-envoi-health]]**

## 🎯 Objectif

`/api/send` exige un jeton d'envoi distinct (`WHATSAPP_SEND_TOKEN`), est idempotent via `client_msg_id` (renvoie l'id WhatsApp du message), journalise chaque envoi ; `GET /api/health` expose `{connected, logged_in, last_event_at}`. Spécifique au fork.

## 📋 Tâches

- [ ] Jeton d'envoi distinct : si défini, seul lui ouvre `/api/send` ; le jeton du pont reste valable pour download/health
- [ ] `client_msg_id` : table d'idempotence SQLite, renvoi → même réponse sans second envoi ; réponse avec `message_id`
- [ ] Journal d'audit JSON Lines (horodatage, destinataire, client_msg_id, sha256 du texte, résultat), 0600
- [ ] `GET /api/health` (jeton requis)
- [ ] README

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-09-api-envoi-health]]

## ✅ Definition of Done

- [ ] Tâches terminées
- [ ] Test [[TC-09-api-envoi-health]] **PASS**
- [ ] Commit fait, message préfixé `[Sprint 09]`
- [ ] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [ ] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé à la demande de ks-support-agent-6b (interface négociée, cf. ADR-0006).

## 🔗 Liens

- ADR : [[ADR-0006-interface-pont-pipeline]]
