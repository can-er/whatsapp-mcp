---
type: sprint
id: sprint-16
title: "Confidentialité : aucun texte dans les journaux, liste blanche stricte, purge hors liste (fork + correctif)"
status: DONE          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-31
closed: 2026-10-31
tags: [sprint]
---

# Sprint 16 — Confidentialité : aucun texte dans les journaux, liste blanche stricte, purge hors liste (fork + correctif)

> **status:: DONE**  ·  **test:: [[TC-16-confidentialite-journaux-stockage]]**

## 🎯 Objectif

Le pont ne journalise plus jamais le texte d'un message (ni noms de contacts en INFO). Une valeur de `WHATSAPP_ALLOWED_CHATS` non vide mais invalide empêche le démarrage au lieu de tout autoriser. `WHATSAPP_PURGE_DISALLOWED_CHATS=true` supprime au démarrage les discussions, messages et médias stockés hors liste blanche, puis compacte la base (VACUUM).

## 📋 Tâches

- [x] Retirer le texte des messages de tous les journaux (direct + synchro d'historique) ; métadonnées en DEBUG seulement
- [x] Noms de contacts / groupes : journalisés en DEBUG
- [x] `loadChatAllowlist` : erreur si une entrée est invalide ou si la valeur ne donne aucune entrée ; avertissement si la liste n'est pas définie
- [x] `PurgeDisallowedChats` : messages + chats + dossiers médias hors liste, VACUUM + checkpoint WAL
- [x] README

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-16-confidentialite-journaux-stockage]]

## ✅ Definition of Done

- [x] Tâches terminées
- [x] Test [[TC-16-confidentialite-journaux-stockage]] **PASS**
- [x] Commit fait, message préfixé `[Sprint 16]`
- [x] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [x] Bip Windows émis

## 🧾 Journal

- `2026-10-31` — Sprint créé après le signalement de ks-support-agent-47 : 823 lignes de contenu de conversations privées hors liste blanche dans journalctl après l'appairage, base locale contenant tout l'historique.
- `2026-10-31` — Livré : commit 00b1298 (fix/16-confidentialite-journaux-stockage), mergé dans main. TC-16 PASS.
- `2026-10-31` — Diagnostic corrigé : liste blanche active et respectée, seul le texte des discussions autorisées était journalisé. Redéploiement en 119e9c3 par ks-support-agent-47, sans purge.

## 🔗 Liens

- ADR : [[ADR-0006-interface-pont-pipeline]]
