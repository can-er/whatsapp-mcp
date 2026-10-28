---
type: test-case
id: TC-06-store-permissions
sprint: sprint-06
status: TODO          # TODO | PASS | FAIL
created: 2026-10-29
last_run:
tags: [test-case]
---

# TC-06 — store/ en 0700, fichiers en 0600

> **sprint:: [[sprint-06-store-permissions]]**  ·  **status:: TODO**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté : tests purement unitaires.

## Étapes

1. `go test -run TestEnsurePrivate` (Linux, conteneur) : création → 0700 / 0600 ; dossier existant en 0755 → resserré à 0700 ; fichier existant en 0644 → 0600 ; contenu préservé.
2. Sous Windows le test est ignoré (modes POSIX non applicables).

## Résultat attendu

- Modes 0700 / 0600 constatés.

## Résultat obtenu

- _(à remplir à l'exécution)_

## Commande(s)

```bash
bash ".vault/tools/gocheck.sh" <worktree>/whatsapp-bridge   # go build + vet + test (Linux, CGO), golang:1.27-bookworm
```
