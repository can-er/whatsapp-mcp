# AGENTS.md — whatsapp-mcp (fork durci)

Fork de [lharries/whatsapp-mcp](https://github.com/lharries/whatsapp-mcp) : pont Go (whatsmeow)
+ serveur MCP Python pour WhatsApp. Ce fork corrige 7 failles de sécurité (API REST du pont,
lecture/écriture de fichiers, liste blanche, mode lecture seule, permissions, dépendances) pour un
usage par KS Esthetic (lecture automatisée de conversations choisies).

- **Remote Git** : `origin` = https://github.com/can-er/whatsapp-mcp (fork) ·
  `upstream` = https://github.com/lharries/whatsapp-mcp
- **Code** : [`whatsapp-bridge/`](whatsapp-bridge/) (Go) · [`whatsapp-mcp-server/`](whatsapp-mcp-server/) (Python, uv)
- **Pilotage projet** : dossier [`.vault/`](.vault/) (coffre Obsidian, voir ci-dessous).

> `AGENTS.md` est la source canonique des consignes agent, lue par tous les harness (Claude Code,
> OpenCode…). Le `CLAUDE.md` ne contient que `@AGENTS.md` : ne pas y ajouter de contenu.

---

## 1. Le coffre `.vault/` (Obsidian)

Tout le suivi du projet vit dans `.vault/`, organisé **à la manière d'Obsidian** :
frontmatter YAML, liens `[[wikilink]]`, tags `#sprint` / `#adr` / `#test-case`.

```
.vault/
├── Home.md                 # MOC (Map of Content) — point d'entrée
├── sprints/                # 1 fichier par sprint, avec un statut
├── test-cases/             # 1 cas de test par sprint
├── decisions/              # ADR — décisions d'architecture
├── templates/              # gabarits sprint / test-case / ADR
└── tools/                  # gocheck.sh, check_vault.py, set_status.py
```

**À lire en premier à chaque session : [`.vault/Home.md`](.vault/Home.md).**

### Statut des tâches — OBLIGATOIRE

Chaque sprint porte un champ `status` dans son frontmatter. Quand tu travailles sur un sprint,
**mets son statut à jour en temps réel** :

| Statut    | Signification                                              |
|-----------|-----------------------------------------------------------|
| `TODO`    | Pas encore commencé.                                       |
| `ONGOING` | En cours — passe le sprint à ce statut **avant** de coder. |
| `DONE`    | Terminé, testé et committé.                                |

> Règle : ne jamais laisser un sprint sur lequel tu travailles à `TODO`.
> Passe-le à `ONGOING` au démarrage, puis à `DONE` une fois le commit du sprint poussé.
> Mets aussi à jour le tableau de bord dans [`.vault/Home.md`](.vault/Home.md)
> (`python .vault/tools/set_status.py .vault NN ONGOING|DONE [--tc PASS]` le fait partout).

---

## 2. Cycle de travail d'un sprint

Pour **chaque** sprint, dans cet ordre :

1. **Ouvre le sprint** dans `.vault/sprints/` → passe `status: ONGOING`.
2. **Élabore / mets à jour un test spécifique** pour la fonctionnalité du sprint, dans
   `.vault/test-cases/` (un cas de test par sprint, lié au sprint via `[[wikilink]]`).
3. **Implémente** la fonctionnalité.
4. **Teste la fonctionnalité** en exécutant le cas de test du sprint. Ne committe **jamais**
   sans avoir fait passer le test au vert. Consigne le résultat dans le fichier test-case.
5. **Committe** en mentionnant le sprint (voir §3).
6. **Passe le sprint à `status: DONE`** et mets à jour `Home.md`.
7. **Émets un bip Windows** (voir §4).

Toute décision structurante → un ADR dans `.vault/decisions/`.

### Spécificités de ce dépôt

- **Organisation Git** ([ADR-0002](.vault/decisions/ADR-0002-organisation-git.md)) : le vault vit
  sur `main` du fork. Chaque faille a sa branche `fix/NN-slug`, développée dans un worktree
  `.worktrees/NN` (exclu de Git) ; les branches ne contiennent **jamais** `.vault/`, `AGENTS.md`
  ni `CLAUDE.md`. Elles sont mergées dans `main` à la clôture du sprint.
- **PR upstream** ([ADR-0001](.vault/decisions/ADR-0001-perimetre-upstream.md)) : uniquement
  sprints 05 et 06, branches basées sur `upstream/main`, en anglais, sans preuve d'exploitation.
- **Tests Go** ([ADR-0003](.vault/decisions/ADR-0003-tests-go-docker.md)) : toujours via
  `bash .vault/tools/gocheck.sh <chemin>/whatsapp-bridge` (build + vet + test, Linux, CGO).
  Tests Python : `cd whatsapp-mcp-server && uv run --with pytest pytest -q tests`.
- **Ne jamais lancer le pont connecté à un vrai compte WhatsApp** pendant les tests.
- **Aucun secret ni donnée personnelle réelle** dans le dépôt, même dans les tests et exemples
  ([ADR-0007](.vault/decisions/ADR-0007-donnees-de-test-fictives.md)) : numéros fictifs `1555…`
  uniquement. Hook `pre-commit` = `.vault/tools/check_no_pii.py --staged` ; lancer
  `python .vault/tools/check_no_pii.py --history` avant chaque push.

---

## 3. Convention de commit

Chaque commit **doit mentionner le sprint concerné**, au format :

```
[Sprint NN] <résumé impératif court>

<corps optionnel : ce qui a été fait, résultat du test>

Sprint: .vault/sprints/sprint-NN-<slug>.md
Test: .vault/test-cases/TC-NN-<slug>.md (PASS)
Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>
```

> Travaille sur une branche par sprint (`fix/NN-slug`), puis merge vers `main`.
> Ne pousse sur `main` que du code dont le test de sprint est PASS.

---

## 4. Bip Windows en fin de tâche — OBLIGATOIRE

**Quand tu termines une tâche** (sprint passé à `DONE`, ou jalon livré), émets un bip sonore
Windows pour signaler la fin :

```powershell
powershell -c "[console]::beep(880,250); [console]::beep(1175,400)"
```

(Deux notes ascendantes = tâche terminée. À lancer juste après le commit du sprint.)
