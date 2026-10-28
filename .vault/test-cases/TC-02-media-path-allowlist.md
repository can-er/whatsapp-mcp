---
type: test-case
id: TC-02-media-path-allowlist
sprint: sprint-02
status: TODO          # TODO | PASS | FAIL
created: 2026-10-29
last_run:
tags: [test-case]
---

# TC-02 — Seuls les fichiers du dossier d'envoi peuvent être lus

> **sprint:: [[sprint-02-media-path-allowlist]]**  ·  **status:: TODO**

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

- _(à remplir à l'exécution)_

## Commande(s)

```bash
bash ".vault/tools/gocheck.sh" <worktree>/whatsapp-bridge   # go build + vet + test (Linux, CGO), golang:1.27-bookworm
cd whatsapp-mcp-server && uv run --with pytest pytest -q tests/test_media_path.py
```
