---
type: sprint
id: sprint-12
title: "Appairage par code (WHATSAPP_PAIR_PHONE) livré par inbox + webhook (fork)"
status: DONE          # TODO | ONGOING | DONE
owner: CK
created: 2026-10-29
closed: 2026-10-29
tags: [sprint]
---

# Sprint 12 — Appairage par code (WHATSAPP_PAIR_PHONE) livré par inbox + webhook (fork)

> **status:: DONE**  ·  **test:: [[TC-12-code-appairage]]**

## 🎯 Objectif

Avec `WHATSAPP_PAIR_PHONE`, un pont non appairé obtient un code d'appairage (`PairPhone`) au lieu du QR et le livre en événement `pairing_code` (inbox + webhook, indice = 2 derniers chiffres), régénéré à chaque session de connexion (≥ 60 s d'intervalle) jusqu'à l'appairage, puis émet `paired`.

## 📋 Tâches

- [x] `pairingNotifier` : format des événements, `phone_hint`, intervalle minimal 60 s
- [x] Livraison générique `eventSink.deliver` (inbox atomique + webhook)
- [x] Boucle d'appairage : nouvelle session tant que non appairé, délai global `WHATSAPP_PAIR_TIMEOUT_MINUTES` (30)
- [x] Mode QR inchangé sans la variable
- [x] README

## 🧪 Test associé (obligatoire)

Chaque sprint a **un test spécifique** qui prouve la fonctionnalité. À exécuter **avant le commit**.

→ [[TC-12-code-appairage]]

## ✅ Definition of Done

- [x] Tâches terminées
- [x] Test [[TC-12-code-appairage]] **PASS**
- [x] Commit fait, message préfixé `[Sprint 12]`
- [x] Statut de ce sprint passé à `DONE` + Home.md mis à jour
- [x] Bip Windows émis

## 🧾 Journal

- `2026-10-29` — Sprint créé à la demande de ks-support-agent-6b (spécification reçue, cf. ADR-0006).
- `2026-10-29` — Appairage par code livré par inbox + webhook ; TC-12 PASS ; mergé dans main.

## 🔗 Liens

- ADR : [[ADR-0006-interface-pont-pipeline]]
