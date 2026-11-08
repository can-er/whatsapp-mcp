---
type: sprint
id: sprint-20
title: "Numéro de l'expéditeur stocké dans messages.db à côté du LID (colonne sender_phone, fork)"
status: TODO          # TODO | ONGOING | DONE
owner: CK
created: 2026-11-08
closed:
tags: [sprint]
---

# Sprint 20 — Numéro de l'expéditeur stocké dans messages.db à côté du LID (colonne sender_phone, fork)

> **status:: TODO**  ·  **test:: [[TC-20-numero-expediteur]]**

## 🎯 Objectif

Chaque message stocké porte le numéro de l'expéditeur dans une nouvelle colonne `sender_phone` (résolu depuis le LID par la table LID de whatsmeow quand l'expéditeur n'est connu que par son LID). La colonne `sender` reste inchangée (elle sert à construire les clés de message des réactions et réponses citées). Les lignes existantes sont complétées au démarrage quand le numéro est connu.

## 📋 Tâches

- [ ] Migration : colonne `sender_phone` ajoutée si absente (base existante conservée)
- [ ] Messages en direct : numéro tiré de Sender / SenderAlt, ou de la table LID ; propre numéro pour `is_from_me`
- [ ] Synchronisation d'historique : même résolution depuis le participant
- [ ] Rattrapage au démarrage des lignes sans `sender_phone`, sans jamais deviner : un identifiant ambigu sans correspondance reste NULL
- [ ] README (requête par numéro : `sender_phone = ?`)

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-20-numero-expediteur]]

## ✅ Definition of Done

- [ ] Tâches terminées
- [ ] Test [[TC-20-numero-expediteur]] **PASS**
- [ ] Commit fait, message préfixé `[Sprint 20]`
- [ ] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [ ] Bip Windows émis

## 🧾 Journal

- `2026-11-08` — Sprint créé à la demande de la session KS Esthetic (« Higgsfield CLI et skills »), pour Caner, après lecture du groupe de l'équipe via le pont de production : `messages.sender` contient le LID de l'auteur dans le groupe, un filtre par numéro sur messages.db le rate.

## 🔗 Liens

- ADR : [[ADR-0006-interface-pont-pipeline]]
