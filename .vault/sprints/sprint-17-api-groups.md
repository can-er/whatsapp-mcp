---
type: sprint
id: sprint-17
title: "GET /api/groups : découverte des JID de groupes, métadonnées seulement (fork)"
status: ONGOING          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-31
closed:
tags: [sprint]
---

# Sprint 17 — GET /api/groups : découverte des JID de groupes, métadonnées seulement (fork)

> **status:: ONGOING**  ·  **test:: [[TC-17-api-groups]]**

## 🎯 Objectif

`GET /api/groups` (jeton du pont, pas le jeton d'envoi) renvoie les groupes rejoints sous la forme `{jid, name, participants_count}`, triés par nom, sans messages ni numéros de participants, pour trouver le JID à mettre dans `WHATSAPP_ALLOWED_CHATS`.

## 📋 Tâches

- [ ] Handler `/api/groups` (GET seulement), source injectable via `apiOptions.groups` (`client.GetJoinedGroups` en production)
- [ ] Réponse limitée à jid / name / participants_count ; 503 si non connecté, 502 si WhatsApp échoue
- [ ] README

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-17-api-groups]]

## ✅ Definition of Done

- [ ] Tâches terminées
- [ ] Test [[TC-17-api-groups]] **PASS**
- [ ] Commit fait, message préfixé `[Sprint 17]`
- [ ] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [ ] Bip Windows émis

## 🧾 Journal

- `2026-10-31` — Sprint créé à la demande de ks-support-agent-47 (sprint 13 du pipeline) : obtenir le JID du groupe « KS Support & improvements » sans stocker de discussion hors liste blanche.

## 🔗 Liens

- ADR : [[ADR-0006-interface-pont-pipeline]]
