---
type: test-case
id: TC-12-code-appairage
sprint: sprint-12
status: PASS          # TODO | PASS | FAIL
created: 2026-10-29
last_run: 2026-10-29
tags: [test-case]
---

# TC-12 — Les codes d'appairage sont livrés, limités en fréquence, puis confirmés

> **sprint:: [[sprint-12-code-appairage]]**  ·  **status:: PASS**

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

- `2026-10-29` — **PASS** (branche `fix/12-code-appairage`, commit `87cbb78`)
- `TestPairingCodeDeliveredToInboxAndWebhook` : événement `pairing_code` conforme (`pair-20261029T093002Z`, `expires_at` = +160 s, `phone_hint` = `…01`, numéro complet absent), même JSON reçu par le webhook ; `paired` écrit ensuite.
- `TestPairingCodeMinInterval` : pas de nouveau code à 59 s, autorisé à 60 s. `TestPairingPhoneHint`, `TestPairTimeout` (30 min par défaut) PASS.
- Suite complète `GOCHECK_OK`. Limite : la boucle `PairPhone` réelle n'est testable qu'avec un compte (déploiement).

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
