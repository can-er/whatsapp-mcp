---
type: test-case
id: TC-06-store-permissions
sprint: sprint-06
status: PASS          # TODO | PASS | FAIL
created: 2026-10-29
last_run: 2026-10-29
tags: [test-case]
---

# TC-06 — store/ en 0700, fichiers en 0600

> **sprint:: [[sprint-06-store-permissions]]**  ·  **status:: PASS**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté : tests purement unitaires.

## Étapes

1. `go test -run TestEnsurePrivate` (Linux, conteneur) : création → 0700 / 0600 ; dossier existant en 0755 → resserré à 0700 ; fichier existant en 0644 → 0600 ; contenu préservé.
2. Sous Windows le test est ignoré (modes POSIX non applicables).

## Résultat attendu

- Modes 0700 / 0600 constatés.

## Résultat obtenu

- `2026-10-29` — **PASS** (branche `fix/06-store-permissions` basée sur `upstream/main`, commit `ca8f262`)
- Linux (`gocheck.sh`, umask 022) : `TestEnsurePrivateDir`, `TestEnsurePrivateFile`, `TestEnsurePrivateMessageStore` PASS — `store/` 0755 préexistant → 0700 ; `messages.db`, `-wal`, `-shm` en 0600 (vrai SQLite, mode WAL) ; contenu préservé. `GOCHECK_OK`.
- Windows natif : les 3 tests sont `SKIP` (modes POSIX non applicables), `go vet` OK.
- Après merge dans `main` (01-06) : 13 tests Go + 30 tests Python PASS.
- PR upstream : https://github.com/lharries/whatsapp-mcp/pull/389

## Commande(s)

```bash
bash ".vault/tools/gocheck.sh" <worktree>/whatsapp-bridge   # go build + vet + test (Linux, CGO), golang:1.27-bookworm
```
