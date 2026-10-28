---
type: adr
id: ADR-0004
title: "Authentification du pont par jeton Bearer partagé"
status: accepted
created: 2026-10-29
tags: [adr]
---

# ADR-0004 — Authentification du pont par jeton Bearer partagé

## Contexte

Le serveur MCP Python et le pont Go tournent sur la même machine mais dans des processus distincts, lancés séparément. Il faut un secret partagé sans configuration obligatoire.

## Décision

Le pont lit `WHATSAPP_BRIDGE_TOKEN` ; à défaut il crée `store/bridge_token` (32 octets aléatoires, hex, mode 0600) au premier lancement. Le client Python lit la même variable ou le même fichier (chemin relatif identique à celui de `messages.db`). En-tête `Authorization: Bearer <jeton>`, comparaison à temps constant. Adresse d'écoute `WHATSAPP_BRIDGE_ADDR`, défaut `127.0.0.1:8080`.

## Conséquences

- Aucune étape de configuration supplémentaire pour l'installation standard.
- Requêtes CSRF bloquées doublement (en-tête Authorization non « simple » + refus de `Origin` + JSON obligatoire).

## Alternatives écartées

- Jeton obligatoire en variable (comme #387) — casse l'installation par défaut.
- Socket Unix — non portable sous Windows.

## 🔗 Liens

- Sprints concernés : [[sprint-01-api-rest-auth]]
