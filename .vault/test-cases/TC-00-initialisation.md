---
type: test-case
id: TC-00-initialisation
sprint: sprint-00
status: PASS          # TODO | PASS | FAIL
created: 2026-10-29
last_run: 2026-10-29
tags: [test-case]
---

# TC-00 — Le vault est complet et navigable, la chaîne de test fonctionne

> **sprint:: [[sprint-00-initialisation]]**  ·  **status:: PASS**

## Préconditions

- Dépôt cloné, vault créé.

## Étapes

1. Vérifier que chaque `[[wikilink]]` du vault pointe vers un fichier existant.
2. Vérifier que Home.md liste les 9 sprints.
3. Lancer `gocheck.sh` sur le code upstream non modifié.

## Résultat attendu

- 0 lien cassé, 9 sprints dans Home.md, `GOCHECK_OK`.

## Résultat obtenu

- `2026-10-29` — **PASS**
- `check_vault.py` : `9 sprints, 0 liens cassés, 0 sprints absents de Home.md`
- `gocheck.sh whatsapp-bridge` (upstream 7d6a06d, non modifié) : `[no test files]` puis `GOCHECK_OK`
- Upstream : pas de SECURITY.md, `private-vulnerability-reporting.enabled=false` ; failles 1-3/4/7 déjà couvertes par des PR ouvertes (cf. [[ADR-0001-perimetre-upstream]]).
- Fork créé : https://github.com/can-er/whatsapp-mcp

## Commande(s)

```bash
python ".vault/tools/check_vault.py" .vault
bash ".vault/tools/gocheck.sh" <worktree>/whatsapp-bridge   # go build + vet + test (Linux, CGO), golang:1.27-bookworm
```
