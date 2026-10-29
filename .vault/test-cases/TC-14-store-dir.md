---
type: test-case
id: TC-14-store-dir
sprint: sprint-14
status: PASS          # TODO | PASS | FAIL
created: 2026-10-29
last_run: 2026-10-29
tags: [test-case]
---

# TC-14 — Toutes les données vont dans WHATSAPP_STORE_DIR

> **sprint:: [[sprint-14-store-dir]]**  ·  **status:: PASS**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm` ; `uv` installé.

## Étapes

1. `go test -run TestStoreDir` : avec `WHATSAPP_STORE_DIR` vers un dossier temporaire, `NewMessageStore`, jeton, audit et médias y sont créés (0700/0600) et rien dans `./store`.
2. Sans variable : chemin `store/` relatif inchangé.
3. Python : `tests/test_store_dir.py` (chemins dérivés de la variable).

## Résultat attendu

- Tests PASS.

## Résultat obtenu

- `2026-10-29` — **PASS** (branche `fix/14-store-dir`, commit `ef71428`)
- `TestStoreDirFromEnv` : avec `WHATSAPP_STORE_DIR` (chemin contenant un espace), messages.db, whatsapp.db (sqlstore réel), bridge_token et send_audit.jsonl créés dedans en 0600, dossier en 0700, médias dessous, rien dans `./store`. `TestStoreDirDefault` : défaut `store/` inchangé.
- Python `tests/test_store_dir.py` : chemins et lecture du jeton dérivés de la variable ; défaut inchangé. 32 passed.
- Suite Go complète `GOCHECK_OK` (38 tests).

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
cd whatsapp-mcp-server && uv run --with pytest pytest -q tests
```
