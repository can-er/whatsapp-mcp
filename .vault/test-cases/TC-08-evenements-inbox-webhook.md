---
type: test-case
id: TC-08-evenements-inbox-webhook
sprint: sprint-08
status: PASS          # TODO | PASS | FAIL
created: 2026-10-29
last_run: 2026-10-29
tags: [test-case]
---

# TC-08 — Un message autorisé produit média + événement JSON + appel webhook

> **sprint:: [[sprint-08-evenements-inbox-webhook]]**  ·  **status:: PASS**

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

- `2026-10-29` — **PASS** (branche `fix/08-evenements-inbox-webhook`, commit `3f7c9db`)
- `TestEventTextMessage` : JSON v1 complet et conforme, nom `20261029T093000Z_<id>.json`, aucun fichier tmp résiduel, mode 0640.
- `TestEventLIDSenderResolvedToPhone` : chat/sender LID → `chat_jid`/`sender_phone` au format numéro + `chat_lid`/`sender_lid` ; LID inconnu → `sender_phone=null`.
- `TestEventMediaWrittenBeforeEvent` : média écrit avant l'événement, sha256 / voice / duration_s / chemin absolu corrects. `TestEventMediaDownloadFailure` : `status=failed`, `path=null`, aucun fichier.
- `TestWebhookRetriesWithToken` : 503, 503 puis 204 → 3 appels, Bearer + JSON. `TestWebhookURLMustBeLoopback` : URL non loopback refusée au démarrage. `TestInboxRetentionPurge` : >30 j et tmp >1 h supprimés, récent conservé.
- Suite complète `GOCHECK_OK` (Linux, CGO) ; `go vet` natif Windows OK ; `check_no_pii.py` : 0.

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
