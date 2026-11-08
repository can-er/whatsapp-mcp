---
type: test-case
id: TC-19-legende-medias
sprint: sprint-19
status: PASS          # TODO | PASS | FAIL
created: 2026-11-08
last_run: 2026-11-08
tags: [test-case]
---

# TC-19 — La légende d'un média est stockée comme contenu du message

> **sprint:: [[sprint-19-legende-medias]]**  ·  **status:: PASS**

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

- `2026-11-08` — **PASS** (branche `fix/19-legende-medias`, commit cb12a33).
  - `TestCaptionStoredAsContent` : texte, texte étendu, légendes image / vidéo / document, média sans légende, audio, nil → PASS ; même texte que `messageText` des événements.
  - `TestCaptionStoredFromLiveMessage` et `TestCaptionStoredFromHistorySync` (SQLite, CGO) → PASS.
  - Suite complète PASS (3 tests de permissions Unix ignorés sous Windows).
  - Première exécution sous Windows (Go 1.27, CGO avec `zig cc`), Docker Desktop étant bloqué ce soir-là.
- `2026-11-08` — **PASS sous Linux** : `gocheck.sh` (build + vet + test, `golang:1.27-bookworm`, CGO) sur cb12a33 → `GOCHECK_OK`, 57 tests et sous-tests PASS, aucun ignoré.

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
