"""Block real phone numbers, WhatsApp JIDs and secrets from entering the repository.

Usage:
  check_no_pii.py [--staged]      added lines of the staged diff (pre-commit hook, default)
  check_no_pii.py --history       added lines of every commit reachable from any ref
  check_no_pii.py FILE...         whole files

Test data must use obviously fictitious numbers: 1555xxxxxxx, 32400000xxx, 123456789..., 999...
Exit code 1 when something is found.
"""
import re
import subprocess
import sys

sys.stdout.reconfigure(encoding="utf-8")

FICTITIOUS = re.compile(r"^(1555\d*|32400000\d*|123456789\d*|9+|0+|1+)$")

PATTERNS = [
    ("Belgian mobile number", re.compile(r"(?<![\w.])(?:\+|00)?32[\s.]?4\d{2}(?:[\s.]?\d{2}){3}(?!\d)")),
    ("Belgian mobile number (national)", re.compile(r"(?<![\w.])04\d{2}(?:[\s./]?\d{2}){3}(?!\d)")),
    ("international phone number", re.compile(r"(?<![\w.])\+\d{2,3}(?:[\s.]?\d){8,13}(?!\d)")),
    ("WhatsApp JID", re.compile(r"(?<![\w.])(\d{8,15})(?::\d+)?@(?:s\.whatsapp\.net|c\.us|lid)\b")),
    ("bare phone number", re.compile(r"(?<![\w.:/])(?:32|33|31|49|44|1)\d{8,10}(?![\w.])")),
    ("Anthropic API key", re.compile(r"sk-ant-[A-Za-z0-9_-]{20,}")),
    ("Telegram bot token", re.compile(r"\b\d{8,10}:[A-Za-z0-9_-]{35}\b")),
    ("GitHub token", re.compile(r"\bgh[pousr]_[A-Za-z0-9]{36,}\b")),
    ("private key", re.compile(r"-----BEGIN [A-Z ]*PRIVATE KEY-----")),
]

SKIP_PATHS = ("check_no_pii.py", "go.sum", "uv.lock", ".png")


def digits(s):
    return re.sub(r"\D", "", s)


def findings(line):
    out = []
    for label, rx in PATTERNS:
        for m in rx.finditer(line):
            number = m.group(1) if rx.groups else m.group(0)
            d = digits(number)
            if label not in ("Anthropic API key", "Telegram bot token", "GitHub token", "private key") and FICTITIOUS.match(d):
                continue
            out.append((label, m.group(0)))
    return out


def scan_diff(diff):
    found, path, commit = [], "?", ""
    for line in diff.splitlines():
        if line.startswith("commit "):
            commit = line.split()[1][:7] + " "
        elif line.startswith("+++ "):
            path = line[6:] if line.startswith("+++ b/") else line[4:]
        elif line.startswith("+") and not path.endswith(SKIP_PATHS):
            for label, text in findings(line[1:]):
                found.append(f"{commit}{path}: {label}: {text}")
    return found


def git(*args):
    return subprocess.run(["git", *args], capture_output=True, text=True, encoding="utf-8", errors="replace", check=True).stdout


def main(argv):
    if argv and argv[0] == "--history":
        found = scan_diff(git("log", "--all", "-p", "--no-color", "--format=commit %H"))
    elif argv and argv[0] != "--staged":
        found = []
        for f in argv:
            if f.endswith(SKIP_PATHS):
                continue
            with open(f, encoding="utf-8", errors="replace") as fh:
                for n, line in enumerate(fh, 1):
                    found += [f"{f}:{n}: {label}: {text}" for label, text in findings(line)]
    else:
        found = scan_diff(git("diff", "--cached", "-U0", "--no-color"))
    if found:
        print("Possible personal data or secret found (use fictitious values such as 15550000001):")
        for f in found:
            print("  " + f)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
