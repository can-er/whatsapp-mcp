---
type: test-case
id: TC-04-allowed-chats
sprint: sprint-04
status: TODO          # TODO | PASS | FAIL
created: 2026-10-29
last_run:
tags: [test-case]
---

# TC-04 — Les conversations hors liste blanche sont ignorées et masquées

> **sprint:: [[sprint-04-allowed-chats]]**  ·  **status:: TODO**

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

- _(à remplir à l'exécution)_

## Commande(s)

```bash
bash ".vault/tools/gocheck.sh" <worktree>/whatsapp-bridge   # go build + vet + test (Linux, CGO), golang:1.27-bookworm
cd whatsapp-mcp-server && uv run --with pytest pytest -q tests/test_allowed_chats.py
```
