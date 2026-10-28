---
type: test-case
id: TC-01-api-rest-auth
sprint: sprint-01
status: PASS          # TODO | PASS | FAIL
created: 2026-10-29
last_run: 2026-10-29
tags: [test-case]
---

# TC-01 — Le pont rejette les requêtes non authentifiées, cross-origin ou non JSON

> **sprint:: [[sprint-01-api-rest-auth]]**  ·  **status:: PASS**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté : tests purement unitaires.

## Étapes

1. `go test -run 'TestAPI|TestBridgeToken|TestListenAddr'` : requêtes httptest sur le mux réel (client nil : un rejet ne touche jamais WhatsApp).
2. Cas : sans jeton, mauvais jeton, `text/plain` (CSRF « simple request »), en-tête `Origin`, corps > 1 MiB, GET, requête valide.
3. Jeton : création, persistance, réutilisation, permissions 0600, priorité de la variable d'environnement.

## Résultat attendu

- 401 / 403 / 415 / 413 / 405 sur les cas rejetés ; le cas valide atteint le handler.
- Adresse par défaut `127.0.0.1:8080`.

## Résultat obtenu

- `2026-10-29` — **PASS** (branche `fix/01-api-rest-auth`, commit `ae28e73`)
- `gocheck.sh -v -run 'TestAPI|TestBridgeToken|TestListenAddr'` : 5 tests, 20 sous-cas du rejet (10 cas × `/api/send` et `/api/download`) tous PASS ; `GOCHECK_OK`.
  - sans jeton / mauvais jeton / sans schéma `Bearer` → 401 ; `text/plain` et Content-Type absent → 415 ; `Origin` → 403 (même sans jeton) ; corps > 1 MiB → 413 ; JSON invalide → 400 ; GET → 405 ; requête valide → handler atteint.
  - jeton : 64 hex, fichier 0600, réutilisé au 2e appel, `WHATSAPP_BRIDGE_TOKEN` prioritaire (aucun fichier créé).
- Contrôles annexes : `CGO_ENABLED=0 go vet` natif Windows OK ; `py_compile whatsapp.py` OK.

## Commande(s)

```bash
bash ".vault/tools/gocheck.sh" <worktree>/whatsapp-bridge   # go build + vet + test (Linux, CGO), golang:1.27-bookworm
```
