---
type: test-case
id: TC-13-api-history
sprint: sprint-13
status: PASS          # TODO | PASS | FAIL
created: 2026-10-29
last_run: 2026-10-29
tags: [test-case]
---

# TC-13 — L'historique est paginé, filtré par la liste blanche et au format v1

> **sprint:: [[sprint-13-api-history]]**  ·  **status:: PASS**

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

- `2026-10-29` — **PASS** (branche `fix/13-api-history`, commit `fe77f23`)
- `TestHistoryPagination` (SQLite réel) : 122 messages, dont 30 à la même seconde et la moitié stockée en UTC+1 → 3 pages de 50, ni trou ni doublon, ordre (timestamp, id) décroissant ; `received_at` null ; média téléchargé → `ok` + chemin + sha256, sinon `not_downloaded` ; curseur RFC 3339 accepté ; chat hors liste jamais renvoyé.
- `TestHistoryRejections` : hors liste → 403, sans jeton → 401, chat absent / limit 0, 501, abc / curseur invalide → 400.
- Suite complète `GOCHECK_OK` ; `go vet` natif Windows OK.

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
