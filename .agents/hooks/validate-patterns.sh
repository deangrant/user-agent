#!/usr/bin/env bash
# afterFileEdit: validate patterns.json parses as JSON (fail-open).
set -eu

input=$(cat || true)
file=""
if command -v node >/dev/null 2>&1; then
  file=$(printf '%s' "$input" | node -e '
    let raw = "";
    process.stdin.on("data", (c) => (raw += c));
    process.stdin.on("end", () => {
      try {
        const p = JSON.parse(raw).file_path;
        if (typeof p === "string" && p) process.stdout.write(p);
      } catch {}
    });
  ' || true)
elif command -v python3 >/dev/null 2>&1; then
  file=$(printf '%s' "$input" | python3 -c '
import json, sys
try:
    p = json.load(sys.stdin).get("file_path") or ""
    if isinstance(p, str):
        sys.stdout.write(p)
except Exception:
    pass
' || true)
fi

[ -z "${file:-}" ] && exit 0
[ ! -f "$file" ] && exit 0

case "$file" in
  */internal/data/patterns.json|internal/data/patterns.json) ;;
  *) exit 0 ;;
esac

ok=false
if command -v python3 >/dev/null 2>&1; then
  if python3 -c 'import json,sys; json.load(open(sys.argv[1]))' "$file" 2>/dev/null; then
    ok=true
  fi
elif command -v jq >/dev/null 2>&1; then
  if jq empty "$file" >/dev/null 2>&1; then
    ok=true
  fi
elif command -v node >/dev/null 2>&1; then
  if node -e 'JSON.parse(require("fs").readFileSync(process.argv[1],"utf8"))' "$file" 2>/dev/null; then
    ok=true
  fi
else
  # No validator available — fail open.
  exit 0
fi

if [ "$ok" != true ]; then
  printf '%s\n' "validate-patterns: invalid JSON in $file" >&2
fi
exit 0
