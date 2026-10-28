---
type: test-case
id: TC-08-integration-rapport
sprint: sprint-08
status: TODO          # TODO | PASS | FAIL
created: 2026-10-29
last_run:
tags: [test-case]
---

# TC-08 — main du fork vert, PR ouvertes

> **sprint:: [[sprint-08-integration-rapport]]**  ·  **status:: TODO**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté : tests purement unitaires.
- `uv` installé.

## Étapes

1. `gocheck.sh` sur `main`.
2. Tous les tests Python sur `main`.
3. `gh pr list --repo lharries/whatsapp-mcp --author can-er` liste les 2 PR.

## Résultat attendu

- `GOCHECK_OK`, pytest PASS, 2 PR ouvertes.

## Résultat obtenu

- _(à remplir à l'exécution)_

## Commande(s)

```bash
bash ".vault/tools/gocheck.sh" <worktree>/whatsapp-bridge   # go build + vet + test (Linux, CGO), golang:1.27-bookworm
cd whatsapp-mcp-server && uv run --with pytest pytest -q tests
gh pr list --repo lharries/whatsapp-mcp --author can-er
```
