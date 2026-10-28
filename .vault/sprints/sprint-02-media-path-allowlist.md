---
type: sprint
id: sprint-02
title: "Envoi de médias restreint au dossier WHATSAPP_MEDIA_DIR"
status: DONE          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed: 2026-10-29
tags: [sprint]
---

# Sprint 02 — Envoi de médias restreint au dossier WHATSAPP_MEDIA_DIR

> **status:: DONE**  ·  **test:: [[TC-02-media-path-allowlist]]**

## 🎯 Objectif

`media_path` ne peut plus désigner qu'un fichier régulier situé sous `WHATSAPP_MEDIA_DIR` (chemins résolus, liens symboliques suivis) ; envoi de médias désactivé si la variable n'est pas définie. Même garde côté Python.

## 📋 Tâches

- [x] Go : `resolveMediaPath` (Abs + EvalSymlinks + Rel), refus sinon
- [x] Go : titre de document via `filepath.Base`
- [x] Python : même garde dans `send_file` / `send_audio_message`, conversion ffmpeg écrite dans le dossier d'envoi
- [x] README

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-02-media-path-allowlist]]

## ✅ Definition of Done

- [x] Tâches terminées
- [x] Test [[TC-02-media-path-allowlist]] **PASS**
- [x] Commit fait, message préfixé `[Sprint 02]`
- [x] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [x] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé.
- `2026-10-29` — Garde media_path Go + Python, conversion audio dans le dossier d'envoi ; TC-02 PASS ; mergé dans main.

## 🔗 Liens

- ADR : [[ADR-0005-envoi-media-desactive-par-defaut]]
