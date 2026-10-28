---
type: adr
id: ADR-0006
title: "Interface pont WhatsApp → pipeline KS (événements v1, API d'envoi)"
status: accepted
created: 2026-10-29
tags: [adr]
---

# ADR-0006 — Interface pont WhatsApp → pipeline KS (événements v1, API d'envoi)

## Contexte

Le fork sert de brique « pont WhatsApp » au pipeline de ks-support-agent-6b (bugs envoyés par
Kevin / Mehdi → tri IA → correctif → Preview Vercel). Le pipeline tourne sur la même VM sous un
autre utilisateur Unix (`ks-agent`, le pont sous `wa-bridge`). Les échanges avec Caner passent par
un bot Telegram : pas de chat de commande WhatsApp. Interface négociée le 2026-10-29 entre
whatsapp-mcp-d8 et ks-support-agent-6b.

## Décision

- **Entrée** : un événement par message temps réel d'un chat autorisé (pas de regroupement côté
  pont ; le pipeline regroupe les rafales), `is_from_me` émis avec le flag.
  - Dossier durable `WHATSAPP_INBOX_DIR` : `events/<YYYYMMDDTHHMMSSZ>_<msgid>.json` écrit
    atomiquement **après** `media/<msgid>.<ext>` (pleine qualité, nom généré). Dossiers 0750,
    fichiers 0640 ; le groupe `ks-inbox` est hérité du bit setgid posé par l'admin. Purge par le
    pont (`WHATSAPP_INBOX_RETENTION_DAYS`, 30).
  - Réveil : POST du même JSON sur `WHATSAPP_EVENT_WEBHOOK_URL` (loopback uniquement),
    `Authorization: Bearer WHATSAPP_EVENT_WEBHOOK_TOKEN`, 3 tentatives.
  - JSON v1 : `version, type, id, chat_jid, chat_name, sender_jid, sender_phone, sender_lid,
    sender_name, is_from_me, timestamp, received_at, text, quoted_message_id,
    media{status, error, type, mime, path, size, sha256, voice, duration_s, original_filename}`.
    `sender_phone` est toujours le vrai numéro (LID résolu), sinon `null` + `sender_lid`.
    `chat_name` / `sender_name` sont informatifs uniquement.
  - v1.1 : `type` = `edit` | `revoke` | `reaction` avec `target_message_id`.
- **Sortie** : `POST /api/send` réservé au jeton d'envoi `WHATSAPP_SEND_TOKEN` (distinct du jeton
  de lecture), idempotent par `client_msg_id`, réponse avec l'id WhatsApp, journal d'audit
  (horodatage, destinataire, client_msg_id, sha256 du texte). Destinataires limités par la liste
  blanche. Envois uniquement depuis le code déterministe de l'orchestrateur.
- **Supervision** : `GET /api/health` (jeton) → `{connected, logged_in, last_event_at}`.
- Tout cela reste dans le fork, **jamais** dans les PR upstream.

## Conséquences

- Pas de polling : inotify sur `events/` + webhook.
- Le pipeline revérifie la liste blanche sur `sender_phone` / `chat_jid`.
- Un processus qui ne peut que lire ne peut pas envoyer.

## Alternatives écartées

- Regroupement des rafales côté pont (`WHATSAPP_EVENT_DEBOUNCE`) — la fenêtre dépend du pipeline (vidéos lentes).
- `WHATSAPP_COMMAND_CHAT` (self-chat de Caner) — Caner a choisi Telegram.

## 🔗 Liens

- Sprints concernés : [[sprint-08-evenements-inbox-webhook]], [[sprint-09-api-envoi-health]], [[sprint-10-evenements-edit-revoke-reaction]]
