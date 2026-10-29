---
type: adr
id: ADR-0007
title: "Données de test fictives et garde-fou pre-commit (incident numéros réels)"
status: accepted
created: 2026-10-29
tags: [adr, incident]
---

# ADR-0007 — Données de test fictives et garde-fou pre-commit

## Contexte

Incident du 2026-10-29 : les sprints 04 et 07 ont repris dans les tests (Go et Python) les vrais
numéros de deux contacts KS Esthetic fournis dans le brief, ainsi que des numéros d'exemple
réalistes (README, commentaire de `chat_allowlist.go`). Ces commits ont été poussés sur le fork
public `can-er/whatsapp-mcp` (branches `main`, `fix/04-allowed-chats`, `fix/07-bump-whatsmeow`).
Signalé par ks-support-agent-6b à la demande de Caner.

## Décision

- Les numéros sont remplacés par des valeurs manifestement fictives (`1555000000x`) dans tout
  l'historique (`git filter-repo --replace-text`, limité aux commits postérieurs à `fix/03`,
  `fix/05` et `fix/06` pour ne pas toucher l'historique upstream ni les branches des PR
  #388 / #389), puis force-push avec `--force-with-lease` (accord de Caner).
- **Règle** : aucune donnée personnelle réelle (numéro, JID, nom de contact) dans le dépôt, même dans
  les tests ou les exemples. Valeurs autorisées : `1555…`, `32400000…`, `123456789…`, `999…`.
  La configuration réelle (`WHATSAPP_ALLOWED_CHATS`…) ne vit que dans l'environnement de la VM.
- **Garde-fou** : `.vault/tools/check_no_pii.py` (numéros belges / internationaux, JID, clés
  `sk-ant`, jetons Telegram / GitHub, clés privées), branché en hook `pre-commit` commun à tous les
  worktrees (`.git/hooks/pre-commit`), et `--history` lancé avant chaque push.

## Conséquences

- Historique public propre ; `check_no_pii.py --history` → 0 occurrence.
- Les anciens commits restent accessibles par SHA sur GitHub tant que le support ne les a pas
  purgés : tips d'origine `082c904` (main), `39b8945` (fix/04), `8d8c48c` (fix/07), puis
  `b369aad`, `a193dee`, `c1f80ab` (deuxième passe, numéro d'exemple seulement). Demande de purge à
  faire par Caner (GitHub Support → « Remove sensitive data »).
- Les PR upstream #388 et #389 n'ont jamais contenu ces données (vérifié).

## Alternatives écartées

- Simple commit correctif — les numéros restaient dans l'historique public.
- Supprimer et recréer le fork — fermerait les PR #388 / #389.

## 🔗 Liens

- Sprints concernés : [[sprint-04-allowed-chats]], [[sprint-07-bump-whatsmeow]]
