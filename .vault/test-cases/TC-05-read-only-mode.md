---
type: test-case
id: TC-05-read-only-mode
sprint: sprint-05
status: TODO          # TODO | PASS | FAIL
created: 2026-10-29
last_run:
tags: [test-case]
---

# TC-05 — Les outils d'envoi disparaissent en lecture seule

> **sprint:: [[sprint-05-read-only-mode]]**  ·  **status:: TODO**

## Préconditions

- `uv` installé.

## Étapes

1. `uv run --with pytest pytest tests/test_read_only.py`.
2. Sans variable : les 12 outils sont listés.
3. `WHATSAPP_MCP_READ_ONLY=1` (et `true`) : les 3 outils d'envoi sont absents, les 9 autres présents.

## Résultat attendu

- Tests PASS.

## Résultat obtenu

- _(à remplir à l'exécution)_

## Commande(s)

```bash
cd whatsapp-mcp-server && uv run --with pytest pytest -q tests/test_read_only.py
```
