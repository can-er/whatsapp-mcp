---
type: test-case
id: TC-04-allowed-chats
sprint: sprint-04
status: PASS          # TODO | PASS | FAIL
created: 2026-10-29
last_run: 2026-10-29
tags: [test-case]
---

# TC-04 — Les conversations hors liste blanche sont ignorées et masquées

> **sprint:: [[sprint-04-allowed-chats]]**  ·  **status:: PASS**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté : tests purement unitaires.
- `uv` installé.

## Étapes

1. `go test -run TestChatAllowlist` : parsing, numéro → JID, groupes, liste vide = tout autorisé.
2. `uv run --with pytest pytest tests/test_allowed_chats.py` : base SQLite temporaire avec 3 chats, liste blanche de 1 → `list_chats`, `list_messages`, `search_contacts`, `get_chat`… ne renvoient que le chat autorisé.

## Résultat attendu

- Seul le chat autorisé est visible ; sans variable, les 3 chats sont visibles.

## Résultat obtenu

- `2026-10-29` — **PASS** (branche `fix/04-allowed-chats`, commit `39b8945`)
- Go : `TestChatAllowlistParsing` (numéro avec +/espaces, JID, groupe, appareil lié `:3`, entrée invalide ignorée, liste vide = tout autorisé) et `TestChatAllowlistAPI` (`/api/send` et `/api/download` → 403 hors liste, client nil jamais touché) PASS ; suite complète 10 tests, `GOCHECK_OK`.
- Python `tests/test_allowed_chats.py` : base SQLite temporaire (Alice autorisée, Bob et un groupe masqués) → les 9 fonctions de lecture ne voient qu'Alice, pour 4 notations de la liste ; sans variable, les 3 chats sont visibles ; le fichier .db n'est pas modifié. 19 passed (avec test_media_path).

## Commande(s)

```bash
bash ".vault/tools/gocheck.sh" <worktree>/whatsapp-bridge   # go build + vet + test (Linux, CGO), golang:1.27-bookworm
cd whatsapp-mcp-server && uv run --with pytest pytest -q tests/test_allowed_chats.py
```
