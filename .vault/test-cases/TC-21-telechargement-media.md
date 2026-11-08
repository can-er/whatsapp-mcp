---
type: test-case
id: TC-21-telechargement-media
sprint: sprint-21
status: TODO          # TODO | PASS | FAIL
created: 2026-11-08
last_run:
tags: [test-case]
---

# TC-21 — Le téléchargement d'un média réussit avec le chemin signé, la copie de l'inbox ou un renouvellement

> **sprint:: [[sprint-21-telechargement-media]]**  ·  **status:: TODO**

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

- _(à remplir à l'exécution)_

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
