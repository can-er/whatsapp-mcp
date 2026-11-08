---
type: sprint
id: sprint-20
title: "Numéro de l'expéditeur stocké dans messages.db à côté du LID (colonne sender_phone, fork)"
status: DONE          # TODO | ONGOING | DONE
owner: CK
created: 2026-11-08
closed: 2026-11-08
tags: [sprint]
---

# Sprint 20 — Numéro de l'expéditeur stocké dans messages.db à côté du LID (colonne sender_phone, fork)

> **status:: DONE**  ·  **test:: [[TC-20-numero-expediteur]]**

## 🎯 Objectif

Chaque message stocké porte le numéro de l'expéditeur dans une nouvelle colonne `sender_phone` (résolu depuis le LID par la table LID de whatsmeow quand l'expéditeur n'est connu que par son LID). La colonne `sender` reste inchangée (elle sert à construire les clés de message des réactions et réponses citées). Les lignes existantes sont complétées au démarrage quand le numéro est connu.

## 📋 Tâches

- [x] Migration : colonne `sender_phone` ajoutée si absente (base existante conservée)
- [x] Messages en direct : numéro tiré de Sender / SenderAlt, ou de la table LID ; propre numéro pour `is_from_me`
- [x] Synchronisation d'historique : même résolution depuis le participant
- [x] Rattrapage au démarrage des lignes sans `sender_phone`, sans jamais deviner : un identifiant ambigu sans correspondance reste NULL
- [x] README (requête par numéro : `sender_phone = ?`)

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-20-numero-expediteur]]

## ✅ Definition of Done

- [x] Tâches terminées
- [x] Test [[TC-20-numero-expediteur]] **PASS**
- [x] Commit fait, message préfixé `[Sprint 20]`
- [x] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [x] Bip Windows émis

## 🧾 Journal

- `2026-11-08` — Sprint créé à la demande de la session KS Esthetic (« Higgsfield CLI et skills »), pour Caner, après lecture du groupe de l'équipe via le pont de production : `messages.sender` contient le LID de l'auteur dans le groupe, un filtre par numéro sur messages.db le rate.

- `2026-11-08` — Corrigé sur `fix/20-numero-expediteur` (4e19227), PR [#2](https://github.com/can-er/whatsapp-mcp/pull/2) empilée sur #1. Choix : nouvelle colonne `sender_phone` plutôt que réécrire `sender`, car les réactions et réponses citées (sprint 18) construisent la clé du message à partir de `sender` (LID nécessaire dans les groupes LID). Le serveur MCP Python filtre encore sur `sender`.

- `2026-11-08` — Fusionné dans `main` sur accord de Caner (gocheck Linux PASS sur `main`). Déploiement sur la VM non fait depuis ce poste : 10.10.10.90 injoignable en SSH.
- `2026-11-08` — Déployé sur la VM vers 22:26 UTC par la session KS Esthetic, avec l'accord de Caner. Binaire compilé depuis `main` (code identique à d8755dc). Sauvegarde des bases et de l'ancien binaire faite avant. Les deux instances se sont reconnectées sans nouveau scan. `sender_phone` rempli rétroactivement sur toutes les lignes existantes (804 sur le pont réel). `/api/download` réussit sur les deux images qui renvoyaient 403 : copie de l'inbox. Reste à observer : la légende du prochain média reçu avec du texte.

## 🔗 Liens

- ADR : [[ADR-0006-interface-pont-pipeline]]
