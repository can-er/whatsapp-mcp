---
type: test-case
id: TC-17-api-groups
sprint: sprint-17
status: PASS          # TODO | PASS | FAIL
created: 2026-10-31
last_run: 2026-10-31
tags: [test-case]
---

# TC-17 — La liste des groupes ne contient que les métadonnées autorisées

> **sprint:: [[sprint-17-api-groups]]**  ·  **status:: PASS**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté. Données de test fictives uniquement ([[ADR-0007-donnees-de-test-fictives]]).

## Étapes

1. `go test -run TestGroupsAPI` : 200 + JSON trié par nom, clés exactement `jid`, `name`, `participants_count` ; aucun numéro de participant dans la réponse.
2. Sans jeton → 401 ; POST → 405 ; avec le jeton d'envoi seul → 401 ; source en erreur → 502.

## Résultat attendu

- Tests PASS.

## Résultat obtenu

- `2026-10-31` — **PASS** : `bash .vault/tools/gocheck.sh .worktrees/17/whatsapp-bridge -v` (Docker golang:1.27-bookworm, CGO) → build + vet OK, 51 tests PASS dont TestGroupsAPI (tri par nom, 3 clés exactement, aucun JID de participant ; 401 sans jeton ou avec le jeton d'envoi, 405 en POST) et TestGroupsAPIErrors (503 non connecté / sans source, 502 erreur WhatsApp).
- Non testé : l'appel réel `GetJoinedGroups` (pas de compte en test), à vérifier sur la VM.

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
