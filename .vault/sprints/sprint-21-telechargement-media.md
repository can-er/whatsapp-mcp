---
type: sprint
id: sprint-21
title: "/api/download : chemin d'accès direct complet, copie de l'inbox et renouvellement du média (fork)"
status: DONE          # TODO | ONGOING | DONE
owner: CK
created: 2026-11-08
closed: 2026-11-08
tags: [sprint]
---

# Sprint 21 — /api/download : chemin d'accès direct complet, copie de l'inbox et renouvellement du média (fork)

> **status:: DONE**  ·  **test:: [[TC-21-telechargement-media]]**

## 🎯 Objectif

`POST /api/download` ne doit plus renvoyer 403 pour un média encore valide. Cause probable : le chemin direct était reconstruit depuis l'URL en coupant ses paramètres signés (`oh`, `oe`) ; il est désormais stocké tel que WhatsApp le donne (`direct_path`). Si le média est déjà dans l'inbox (même empreinte SHA-256), ce fichier est renvoyé. Si le CDN répond 403/404/410, le pont demande un renouvellement au téléphone (media retry de whatsmeow) puis retélécharge.

## 📋 Tâches

- [x] Colonne `direct_path` (migration) remplie à la réception et à la synchronisation ; repli sur l'URL en gardant la requête pour les lignes anciennes
- [x] Recherche dans `WHATSAPP_INBOX_DIR/media/<id>.*` avec vérification de l'empreinte `file_sha256`
- [x] Sur 403/404/410 : `SendMediaRetryReceipt`, attente de l'événement `MediaRetry` (délai borné), nouveau chemin direct stocké, nouveau téléchargement
- [x] README

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-21-telechargement-media]]

## ✅ Definition of Done

- [x] Tâches terminées
- [x] Test [[TC-21-telechargement-media]] **PASS**
- [x] Commit fait, message préfixé `[Sprint 21]`
- [x] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [x] Bip Windows émis

## 🧾 Journal

- `2026-11-08` — Sprint créé à la demande de la session KS Esthetic (« Higgsfield CLI et skills »), pour Caner, après lecture du groupe de l'équipe via le pont de production : `/api/download` renvoie « download failed with status code 403 » pour deux images reçues deux heures plus tôt, déjà présentes dans l'inbox.

- `2026-11-08` — Corrigé sur `fix/21-telechargement-media` (ed16824), PR [#3](https://github.com/can-er/whatsapp-mcp/pull/3) empilée sur #2. L'inbox renvoie son propre chemin (pas de copie dans `store/`), ce qui respecte la purge de rétention de l'inbox.

## 🔗 Liens

- ADR : [[ADR-0006-interface-pont-pipeline]]
