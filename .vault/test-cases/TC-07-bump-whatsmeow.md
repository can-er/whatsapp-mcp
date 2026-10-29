---
type: test-case
id: TC-07-bump-whatsmeow
sprint: sprint-07
status: PASS          # TODO | PASS | FAIL
created: 2026-10-29
last_run: 2026-10-29
tags: [test-case]
---

# TC-07 — Le pont compile et fonctionne avec whatsmeow à jour

> **sprint:: [[sprint-07-bump-whatsmeow]]**  ·  **status:: PASS**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté : tests purement unitaires.

## Étapes

1. `go list -m go.mau.fi/whatsmeow` plus récent que `v0.0.0-20250318233852`.
2. `gocheck.sh` : build + vet + tous les tests (dont `TestSQLStoreOpens`, CGO).

## Résultat attendu

- Version récente, `GOCHECK_OK`.

## Résultat obtenu

- `2026-10-29` — **PASS** (branche `fix/07-bump-whatsmeow` = fix/04 + fix/05 + fix/06 + bump, commit `8d8c48c`)
- `go list -m go.mau.fi/whatsmeow` → `v0.0.0-20261028111316-d24d079a3755` (au lieu de `v0.0.0-20250318233852-06705625cf82`) ; `go 1.26.0` dans go.mod.
- `gocheck.sh` : 15 tests PASS dont `TestSQLStoreOpens` (sqlstore réel + client non appairé, aucune connexion) et `TestChatAllowlistResolvesLID` (LID → numéro via le LID store) ; `GOCHECK_OK`.
- Windows natif : `CGO_ENABLED=0 go vet` + `go test` OK (tests CGO exclus par build tag).
- Note : branche partie de fix/04 avec merge de 05 et 06 (au lieu de fix/04 seul, cf. ADR-0002) pour adapter l'API sur le code complet ; merge dans `main` sans conflit, code identique.

## Commande(s)

```bash
bash ".vault/tools/gocheck.sh" <worktree>/whatsapp-bridge   # go build + vet + test (Linux, CGO), golang:1.27-bookworm
```
