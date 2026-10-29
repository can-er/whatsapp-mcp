---
type: test-case
id: TC-12-code-appairage
sprint: sprint-12
status: TODO          # TODO | PASS | FAIL
created: 2026-10-29
last_run:
tags: [test-case]
---

# TC-12 — Les codes d'appairage sont livrés, limités en fréquence, puis confirmés

> **sprint:: [[sprint-12-code-appairage]]**  ·  **status:: TODO**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté. Données de test fictives uniquement ([[ADR-0007-donnees-de-test-fictives]]).

## Étapes

1. `go test -run TestPairing` : événement `pairing_code` conforme (id `pair-<ts>`, `expires_at`, `phone_hint` sans le numéro complet) écrit dans l'inbox et posté au webhook.
2. Deux codes à moins de 60 s → le second n'est pas demandé ; après 60 s → demandé.
3. `paired` émis après succès.
4. La boucle réseau (`PairPhone` réel) n'est pas testable sans compte : validée au déploiement.

## Résultat attendu

- Tests PASS.

## Résultat obtenu

- _(à remplir à l'exécution)_

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
