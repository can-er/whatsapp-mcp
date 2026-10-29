"""Update a sprint's status everywhere in the vault.

Usage: set_status.py <vault> <NN> <TODO|ONGOING|DONE> [--tc PASS|FAIL] [--result FILE] [--log "journal entry"]
"""
import argparse
import re
from datetime import date
from pathlib import Path

ap = argparse.ArgumentParser()
ap.add_argument("vault")
ap.add_argument("n")
ap.add_argument("status")
ap.add_argument("--tc")
ap.add_argument("--result", help="file whose content replaces the 'Résultat obtenu' section")
ap.add_argument("--log")
a = ap.parse_args()
V = Path(a.vault)
today = date.today().isoformat()

sp = next((V / "sprints").glob(f"sprint-{a.n}-*.md"))
slug = sp.stem
t = sp.read_text(encoding="utf-8")
t = re.sub(r"^status: \w+", f"status: {a.status}", t, count=1, flags=re.M)
t = re.sub(r"\*\*status:: \w+\*\*", f"**status:: {a.status}**", t, count=1)
if a.status == "DONE":
    t = re.sub(r"^closed:.*$", f"closed: {today}", t, count=1, flags=re.M)
    t = t.replace("- [ ] ", "- [x] ")
if a.log:
    t = re.sub(r"(## 🧾 Journal\n\n(?:- .*\n)*)", lambda m: m.group(1) + f"- `{today}` — {a.log}\n", t, count=1)
sp.write_text(t, encoding="utf-8", newline=chr(10))

if a.tc:
    tc = next((V / "test-cases").glob(f"TC-{a.n}-*.md"))
    u = tc.read_text(encoding="utf-8")
    u = re.sub(r"^status: \w+", f"status: {a.tc}", u, count=1, flags=re.M)
    u = re.sub(r"\*\*status:: \w+\*\*", f"**status:: {a.tc}**", u, count=1)
    u = re.sub(r"^last_run:.*$", f"last_run: {today}", u, count=1, flags=re.M)
    if a.result:
        body = Path(a.result).read_text(encoding="utf-8").strip()
        u = re.sub(r"(## Résultat obtenu\n\n).*?(\n## )", lambda m: m.group(1) + body + "\n" + m.group(2), u, count=1, flags=re.S)
    tc.write_text(u, encoding="utf-8", newline=chr(10))

for f in (V / "Home.md", V / "sprints" / "_index.md"):
    h = f.read_text(encoding="utf-8")
    h = re.sub(rf"(\| \[\[{re.escape(slug)}\]\] \| )`\w+`", rf"\1`{a.status}`", h)
    h = re.sub(r"^updated:.*$", f"updated: {today}", h, count=1, flags=re.M)
    f.write_text(h, encoding="utf-8", newline=chr(10))
print(f"{slug} -> {a.status}" + (f", TC -> {a.tc}" if a.tc else ""))
