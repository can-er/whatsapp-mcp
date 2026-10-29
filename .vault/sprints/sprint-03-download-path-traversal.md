---
type: sprint
id: sprint-03
title: "Noms de fichiers médias assainis (path traversal)"
status: DONE          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed: 2026-10-29
tags: [sprint]
---

# Sprint 03 — Noms de fichiers médias assainis (path traversal)

> **status:: DONE**  ·  **test:: [[TC-03-download-path-traversal]]**

## 🎯 Objectif

Un nom de document fourni par l'expéditeur (ou un JID) ne peut plus faire écrire un fichier hors de `store/<chat>/`.

## 📋 Tâches

- [x] `sanitizeFilename` : `/` et `\` → dernier composant, `..`, caractères de contrôle et réservés Windows, noms réservés (CON, NUL, COM1…), longueur, repli sur un nom généré
- [x] Assainir à l'extraction (stockage) **et** au téléchargement (anciennes lignes en base)
- [x] Assainir le composant JID du dossier
- [x] Vérification finale : chemin résolu dans `chatDir`

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-03-download-path-traversal]]

## ✅ Definition of Done

- [x] Tâches terminées
- [x] Test [[TC-03-download-path-traversal]] **PASS**
- [x] Commit fait, message préfixé `[Sprint 03]`
- [x] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [x] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé.
- `2026-10-29` — Noms de médias et JID assainis (extraction + téléchargement) ; TC-03 PASS ; mergé dans main.

## 🔗 Liens

- ADR : —
