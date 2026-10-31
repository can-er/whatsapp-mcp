---
type: test-case
id: TC-18-reaction-reponse-citee
sprint: sprint-18
status: PASS          # TODO | PASS | FAIL
created: 2026-10-31
last_run: 2026-10-31
tags: [test-case]
---

# TC-18 — Réactions et réponses citées visent le bon message, sous les mêmes garde-fous que l'envoi

> **sprint:: [[sprint-18-reaction-reponse-citee]]**  ·  **status:: PASS**

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

- `2026-10-31` — **PASS** : `bash .vault/tools/gocheck.sh .worktrees/18/whatsapp-bridge -v` (Docker golang:1.27-bookworm, CGO) → build + vet OK, 54 tests PASS dont TestReactAPI (envoi unique avec doublon `client_msg_id`, 401 jeton du pont, 403 hors liste, 404 message inconnu ou d'une autre discussion, 400 emoji vide/long/ID manquant, audit `kind`/`target_message_id` sans emoji), TestReplyAPI (groupe + discussion directe donnée par numéro, 400 autre discussion / inconnu / média, 403 hors liste, texte absent de l'audit) et TestReplyAndReactionMessages (clé : participant pour un groupe, from_me pour un message propre, rien en direct ; ContextInfo StanzaID/Participant/QuotedMessage). TestNoMessageTextInLogs reste vert.
- Non testé : envoi réel à WhatsApp (pas de compte en test), à vérifier sur la VM ; en particulier le participant d'un groupe en adressage LID.

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
