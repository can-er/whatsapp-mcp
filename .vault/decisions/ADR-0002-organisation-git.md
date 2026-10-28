---
type: adr
id: ADR-0002
title: "Organisation Git : vault sur main du fork, une branche par faille en worktree"
status: accepted
created: 2026-10-29
tags: [adr]
---

# ADR-0002 — Organisation Git : vault sur main du fork, une branche par faille en worktree

## Contexte

Le vault (`.vault/`, `AGENTS.md`, `CLAUDE.md`) ne doit pas polluer les branches proposées en PR, mais doit rester visible et versionné pendant tout le travail.

## Décision

Le vault est versionné sur `main` du fork (racine du dépôt). Chaque faille a sa branche `fix/NN-slug` développée dans un worktree `.worktrees/NN` (exclu via `.git/info/exclude`). 01→04 sont empilées (même fichiers) ; 05 et 06 partent de `upstream/main` (PR indépendantes) ; 07 part de 04. Chaque branche est mergée dans `main` à la clôture du sprint.

## Conséquences

- Branches de PR propres, sans fichiers de pilotage.
- Conflits 03/06 et 07/05-06 résolus une seule fois, au merge dans `main`.

## Alternatives écartées

- Vault non versionné — perte de la traçabilité.
- Vault dans chaque branche — pollue les PR upstream.

## 🔗 Liens

- Sprints concernés : [[sprint-00-initialisation]], [[sprint-11-integration-rapport]]
