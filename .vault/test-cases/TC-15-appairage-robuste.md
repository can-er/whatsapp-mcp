---
type: test-case
id: TC-15-appairage-robuste
sprint: sprint-15
status: TODO          # TODO | PASS | FAIL
created: 2026-10-29
last_run:
tags: [test-case]
---

# TC-15 — L'appairage livre toujours ses événements et s'arrête proprement

> **sprint:: [[sprint-15-appairage-robuste]]**  ·  **status:: TODO**

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

- _(à remplir à l'exécution)_

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
