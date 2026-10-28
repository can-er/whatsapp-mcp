---
type: test-case
id: TC-08-evenements-inbox-webhook
sprint: sprint-08
status: TODO          # TODO | PASS | FAIL
created: 2026-10-29
last_run:
tags: [test-case]
---

# TC-08 — Un message autorisé produit média + événement JSON + appel webhook

> **sprint:: [[sprint-08-evenements-inbox-webhook]]**  ·  **status:: TODO**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté : tests unitaires (événements construits à la main).

## Étapes

1. `go test -run 'TestInbox|TestWebhook|TestEvent'`.
2. Événement texte → fichier JSON valide conforme au schéma v1, aucun fichier `.tmp` résiduel.
3. Événement image (téléchargeur factice) → média écrit avant le JSON, `media.path` existe, sha256 correct ; téléchargeur en erreur → `status=failed`, `path=null`.
4. Webhook : serveur httptest reçoit le JSON + jeton ; 2 échecs puis succès → 3 appels ; URL non loopback refusée au démarrage.
5. Purge : fichiers plus vieux que la rétention supprimés, récents conservés.

## Résultat attendu

- Tous les tests PASS.

## Résultat obtenu

- _(à remplir à l'exécution)_

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
