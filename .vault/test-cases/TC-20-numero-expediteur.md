---
type: test-case
id: TC-20-numero-expediteur
sprint: sprint-20
status: TODO          # TODO | PASS | FAIL
created: 2026-11-08
last_run:
tags: [test-case]
---

# TC-20 — Le numéro de l'expéditeur est stocké même quand WhatsApp n'envoie que le LID

> **sprint:: [[sprint-20-numero-expediteur]]**  ·  **status:: TODO**

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

- _(à remplir à l'exécution)_

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
