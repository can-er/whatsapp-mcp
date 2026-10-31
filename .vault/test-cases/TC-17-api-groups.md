---
type: test-case
id: TC-17-api-groups
sprint: sprint-17
status: TODO          # TODO | PASS | FAIL
created: 2026-10-31
last_run:
tags: [test-case]
---

# TC-17 — La liste des groupes ne contient que les métadonnées autorisées

> **sprint:: [[sprint-17-api-groups]]**  ·  **status:: TODO**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté. Données de test fictives uniquement ([[ADR-0007-donnees-de-test-fictives]]).

## Étapes

1. `go test -run TestGroupsAPI` : 200 + JSON trié par nom, clés exactement `jid`, `name`, `participants_count` ; aucun numéro de participant dans la réponse.
2. Sans jeton → 401 ; POST → 405 ; avec le jeton d'envoi seul → 401 ; source en erreur → 502.

## Résultat attendu

- Tests PASS.

## Résultat obtenu

- _(à remplir à l'exécution)_

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
