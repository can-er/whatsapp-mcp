---
type: test-case
id: TC-02-media-path-allowlist
sprint: sprint-02
status: PASS          # TODO | PASS | FAIL
created: 2026-10-29
last_run: 2026-10-29
tags: [test-case]
---

# TC-02 — Seuls les fichiers du dossier d'envoi peuvent être lus

> **sprint:: [[sprint-02-media-path-allowlist]]**  ·  **status:: PASS**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté : tests purement unitaires.
- `uv` installé (test Python).

## Étapes

1. `go test -run TestResolveMediaPath` : fichier autorisé, `../`, chemin absolu hors dossier, lien symbolique sortant, dossier, variable non définie.
2. `uv run --with pytest pytest tests/test_media_path.py` : même matrice côté Python.

## Résultat attendu

- Seul le fichier sous le dossier d'envoi est accepté ; tous les autres cas renvoient une erreur.

## Résultat obtenu

- `2026-10-29` — **PASS** (branche `fix/02-media-path-allowlist`, commit `f9c2225`)
- Go `TestResolveMediaPath` : 12 sous-cas PASS — acceptés : absolu dedans, imbriqué, relatif au dossier ; refusés : absolu hors dossier, `../`, traversée depuis l'intérieur, le dossier lui-même, un sous-dossier, fichier absent, voisin à préfixe commun (`outbox-evil`), lien symbolique sortant, variable non définie. `GOCHECK_OK`.
- Python `tests/test_media_path.py` : 12 passed sous Windows (uv) et sous Linux (`python:3.12-slim`, cas symlink réellement exécuté) ; `send_file` / `send_audio_message` n'appellent jamais le pont pour un fichier hors dossier.

## Commande(s)

```bash
bash ".vault/tools/gocheck.sh" <worktree>/whatsapp-bridge   # go build + vet + test (Linux, CGO), golang:1.27-bookworm
cd whatsapp-mcp-server && uv run --with pytest pytest -q tests/test_media_path.py
```
