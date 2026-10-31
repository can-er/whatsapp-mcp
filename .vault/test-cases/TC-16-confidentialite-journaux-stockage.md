---
type: test-case
id: TC-16-confidentialite-journaux-stockage
sprint: sprint-16
status: TODO          # TODO | PASS | FAIL
created: 2026-10-31
last_run:
tags: [test-case]
---

# TC-16 — Aucun texte journalisé, liste blanche stricte, purge effective

> **sprint:: [[sprint-16-confidentialite-journaux-stockage]]**  ·  **status:: TODO**

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

- _(à remplir à l'exécution)_

## Commande(s)

```bash
bash .vault/tools/gocheck.sh <worktree>/whatsapp-bridge
```
