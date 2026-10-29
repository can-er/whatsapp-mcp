---
type: test-case
id: TC-13-api-history
sprint: sprint-13
status: TODO          # TODO | PASS | FAIL
created: 2026-10-29
last_run:
tags: [test-case]
---

# TC-13 — L'historique est paginé, filtré par la liste blanche et au format v1

> **sprint:: [[sprint-13-api-history]]**  ·  **status:: TODO**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté. Données de test fictives uniquement ([[ADR-0007-donnees-de-test-fictives]]).

## Étapes

1. `go test -run TestHistory` (SQLite réel, CGO) : 120 messages dont 30 à la même seconde → pagination par 50 sans trou ni doublon.
2. Chat hors liste blanche → 403 ; jeton absent → 401 ; `limit` invalide → 400.
3. Média téléchargé → `status: ok` + chemin + sha256 ; sinon `not_downloaded`, `path: null`.

## Résultat attendu

- Tests PASS.

## Résultat obtenu

- _(à remplir à l'exécution)_

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
