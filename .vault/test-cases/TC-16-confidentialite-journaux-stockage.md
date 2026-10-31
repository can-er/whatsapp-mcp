---
type: test-case
id: TC-16-confidentialite-journaux-stockage
sprint: sprint-16
status: PASS          # TODO | PASS | FAIL
created: 2026-10-31
last_run: 2026-10-31
tags: [test-case]
---

# TC-16 — Aucun texte journalisé, liste blanche stricte, purge effective

> **sprint:: [[sprint-16-confidentialite-journaux-stockage]]**  ·  **status:: PASS**

## Préconditions

- Docker Desktop démarré, image `golang:1.27-bookworm`.
- Aucun compte WhatsApp réel connecté. Données de test fictives uniquement ([[ADR-0007-donnees-de-test-fictives]]).

## Étapes

1. Test de source : aucun appel de journalisation dans le code du pont ne reçoit la variable `content` ni le texte d'un message.
2. `go test -run TestChatAllowlistStrict` : entrée invalide ou liste sans entrée valide → erreur ; liste vide → nil.
3. `go test -run TestPurgeDisallowedChats` (SQLite, CGO) : chats/messages hors liste supprimés (y compris JID illisible), chats autorisés intacts, dossiers médias hors liste supprimés, autres fichiers du store intacts, contenu supprimé absent du fichier après VACUUM.

## Résultat attendu

- Tests PASS.

## Résultat obtenu

- `2026-10-31` — **PASS** : `bash .vault/tools/gocheck.sh .worktrees/16/whatsapp-bridge -v` (Docker golang:1.27-bookworm, CGO) → build + vet OK, 49 tests PASS dont TestNoMessageTextInLogs, TestChatAllowlistStrict, TestPurgeDisallowedChats (3 chats / 9 messages / 1 dossier média supprimés, chat autorisé intact, texte supprimé absent de messages.db*).
- Contre-épreuve : une ligne `logger.Debugf("leak %s", content)` injectée dans la synchro d'historique fait échouer TestNoMessageTextInLogs (main.go:1334), puis retirée.
- Constat : le filtrage du stockage par la liste blanche existait déjà (sprint 04), en amont des lignes « Message content ». Leur présence implique une liste blanche vide ou entièrement invalide sur l'unité de prod (comportement « tout autoriser », désormais refusé au démarrage). À confirmer par la ligne « Chat allowlist enabled » du journal VM.

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
