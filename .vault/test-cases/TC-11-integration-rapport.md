---
type: test-case
id: TC-11-integration-rapport
sprint: sprint-11
status: PASS          # TODO | PASS | FAIL
created: 2026-10-29
last_run: 2026-10-29
tags: [test-case]
---

# TC-11 — main du fork vert, PR ouvertes

> **sprint:: [[sprint-11-integration-rapport]]**  ·  **status:: PASS**

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

- `2026-10-29` — **PASS** (`main` = `origin/main` = `7ae8eb3` avant ce commit)
- Go, Linux (CGO) : 32 tests PASS, `GOCHECK_OK`. Go, Windows natif : `go vet` et `go test` OK.
- Python : 30 passed.
- `check_no_pii.py --history` : 0 occurrence ; diffs des PR #388 et #389 : 0 occurrence.
- PR upstream ouvertes : #388 (lecture seule) https://github.com/lharries/whatsapp-mcp/pull/388 et #389 (permissions) https://github.com/lharries/whatsapp-mcp/pull/389.
- Fork : https://github.com/can-er/whatsapp-mcp, `main` + branches `fix/01` à `fix/10`.

## Commande(s)

```bash
bash ".vault/tools/gocheck.sh" <worktree>/whatsapp-bridge   # go build + vet + test (Linux, CGO), golang:1.27-bookworm
cd whatsapp-mcp-server && uv run --with pytest pytest -q tests
gh pr list --repo lharries/whatsapp-mcp --author can-er
```
