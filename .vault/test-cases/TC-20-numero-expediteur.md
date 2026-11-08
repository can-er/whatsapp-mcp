---
type: test-case
id: TC-20-numero-expediteur
sprint: sprint-20
status: PASS          # TODO | PASS | FAIL
created: 2026-11-08
last_run: 2026-11-08
tags: [test-case]
---

# TC-20 — Le numéro de l'expéditeur est stocké même quand WhatsApp n'envoie que le LID

> **sprint:: [[sprint-20-numero-expediteur]]**  ·  **status:: PASS**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté. Données de test fictives uniquement ([[ADR-0007-donnees-de-test-fictives]]).

## Étapes

1. `go test -run TestSenderPhone` : expéditeur LID avec SenderAlt numéro → `sender_phone` = numéro ; LID seul résolu par la table LID → numéro ; LID inconnu → NULL ; message envoyé par soi → propre numéro.
2. Base créée avec l'ancien schéma (sans colonne) → ouverture : colonne ajoutée, lignes conservées ; rattrapage : `xxx@lid` et LID nu résolus remplis, numéro au format `@s.whatsapp.net` rempli, inconnus laissés NULL.
3. `sender` n'est pas modifié.

## Résultat attendu

- Tests PASS.

## Résultat obtenu

- `2026-11-08` — **PASS** (branche `fix/20-numero-expediteur`, commit 4e19227, empilée sur `fix/19`).
  - `TestSenderPhoneOfMessage` : LID + SenderAlt, LID résolu par la table LID, LID inconnu (vide), numéro, message envoyé par soi → PASS.
  - `TestSenderPhoneStoredFromLiveMessage` et `TestSenderPhoneStoredFromHistorySync` (SQLite, CGO) : `sender_phone` rempli, `sender` inchangé → PASS.
  - `TestSenderPhoneMigrationAndBackfill` : base à l'ancien schéma ouverte, colonne ajoutée, réouverture OK ; rattrapage de 5 lignes sur 7, les 2 identifiants inconnus restent NULL ; second passage : 0 → PASS.
  - Suite complète PASS. ⚠️ Exécutée sous Windows (CGO avec `zig cc`) car Docker Desktop était bloqué ; `gocheck.sh` (Linux) à relancer avant le merge.

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
