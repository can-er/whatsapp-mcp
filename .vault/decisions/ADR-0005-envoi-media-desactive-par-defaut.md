---
type: adr
id: ADR-0005
title: "Envoi de médias désactivé tant que WHATSAPP_MEDIA_DIR n'est pas défini"
status: accepted
created: 2026-10-29
tags: [adr]
---

# ADR-0005 — Envoi de médias désactivé tant que WHATSAPP_MEDIA_DIR n'est pas défini

## Contexte

Restreindre `media_path` exige un dossier autorisé connu à la fois du pont et du serveur MCP, qui n'ont pas le même répertoire courant.

## Décision

Le pont et le serveur MCP lisent `WHATSAPP_MEDIA_DIR`. Si elle n'est pas définie, l'envoi de fichiers est refusé avec un message explicite (l'envoi de texte reste possible). La conversion audio ffmpeg écrit dans ce dossier.

## Conséquences

- Sûr par défaut : aucune lecture de fichier arbitraire possible.
- Changement de comportement documenté dans le README.

## Alternatives écartées

- Dossier par défaut implicite — divergences de chemin entre les deux processus, surprises.

## 🔗 Liens

- Sprints concernés : [[sprint-02-media-path-allowlist]]
