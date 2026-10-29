---
type: moc
title: whatsapp-mcp (durcissement sécurité) — Vault
tags: [moc, home]
updated: 2026-10-29
---

# 🗂️ whatsapp-mcp — Vault de pilotage

Vault de type **Obsidian** : pilotage des **sprints**, des **test cases** et des **décisions
d'architecture** du durcissement sécurité de `whatsapp-mcp` (fork `can-er/whatsapp-mcp`).
Navigation par liens `[[...]]`.

> Ce vault est la **source de vérité du suivi de dev**.

## 🚀 Sprints

→ Board complet : [[sprints/_index|Board des sprints]]

| Sprint | Statut | Test |
|---|---|---|
| [[sprint-00-initialisation]] | `DONE` | [[TC-00-initialisation]] |
| [[sprint-01-api-rest-auth]] | `DONE` | [[TC-01-api-rest-auth]] |
| [[sprint-02-media-path-allowlist]] | `DONE` | [[TC-02-media-path-allowlist]] |
| [[sprint-03-download-path-traversal]] | `DONE` | [[TC-03-download-path-traversal]] |
| [[sprint-04-allowed-chats]] | `DONE` | [[TC-04-allowed-chats]] |
| [[sprint-05-read-only-mode]] | `DONE` | [[TC-05-read-only-mode]] |
| [[sprint-06-store-permissions]] | `DONE` | [[TC-06-store-permissions]] |
| [[sprint-07-bump-whatsmeow]] | `DONE` | [[TC-07-bump-whatsmeow]] |
| [[sprint-08-evenements-inbox-webhook]] | `DONE` | [[TC-08-evenements-inbox-webhook]] |
| [[sprint-09-api-envoi-health]] | `DONE` | [[TC-09-api-envoi-health]] |
| [[sprint-10-evenements-edit-revoke-reaction]] | `DONE` | [[TC-10-evenements-edit-revoke-reaction]] |
| [[sprint-11-integration-rapport]] | `DONE` | [[TC-11-integration-rapport]] |
| [[sprint-12-code-appairage]] | `DONE` | [[TC-12-code-appairage]] |
| [[sprint-13-api-history]] | `TODO` | [[TC-13-api-history]] |

## ✅ Test cases

→ Index : [[test-cases/_index|Index des test cases]]

## 🧭 Décisions (ADR)

- [[ADR-0001-perimetre-upstream]] — Périmètre upstream : fork complet, PR upstream pour 05 et 06 seulement
- [[ADR-0002-organisation-git]] — Organisation Git : vault sur main du fork, une branche par faille en worktree
- [[ADR-0003-tests-go-docker]] — Tests Go dans un conteneur Linux avec CGO
- [[ADR-0004-auth-pont-jeton]] — Authentification du pont par jeton Bearer partagé
- [[ADR-0005-envoi-media-desactive-par-defaut]] — Envoi de médias désactivé tant que WHATSAPP_MEDIA_DIR n'est pas défini
- [[ADR-0006-interface-pont-pipeline]] — Interface pont WhatsApp → pipeline KS (événements v1, API d'envoi)
- [[ADR-0007-donnees-de-test-fictives]] — Données de test fictives et garde-fou pre-commit (incident numéros réels)

## 🧩 Templates

- [[sprint-template]]
- [[test-case-template]]
- [[adr-template]]

## Convention de statut

`TODO` → `ONGOING` → `DONE`. **On bascule le statut d'un sprint dès qu'on commence à travailler
dessus** (→ `ONGOING`) et à la clôture (→ `DONE`, une fois le test passé et le commit fait).
Cf. `AGENTS.md`.
