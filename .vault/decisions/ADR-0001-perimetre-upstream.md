---
type: adr
id: ADR-0001
title: "Périmètre upstream : fork complet, PR upstream pour 05 et 06 seulement"
status: accepted
created: 2026-10-29
tags: [adr]
---

# ADR-0001 — Périmètre upstream : fork complet, PR upstream pour 05 et 06 seulement

## Contexte

L'upstream (lharries/whatsapp-mcp) n'a rien mergé depuis avril 2025 (168 PR ouvertes, issue #220). Pas de SECURITY.md, signalement privé désactivé. Les failles 1-3 sont déjà publiques (issues #215, #241, #57) et couvertes par des PR ouvertes (#387, #224, #275, #271, #189, #143) ; la 4 par #319 ; la 7 par #343.

## Décision

Les 7 correctifs sont faits (avec tests) sur le fork `can-er/whatsapp-mcp`. Seules les failles 05 (mode lecture seule) et 06 (permissions de `store/`), non couvertes, font l'objet d'une PR upstream. Décision de Caner, 2026-10-29.

## Conséquences

- Pas de doublon dans un dépôt déjà saturé de PR.
- La divulgation coordonnée est sans objet pour 1-3 (déjà publiques).
- Le fork reste la version utilisée par KS Esthetic.

## Alternatives écartées

- 7 PR upstream (brief initial) — 5 doublons, chances de merge quasi nulles.
- Aucune PR — 05 et 06 apportent une valeur nouvelle à la communauté.

## 🔗 Liens

- Sprints concernés : [[sprint-00-initialisation]], [[sprint-05-read-only-mode]], [[sprint-06-store-permissions]]
