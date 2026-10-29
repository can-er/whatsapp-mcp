---
type: test-case
id: TC-05-read-only-mode
sprint: sprint-05
status: PASS          # TODO | PASS | FAIL
created: 2026-10-29
last_run: 2026-10-29
tags: [test-case]
---

# TC-05 — Les outils d'envoi disparaissent en lecture seule

> **sprint:: [[sprint-05-read-only-mode]]**  ·  **status:: PASS**

## Préconditions

- `uv` installé.

## Étapes

1. `uv run --with pytest pytest tests/test_read_only.py`.
2. Sans variable : les 12 outils sont listés.
3. `WHATSAPP_MCP_READ_ONLY=1` (et `true`) : les 3 outils d'envoi sont absents, les 9 autres présents.

## Résultat attendu

- Tests PASS.

## Résultat obtenu

- `2026-10-29` — **PASS** (branche `fix/05-read-only-mode` basée sur `upstream/main`, commit `fdef102`)
- `uv run --with pytest pytest -q tests/test_read_only.py` → `11 passed` : 12 outils par défaut (non défini, `""`, `0`, `false`, `no`) ; 9 outils de lecture seulement pour `1`, `true`, `TRUE`, `yes`, `on`, `" 1 "`.
- Après merge dans `main` (avec 01-04) : 30 tests Python PASS.
- PR upstream : https://github.com/lharries/whatsapp-mcp/pull/388

## Commande(s)

```bash
cd whatsapp-mcp-server && uv run --with pytest pytest -q tests/test_read_only.py
```
