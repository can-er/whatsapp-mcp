---
type: test-case
id: TC-01-api-rest-auth
sprint: sprint-01
status: TODO          # TODO | PASS | FAIL
created: 2026-10-29
last_run:
tags: [test-case]
---

# TC-01 — Le pont rejette les requêtes non authentifiées, cross-origin ou non JSON

> **sprint:: [[sprint-01-api-rest-auth]]**  ·  **status:: TODO**

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

- _(à remplir à l'exécution)_

## Commande(s)

```bash
bash ".vault/tools/gocheck.sh" <worktree>/whatsapp-bridge   # go build + vet + test (Linux, CGO), golang:1.27-bookworm
```
