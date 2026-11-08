---
type: sprint
id: sprint-19
title: "Légende des images, vidéos et documents stockée dans messages.db (fork)"
status: DONE          # TODO | ONGOING | DONE
owner: CK
created: 2026-11-08
closed: 2026-11-08
tags: [sprint]
---

# Sprint 19 — Légende des images, vidéos et documents stockée dans messages.db (fork)

> **status:: DONE**  ·  **test:: [[TC-19-legende-medias]]**

## 🎯 Objectif

La légende (caption) d'une image, d'une vidéo ou d'un document est enregistrée dans `messages.content`, en direct comme à la synchronisation d'historique, comme le fait déjà l'événement exporté. Aujourd'hui `content` reste vide pour ces messages.

## 📋 Tâches

- [x] `extractTextContent` : texte, texte étendu, puis légende image / vidéo / document (même règle que `messageText` des événements)
- [x] Synchronisation d'historique : même extraction (la légende n'y était pas lue non plus)
- [x] README

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-19-legende-medias]]

## ✅ Definition of Done

- [x] Tâches terminées
- [x] Test [[TC-19-legende-medias]] **PASS**
- [x] Commit fait, message préfixé `[Sprint 19]`
- [x] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [x] Bip Windows émis

## 🧾 Journal

- `2026-11-08` — Sprint créé à la demande de la session KS Esthetic (« Higgsfield CLI et skills »), pour Caner, après lecture du groupe de l'équipe via le pont de production : une image avec légende du groupe a `content` vide dans messages.db alors que l'événement JSON contient le texte.

- `2026-11-08` — Corrigé sur `fix/19-legende-medias` (cb12a33), PR [#1](https://github.com/can-er/whatsapp-mcp/pull/1) ouverte sur le fork, merge par Caner. Les lignes déjà stockées gardent `content` vide.

- `2026-11-08` — Fusionné dans `main` sur accord de Caner (gocheck Linux PASS sur `main`). Déploiement sur la VM non fait depuis ce poste : 10.10.10.90 injoignable en SSH.
- `2026-11-08` — Déployé sur la VM vers 22:26 UTC par la session KS Esthetic, avec l'accord de Caner. Binaire compilé depuis `main` (code identique à d8755dc). Sauvegarde des bases et de l'ancien binaire faite avant. Les deux instances se sont reconnectées sans nouveau scan. `sender_phone` rempli rétroactivement sur toutes les lignes existantes (804 sur le pont réel). `/api/download` réussit sur les deux images qui renvoyaient 403 : copie de l'inbox. Reste à observer : la légende du prochain média reçu avec du texte.

## 🔗 Liens

- ADR : [[ADR-0006-interface-pont-pipeline]]
