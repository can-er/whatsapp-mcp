---
type: sprint
id: sprint-19
title: "Légende des images, vidéos et documents stockée dans messages.db (fork)"
status: ONGOING          # TODO | ONGOING | DONE
owner: CK
created: 2026-11-08
closed:
tags: [sprint]
---

# Sprint 19 — Légende des images, vidéos et documents stockée dans messages.db (fork)

> **status:: ONGOING**  ·  **test:: [[TC-19-legende-medias]]**

## 🎯 Objectif

La légende (caption) d'une image, d'une vidéo ou d'un document est enregistrée dans `messages.content`, en direct comme à la synchronisation d'historique, comme le fait déjà l'événement exporté. Aujourd'hui `content` reste vide pour ces messages.

## 📋 Tâches

- [ ] `extractTextContent` : texte, texte étendu, puis légende image / vidéo / document (même règle que `messageText` des événements)
- [ ] Synchronisation d'historique : même extraction (la légende n'y était pas lue non plus)
- [ ] README

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-19-legende-medias]]

## ✅ Definition of Done

- [ ] Tâches terminées
- [ ] Test [[TC-19-legende-medias]] **PASS**
- [ ] Commit fait, message préfixé `[Sprint 19]`
- [ ] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [ ] Bip Windows émis

## 🧾 Journal

- `2026-11-08` — Sprint créé à la demande de la session KS Esthetic (« Higgsfield CLI et skills »), pour Caner, après lecture du groupe de l'équipe via le pont de production : une image avec légende du groupe a `content` vide dans messages.db alors que l'événement JSON contient le texte.

## 🔗 Liens

- ADR : [[ADR-0006-interface-pont-pipeline]]
