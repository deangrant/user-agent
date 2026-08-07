#!/usr/bin/env bash
# beforeShellExecution: ask before adding non-stdlib Go module deps (fail-open).
set -eu

input=$(cat || true)
cmd=""
if command -v node >/dev/null 2>&1; then
  cmd=$(printf '%s' "$input" | node -e '
    let raw = "";
    process.stdin.on("data", (c) => (raw += c));
    process.stdin.on("end", () => {
      try {
        const c = JSON.parse(raw).command;
        if (typeof c === "string") process.stdout.write(c);
      } catch {}
    });
  ' || true)
elif command -v python3 >/dev/null 2>&1; then
  cmd=$(printf '%s' "$input" | python3 -c '
import json, sys
try:
    c = json.load(sys.stdin).get("command") or ""
    if isinstance(c, str):
        sys.stdout.write(c)
except Exception:
    pass
' || true)
fi

if [ -z "${cmd:-}" ]; then
  printf '%s\n' '{ "permission": "allow" }'
  exit 0
fi

ask=false
reason=""

# go get with a module path (not bare "go get" / "go get -h")
if printf '%s' "$cmd" | grep -Eq -- '(^|[[:space:];|&])go[[:space:]]+get([[:space:]]|$)' \
  && printf '%s' "$cmd" | grep -Eq -- '[[:space:]][a-zA-Z0-9][a-zA-Z0-9._~+/-]*\.[a-zA-Z]{2,}(/|@|[[:space:]]|$)'; then
  ask=true
  reason="Command may add a third-party Go module (stdlib-only repo)."
fi

# go mod edit -require / -droprequire / replace that pulls external modules
if printf '%s' "$cmd" | grep -Eq -- 'go[[:space:]]+mod[[:space:]]+edit' \
  && printf '%s' "$cmd" | grep -Eq -- '(-require|-droprequire|-replace|-dropreplace)'; then
  ask=true
  if [ -n "$reason" ]; then
    reason="$reason Also edits go.mod require/replace."
  else
    reason="Command edits go.mod require/replace (stdlib-only repo)."
  fi
fi

if [ "$ask" = true ]; then
  if command -v python3 >/dev/null 2>&1; then
    reason_json=$(printf '%s' "$reason" | python3 -c 'import json,sys; print(json.dumps(sys.stdin.read()))')
  else
    reason_json='"Review dependency change (stdlib-only repo)."'
  fi
  printf '{"permission":"ask","user_message":%s,"agent_message":%s}\n' \
    "$reason_json" "$reason_json"
  exit 0
fi

printf '%s\n' '{ "permission": "allow" }'
exit 0
