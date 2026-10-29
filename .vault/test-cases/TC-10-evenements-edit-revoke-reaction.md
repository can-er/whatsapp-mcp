---
type: test-case
id: TC-10-evenements-edit-revoke-reaction
sprint: sprint-10
status: PASS          # TODO | PASS | FAIL
created: 2026-10-29
last_run: 2026-10-29
tags: [test-case]
---

# TC-10 — Edit / revoke / reaction produisent le bon type d'événement

> **sprint:: [[sprint-10-evenements-edit-revoke-reaction]]**  ·  **status:: PASS**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté : tests unitaires (événements construits à la main).

## Étapes

1. `go test -run TestEventKinds` : protos construits à la main pour edit, revoke, reaction → type et target corrects.

## Résultat attendu

- Tests PASS.

## Résultat obtenu

- `2026-10-29` — **PASS** (branche `fix/10-evenements-edit-revoke-reaction`, commit `8ca6d8f`)
- `TestEventKinds` : edit (nouveau texte), revoke, reaction 👍, reaction retirée (`""`) → bon `type` et `target_message_id`, aucun média ; message normal → `type: message` sans champs de mise à jour.
- `TestEventUpdateWrittenToInbox` : réaction écrite dans `events/` avec `target_message_id` et `reaction`.
- `TestUpdateMessageContent` (SQLite réel) : le texte stocké est remplacé après un edit.
- Suite complète `GOCHECK_OK`.

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
