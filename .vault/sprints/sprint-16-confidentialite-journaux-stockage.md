---
type: sprint
id: sprint-16
title: "Confidentialité : aucun texte dans les journaux, liste blanche stricte, purge hors liste (fork + correctif)"
status: ONGOING          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-31
closed:
tags: [sprint]
---

# Sprint 16 — Confidentialité : aucun texte dans les journaux, liste blanche stricte, purge hors liste (fork + correctif)

> **status:: ONGOING**  ·  **test:: [[TC-16-confidentialite-journaux-stockage]]**

## 🎯 Objectif

Le pont ne journalise plus jamais le texte d'un message (ni noms de contacts en INFO). Une valeur de `WHATSAPP_ALLOWED_CHATS` non vide mais invalide empêche le démarrage au lieu de tout autoriser. `WHATSAPP_PURGE_DISALLOWED_CHATS=true` supprime au démarrage les discussions, messages et médias stockés hors liste blanche, puis compacte la base (VACUUM).

## 📋 Tâches

- [ ] Retirer le texte des messages de tous les journaux (direct + synchro d'historique) ; métadonnées en DEBUG seulement
- [ ] Noms de contacts / groupes : journalisés en DEBUG
- [ ] `loadChatAllowlist` : erreur si une entrée est invalide ou si la valeur ne donne aucune entrée ; avertissement si la liste n'est pas définie
- [ ] `PurgeDisallowedChats` : messages + chats + dossiers médias hors liste, VACUUM + checkpoint WAL
- [ ] README

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-16-confidentialite-journaux-stockage]]

## ✅ Definition of Done

- [ ] Tâches terminées
- [ ] Test [[TC-16-confidentialite-journaux-stockage]] **PASS**
- [ ] Commit fait, message préfixé `[Sprint 16]`
- [ ] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [ ] Bip Windows émis

## 🧾 Journal

- `2026-10-31` — Sprint créé après le signalement de ks-support-agent-47 : 823 lignes de contenu de conversations privées hors liste blanche dans journalctl après l'appairage, base locale contenant tout l'historique.

## 🔗 Liens

- ADR : [[ADR-0006-interface-pont-pipeline]]
