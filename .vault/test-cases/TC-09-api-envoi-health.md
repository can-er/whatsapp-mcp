---
type: test-case
id: TC-09-api-envoi-health
sprint: sprint-09
status: TODO          # TODO | PASS | FAIL
created: 2026-10-29
last_run:
tags: [test-case]
---

# TC-09 — Envoi idempotent, audité, réservé au jeton d'envoi ; health OK

> **sprint:: [[sprint-09-api-envoi-health]]**  ·  **status:: TODO**

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

- _(à remplir à l'exécution)_

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
