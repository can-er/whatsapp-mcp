---
type: test-case
id: TC-09-api-envoi-health
sprint: sprint-09
status: PASS          # TODO | PASS | FAIL
created: 2026-10-29
last_run: 2026-10-29
tags: [test-case]
---

# TC-09 — Envoi idempotent, audité, réservé au jeton d'envoi ; health OK

> **sprint:: [[sprint-09-api-envoi-health]]**  ·  **status:: PASS**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté : tests unitaires (événements construits à la main).

## Étapes

1. `go test -run 'TestSendToken|TestIdempotency|TestAudit|TestHealth'` (envoi WhatsApp remplacé par une fonction factice).
2. Jeton de lecture sur `/api/send` → 401 quand `WHATSAPP_SEND_TOKEN` est défini ; jeton d'envoi → 200.
3. Même `client_msg_id` deux fois → 1 seul envoi factice, même `message_id`.
4. Ligne d'audit écrite avec le sha256 du texte, sans le texte en clair.
5. `/api/health` → JSON attendu.

## Résultat attendu

- Tous les tests PASS.

## Résultat obtenu

- `2026-10-29` — **PASS** (branche `fix/09-api-envoi-health`, commit `1cd4350`)
- `TestSendTokenSeparation` : jeton du pont sur `/api/send` (et `/api//send`) → refusé ; jeton d'envoi → 200 avec `message_id` ; jeton d'envoi sur `/api/health` → 401.
- `TestSendIdempotency` : même `client_msg_id` → 1 seul envoi, même `message_id`, `duplicate: true` ; 10 tentatives concurrentes → 1 envoi. `TestSendFailureIsNotCached` : échec puis nouvel essai → renvoi.
- `TestSendAuditLog` : ligne JSONL avec sha256 du texte (texte absent du fichier), mode 0600, idempotence conservée après réouverture. `TestSendRespectsAllowlist` : 403, aucun envoi. `TestHealth` : `connected`, `logged_in`, `last_event_at` ; 401 sans jeton.
- Suite complète 29 tests `GOCHECK_OK` ; pytest 30 passed ; en-têtes Python (jeton d'envoi) vérifiés.

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
