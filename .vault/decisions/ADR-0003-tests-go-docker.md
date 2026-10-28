---
type: adr
id: ADR-0003
title: "Tests Go dans un conteneur Linux avec CGO"
status: accepted
created: 2026-10-29
tags: [adr]
---

# ADR-0003 — Tests Go dans un conteneur Linux avec CGO

## Contexte

`go-sqlite3` exige CGO et aucun compilateur C n'est installé sur le poste Windows ; de plus Windows ignore les modes POSIX (0600/0700) que la faille 06 doit vérifier.

## Décision

`go build`, `go vet` et `go test` sont exécutés dans `golang:1.27-bookworm` (CGO activé, gcc inclus) via `gocheck.sh`, avec des volumes de cache de modules. Les tests de modes POSIX se mettent en `t.Skip` sous Windows.

## Conséquences

- Build réel avec sqlite (pas le stub `!cgo`).
- Environnement identique à la CI d'un mainteneur Linux.

## Alternatives écartées

- Installer MSYS2/gcc sur le poste — plus lourd, et ne règle pas les modes POSIX.

## 🔗 Liens

- Sprints concernés : [[sprint-00-initialisation]], [[sprint-06-store-permissions]]
