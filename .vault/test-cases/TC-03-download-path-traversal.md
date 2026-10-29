---
type: test-case
id: TC-03-download-path-traversal
sprint: sprint-03
status: PASS          # TODO | PASS | FAIL
created: 2026-10-29
last_run: 2026-10-29
tags: [test-case]
---

# TC-03 — Aucun nom malveillant ne sort du dossier du chat

> **sprint:: [[sprint-03-download-path-traversal]]**  ·  **status:: PASS**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté : tests purement unitaires.

## Étapes

1. `go test -run 'TestSanitizeFilename|TestMediaFilePath'` : `../../x`, `..\..\x`, `/etc/passwd`, `C:\x`, `CON`, octet nul, nom vide, `..`, nom très long.
2. Pour chaque cas, le chemin final est un enfant direct de `store/<chat>`.

## Résultat attendu

- Tous les chemins restent dans le dossier du chat, aucun composant `..`.

## Résultat obtenu

- `2026-10-29` — **PASS** (branche `fix/03-download-path-traversal`, commit `4f5227f`)
- `TestSanitizeFilename` : 21 cas + troncature UTF-8 (`../../evil.sh`→`evil.sh`, `..\..\evil.bat`→`evil.bat`, `C:\Windows\win.ini`→`win.ini`, `..`→vide, `CON`→`_CON`, octet nul et U+202E supprimés…).
- `TestMediaFilePathStaysInChatDir` : 4 JID (dont `../../outside`) × 8 noms → chaque chemin est un enfant direct de `store/<chat>` ; noms vides / `..` refusés.
- Exécuté sous Linux (`gocheck.sh`, `GOCHECK_OK`) **et** en natif Windows (`CGO_ENABLED=0 go test`).

## Commande(s)

```bash
bash ".vault/tools/gocheck.sh" <worktree>/whatsapp-bridge   # go build + vet + test (Linux, CGO), golang:1.27-bookworm
```
