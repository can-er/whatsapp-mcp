---
type: sprint
id: sprint-15
title: "Appairage par code robuste : livraison garantie, limite de codes, recul sur 429 (fork)"
status: ONGOING          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed:
tags: [sprint]
---

# Sprint 15 — Appairage par code robuste : livraison garantie, limite de codes, recul sur 429 (fork)

> **status:: ONGOING**  ·  **test:: [[TC-15-appairage-robuste]]**

## 🎯 Objectif

Corriger les problèmes vus à la mise en service (VM ks-agent) : les événements `pairing_code` ne doivent plus se perdre (dossier `events/` disparu, webhook sauté quand l'écriture inbox échoue), et le pont ne doit plus redemander des codes indéfiniment : arrêt après `WHATSAPP_PAIR_MAX_CODES` demandes (5) ou au délai, événement `pairing_stopped`, code de sortie dédié 3 (systemd `RestartPreventExitStatus=3`), recul exponentiel sur 429 rate-overlimit.

## 📋 Tâches

- [ ] `deliver` : le webhook est appelé même si l'écriture inbox échoue ; sous-dossiers recréés à chaque écriture
- [ ] Journal explicite au démarrage : destination des codes d'appairage (inbox / webhook / terminal seul)
- [ ] Limite `WHATSAPP_PAIR_MAX_CODES` (5) comptant chaque demande, réussie ou non ; une seule demande par session
- [ ] 429 rate-overlimit : session fermée, recul 5 min doublé jusqu'à 1 h
- [ ] Abandon (limite, délai, recul au-delà du délai) : événement `pairing_stopped` + sortie code 3
- [ ] README

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-15-appairage-robuste]]

## ✅ Definition of Done

- [ ] Tâches terminées
- [ ] Test [[TC-15-appairage-robuste]] **PASS**
- [ ] Commit fait, message préfixé `[Sprint 15]`
- [ ] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [ ] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé après le retour de mise en service de ks-support-agent-47 (événements d'appairage perdus après 18:03, 429 rate-overlimit sur le simulateur).

## 🔗 Liens

- ADR : [[ADR-0006-interface-pont-pipeline]]
