---
type: test-case
id: TC-15-appairage-robuste
sprint: sprint-15
status: PASS          # TODO | PASS | FAIL
created: 2026-10-29
last_run: 2026-10-29
tags: [test-case]
---

# TC-15 — L'appairage livre toujours ses événements et s'arrête proprement

> **sprint:: [[sprint-15-appairage-robuste]]**  ·  **status:: PASS**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté. Données de test fictives uniquement ([[ADR-0007-donnees-de-test-fictives]]).

## Étapes

1. `go test -run TestDeliver` : `events/` supprimé → recréé et événement écrit ; inbox inutilisable → webhook quand même appelé.
2. `go test -run TestPairing` : compteur de demandes, limite atteinte après 5, recul 5/10/20… min plafonné à 1 h, remis à zéro après un code, détection `ErrIQRateOverLimit`.
3. Événement `pairing_stopped` (raison `max_codes` / `timeout` / `rate_limited`) écrit + posté ; erreur d'abandon reconnue (`errors.Is`) → code de sortie 3.
4. La boucle réseau réelle reste à valider au redéploiement (pas de compte WhatsApp en test).

## Résultat attendu

- Tests PASS.

## Résultat obtenu

- `2026-10-29` — **PASS** : `bash .vault/tools/gocheck.sh .worktrees/15/whatsapp-bridge -v` (Docker golang:1.27-bookworm, CGO) → build + vet OK, 44 tests PASS dont TestDeliverRecreatesMissingEventsDir, TestDeliverCallsWebhookWhenInboxFails, TestPairingRequestLimit, TestPairingRateLimitBackoff (5/10/20/40/60/60 min), TestPairingStoppedEvent, TestPairMaxCodes.
- Non testé : la boucle réseau réelle (`PairPhone`, fermeture de session après 429) — à valider au redéploiement sur la VM ks-agent.
- Diagnostic confirmé par le journal de la VM (ks-support-agent-47) : la livraison n'était pas en cause (7 codes écrits entre 17:49 et 18:03, aucun « Failed to deliver », inbox et webhook actifs à chaque démarrage). Cause réelle : 429 rate-overlimit dès 18:06, puis 173 refus, l'ancien code redemandant un code à chaque rotation de QR (environ toutes les 20 s). La limite de demandes et le recul sur 429 de ce sprint traitent ce cas ; le correctif de `deliver` reste valable mais n'était pas la cause.

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
