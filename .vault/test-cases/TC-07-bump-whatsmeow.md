---
type: test-case
id: TC-07-bump-whatsmeow
sprint: sprint-07
status: TODO          # TODO | PASS | FAIL
created: 2026-10-29
last_run:
tags: [test-case]
---

# TC-07 — Le pont compile et fonctionne avec whatsmeow à jour

> **sprint:: [[sprint-07-bump-whatsmeow]]**  ·  **status:: TODO**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté : tests purement unitaires.

## Étapes

1. `go list -m go.mau.fi/whatsmeow` plus récent que `v0.0.0-20250318233852`.
2. `gocheck.sh` : build + vet + tous les tests (dont `TestSQLStoreOpens`, CGO).

## Résultat attendu

- Version récente, `GOCHECK_OK`.

## Résultat obtenu

- _(à remplir à l'exécution)_

## Commande(s)

```bash
bash ".vault/tools/gocheck.sh" <worktree>/whatsapp-bridge   # go build + vet + test (Linux, CGO), golang:1.27-bookworm
```
