---
type: sprint
id: sprint-08
title: "Événements entrants : dossier inbox + webhook local (fork)"
status: DONE          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed: 2026-10-29
tags: [sprint]
---

# Sprint 08 — Événements entrants : dossier inbox + webhook local (fork)

> **status:: DONE**  ·  **test:: [[TC-08-evenements-inbox-webhook]]**

## 🎯 Objectif

Chaque message temps réel d'un chat autorisé produit un événement JSON v1 écrit atomiquement dans `WHATSAPP_INBOX_DIR/events/` (média complet téléchargé avant, dans `media/`, nom généré) et notifié par POST au webhook local `WHATSAPP_EVENT_WEBHOOK_URL`. Spécifique au fork (pipeline KS).

## 📋 Tâches

- [x] Schéma JSON v1 ([[ADR-0006-interface-pont-pipeline]]) : id, chat_jid, sender_jid/phone/lid, is_from_me, timestamp, received_at, text, quoted_message_id, media{status, type, mime, path, size, sha256, voice, duration_s, original_filename, error}
- [x] Téléchargement du média avant l'écriture de l'événement ; `media.status=failed` + `error` si échec
- [x] Écriture atomique (tmp + rename), fichiers 0640, dossiers 0750 (groupe hérité du setgid posé par l'admin)
- [x] Webhook : URL loopback uniquement, `Authorization: Bearer WHATSAPP_EVENT_WEBHOOK_TOKEN`, timeout 10 s, 3 tentatives avec backoff, asynchrone
- [x] LID → numéro (sender_phone réel ou null + sender_lid)
- [x] Purge `WHATSAPP_INBOX_RETENTION_DAYS` (30 par défaut)
- [x] README (section fork)

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-08-evenements-inbox-webhook]]

## ✅ Definition of Done

- [x] Tâches terminées
- [x] Test [[TC-08-evenements-inbox-webhook]] **PASS**
- [x] Commit fait, message préfixé `[Sprint 08]`
- [x] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [x] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé à la demande de ks-support-agent-6b (interface négociée, cf. ADR-0006).
- `2026-10-29` — Événements inbox + webhook (ADR-0006) ; TC-08 PASS ; mergé dans main.

## 🔗 Liens

- ADR : [[ADR-0006-interface-pont-pipeline]]
