---
type: test-case
id: TC-14-store-dir
sprint: sprint-14
status: TODO          # TODO | PASS | FAIL
created: 2026-10-29
last_run:
tags: [test-case]
---

# TC-14 — Toutes les données vont dans WHATSAPP_STORE_DIR

> **sprint:: [[sprint-14-store-dir]]**  ·  **status:: TODO**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm` ; `uv` installé.

## Étapes

1. `go test -run TestStoreDir` : avec `WHATSAPP_STORE_DIR` vers un dossier temporaire, `NewMessageStore`, jeton, audit et médias y sont créés (0700/0600) et rien dans `./store`.
2. Sans variable : chemin `store/` relatif inchangé.
3. Python : `tests/test_store_dir.py` (chemins dérivés de la variable).

## Résultat attendu

- Tests PASS.

## Résultat obtenu

- _(à remplir à l'exécution)_

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
cd whatsapp-mcp-server && uv run --with pytest pytest -q tests
```
