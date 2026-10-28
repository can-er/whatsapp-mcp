---
type: test-case
id: TC-03-download-path-traversal
sprint: sprint-03
status: TODO          # TODO | PASS | FAIL
created: 2026-10-29
last_run:
tags: [test-case]
---

# TC-03 — Aucun nom malveillant ne sort du dossier du chat

> **sprint:: [[sprint-03-download-path-traversal]]**  ·  **status:: TODO**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté : tests purement unitaires.

## Étapes

1. `go test -run 'TestSanitizeFilename|TestMediaFilePath'` : `../../x`, `..\..\x`, `/etc/passwd`, `C:\x`, `CON`, octet nul, nom vide, `..`, nom très long.
2. Pour chaque cas, le chemin final est un enfant direct de `store/<chat>`.

## Résultat attendu

- Tous les chemins restent dans le dossier du chat, aucun composant `..`.

## Résultat obtenu

- _(à remplir à l'exécution)_

## Commande(s)

```bash
bash ".vault/tools/gocheck.sh" <worktree>/whatsapp-bridge   # go build + vet + test (Linux, CGO), golang:1.27-bookworm
```
