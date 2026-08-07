#!/usr/bin/env bash
# afterFileEdit: wrap edited Go files with golines max-len 80 (fail-open).
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
  *.go) ;;
  *) exit 0 ;;
esac

golines_bin=""
if command -v golines >/dev/null 2>&1; then
  golines_bin=$(command -v golines)
fi

if [ -n "$golines_bin" ]; then
  "$golines_bin" -w -m 80 --base-formatter=gofmt "$file" || true
fi
exit 0
