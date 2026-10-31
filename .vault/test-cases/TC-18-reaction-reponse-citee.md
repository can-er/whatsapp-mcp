---
type: test-case
id: TC-18-reaction-reponse-citee
sprint: sprint-18
status: TODO          # TODO | PASS | FAIL
created: 2026-10-31
last_run:
tags: [test-case]
---

# TC-18 — Réactions et réponses citées visent le bon message, sous les mêmes garde-fous que l'envoi

> **sprint:: [[sprint-18-reaction-reponse-citee]]**  ·  **status:: TODO**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté. Données de test fictives uniquement ([[ADR-0007-donnees-de-test-fictives]]).

## Étapes

1. `go test -run 'TestReact|TestReply'` (SQLite, CGO) : réaction envoyée avec la discussion, l'ID et l'expéditeur stockés ; 403 hors liste ; 404 message inconnu ; 400 emoji vide ; doublon `client_msg_id` → un seul envoi, `duplicate: true` ; jeton du pont seul → 401.
2. Réponse citée : le message construit porte StanzaID, Participant (expéditeur du message cité, propre JID si is_from_me) et le texte cité ; message d'une autre discussion → 400 sans envoi ; média + citation → 400.
3. Journal d'audit : `kind` et `target_message_id` présents, ni texte ni emoji en clair.

## Résultat attendu

- Tests PASS.

## Résultat obtenu

- _(à remplir à l'exécution)_

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
