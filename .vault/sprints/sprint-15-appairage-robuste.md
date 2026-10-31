---
type: sprint
id: sprint-15
title: "Appairage par code robuste : livraison garantie, limite de codes, recul sur 429 (fork)"
status: DONE          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed: 2026-10-29
tags: [sprint]
---

# Sprint 15 — Appairage par code robuste : livraison garantie, limite de codes, recul sur 429 (fork)

> **status:: DONE**  ·  **test:: [[TC-15-appairage-robuste]]**

## 🎯 Objectif

Corriger les problèmes vus à la mise en service (VM ks-agent) : les événements `pairing_code` ne doivent plus se perdre (dossier `events/` disparu, webhook sauté quand l'écriture inbox échoue), et le pont ne doit plus redemander des codes indéfiniment : arrêt après `WHATSAPP_PAIR_MAX_CODES` demandes (5) ou au délai, événement `pairing_stopped`, code de sortie dédié 3 (systemd `RestartPreventExitStatus=3`), recul exponentiel sur 429 rate-overlimit.

## 📋 Tâches

- [x] `deliver` : le webhook est appelé même si l'écriture inbox échoue ; sous-dossiers recréés à chaque écriture
- [x] Journal explicite au démarrage : destination des codes d'appairage (inbox / webhook / terminal seul)
- [x] Limite `WHATSAPP_PAIR_MAX_CODES` (5) comptant chaque demande, réussie ou non ; une seule demande par session
- [x] 429 rate-overlimit : session fermée, recul 5 min doublé jusqu'à 1 h
- [x] Abandon (limite, délai, recul au-delà du délai) : événement `pairing_stopped` + sortie code 3
- [x] README

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-15-appairage-robuste]]

## ✅ Definition of Done

- [x] Tâches terminées
- [x] Test [[TC-15-appairage-robuste]] **PASS**
- [x] Commit fait, message préfixé `[Sprint 15]`
- [x] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [x] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé après le retour de mise en service de ks-support-agent-47 (événements d'appairage perdus après 18:03, 429 rate-overlimit sur le simulateur).
- `2026-10-29` — Livré : commit 55e6ed3 (fix/15-appairage-robuste), mergé dans main. TC-15 PASS.
- `2026-10-29` — Redéployé par ks-support-agent-47 en c2e0cae (RestartPreventExitStatus=3, WHATSAPP_PAIR_MAX_CODES=5, démarrage manuel à l'appairage, pairing_stopped relayé sur Telegram). Cause de l'arrêt de 18:03 confirmée : 429, pas la livraison.
- `2026-10-31` — Les deux ponts sont appairés par code (prod 02:45:13 UTC, 1er code ; simulateur 02:49:38 UTC, 3e code), sans 429 ni abandon.

## 🔗 Liens

- ADR : [[ADR-0006-interface-pont-pipeline]]
