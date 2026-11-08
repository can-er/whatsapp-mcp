---
type: test-case
id: TC-19-legende-medias
sprint: sprint-19
status: TODO          # TODO | PASS | FAIL
created: 2026-11-08
last_run:
tags: [test-case]
---

# TC-19 — La légende d'un média est stockée comme contenu du message

> **sprint:: [[sprint-19-legende-medias]]**  ·  **status:: TODO**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté. Données de test fictives uniquement ([[ADR-0007-donnees-de-test-fictives]]).

## Étapes

1. `go test -run TestCaption` : image, vidéo et document avec légende → `content` = légende ; média sans légende → `content` vide, `media_type` renseigné ; texte simple inchangé.
2. Message d'historique (HistorySyncMsg) image avec légende → ligne stockée avec la légende (SQLite, CGO).
3. L'extraction du stockage et celle des événements donnent le même texte pour chaque type.

## Résultat attendu

- Tests PASS.

## Résultat obtenu

- _(à remplir à l'exécution)_

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
