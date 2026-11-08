---
type: test-case
id: TC-21-telechargement-media
sprint: sprint-21
status: PASS          # TODO | PASS | FAIL
created: 2026-11-08
last_run: 2026-11-08
tags: [test-case]
---

# TC-21 — Le téléchargement d'un média réussit avec le chemin signé, la copie de l'inbox ou un renouvellement

> **sprint:: [[sprint-21-telechargement-media]]**  ·  **status:: PASS**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté. Données de test fictives uniquement ([[ADR-0007-donnees-de-test-fictives]]).

## Étapes

1. `go test -run TestDirectPath` : le chemin tiré d'une URL garde ses paramètres de requête (`oh`, `oe`) ; `direct_path` stocké est utilisé en priorité.
2. `go test -run TestMediaDownload` : fichier présent dans l'inbox avec la bonne empreinte → renvoyé sans appel réseau ; mauvaise empreinte → ignoré ; 403 → demande de renouvellement puis succès avec le nouveau chemin, stocké en base ; renouvellement en échec ou délai dépassé → erreur claire ; autre erreur → pas de renouvellement.
3. Le média téléchargé est écrit avec les permissions privées (0600).

## Résultat attendu

- Tests PASS.

## Résultat obtenu

- `2026-11-08` — **PASS** (branche `fix/21-telechargement-media`, commit ed16824, empilée sur `fix/20`).
  - `TestDirectPathFromURL` : la requête signée (`oh`, `oe`) est conservée ; l'ancien code la coupait → PASS.
  - `TestMediaRetryWaiter` → PASS.
  - `TestMediaDownloadUsesStoredDirectPath`, `…FallsBackToURLWithQuery`, `…FromInbox`, `…RetriesExpiredMedia` (403, 404, 410), `…RetryFailures`, `TestDirectPathStoredFromLiveMessage` (SQLite, CGO) → PASS.
  - Suite complète sous Windows (CGO avec `zig cc`, Docker Desktop bloqué) : 66 PASS, 3 ignorés (permissions Unix).
- `2026-11-08` — **PASS sous Linux** : `gocheck.sh` sur ed16824 (contient les trois correctifs) → `GOCHECK_OK`, 69 tests et sous-tests PASS, aucun ignoré (les tests de permissions Unix tournent).
  - Non testé en réel : aucun compte WhatsApp connecté pendant les tests. À vérifier après déploiement : un média récent, puis un média plus ancien que la rétention du CDN.

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
