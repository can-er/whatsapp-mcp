---
type: test-case
id: TC-10-evenements-edit-revoke-reaction
sprint: sprint-10
status: TODO          # TODO | PASS | FAIL
created: 2026-10-29
last_run:
tags: [test-case]
---

# TC-10 — Edit / revoke / reaction produisent le bon type d'événement

> **sprint:: [[sprint-10-evenements-edit-revoke-reaction]]**  ·  **status:: TODO**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté : tests unitaires (événements construits à la main).

## Étapes

1. `go test -run TestEventKinds` : protos construits à la main pour edit, revoke, reaction → type et target corrects.

## Résultat attendu

- Tests PASS.

## Résultat obtenu

- _(à remplir à l'exécution)_

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
