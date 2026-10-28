---
type: sprint
id: sprint-03
title: "Noms de fichiers médias assainis (path traversal)"
status: TODO          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed:
tags: [sprint]
---

# Sprint 03 — Noms de fichiers médias assainis (path traversal)

> **status:: TODO**  ·  **test:: [[TC-03-download-path-traversal]]**

## 🎯 Objectif

Un nom de document fourni par l'expéditeur (ou un JID) ne peut plus faire écrire un fichier hors de `store/<chat>/`.

## 📋 Tâches

- [ ] `sanitizeFilename` : `/` et `\` → dernier composant, `..`, caractères de contrôle et réservés Windows, noms réservés (CON, NUL, COM1…), longueur, repli sur un nom généré
- [ ] Assainir à l'extraction (stockage) **et** au téléchargement (anciennes lignes en base)
- [ ] Assainir le composant JID du dossier
- [ ] Vérification finale : chemin résolu dans `chatDir`

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-03-download-path-traversal]]

## ✅ Definition of Done

- [ ] Tâches terminées
- [ ] Test [[TC-03-download-path-traversal]] **PASS**
- [ ] Commit fait, message préfixé `[Sprint 03]`
- [ ] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [ ] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé.

## 🔗 Liens

- ADR : —
