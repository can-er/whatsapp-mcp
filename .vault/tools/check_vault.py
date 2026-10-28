"""Check that every [[wikilink]] in the vault resolves and Home.md lists every sprint."""
import re
import sys
from pathlib import Path

sys.stdout.reconfigure(encoding="utf-8")

V = Path(sys.argv[1])
names = {p.stem for p in V.rglob("*.md")} | {str(p.relative_to(V).with_suffix("")).replace("\\", "/") for p in V.rglob("*.md")}
INLINE_CODE = re.compile(r"`[^`\n]*`")
broken = []
for p in V.rglob("*.md"):
    if p.parent.name == "templates":
        continue
    text = INLINE_CODE.sub("", p.read_text(encoding="utf-8"))
    for m in re.finditer(r"\[\[([^\]|#]+)", text):
        if m.group(1).strip() not in names:
            broken.append(f"{p.relative_to(V)} -> {m.group(1)}")
sprints = sorted(p.stem for p in (V / "sprints").glob("sprint-*.md"))
home = (V / "Home.md").read_text(encoding="utf-8")
missing = [s for s in sprints if f"[[{s}]]" not in home]
print(f"{len(sprints)} sprints, {len(broken)} liens cassés, {len(missing)} sprints absents de Home.md")
for b in broken + missing:
    print("  ", b)
sys.exit(1 if broken or missing else 0)
