# Parse UA

Reproduce User-Agent and Client Hints parsing with the project CLI. Use this
when debugging detection, merge, or a reported string.

## CLI

From the repo root:

```bash
go run ./cmd/useragent parse [flags] '<user-agent>'
go run ./cmd/useragent parse [flags] -   # HTTP headers on stdin
```

Useful flags:

- `-json` — indented JSON `Result`
- `-compact` — compact JSON
- `-sec-ch-ua`, `-sec-ch-ua-full-version-list`, `-sec-ch-ua-full-version`
- `-sec-ch-ua-platform`, `-sec-ch-ua-platform-version`
- `-sec-ch-ua-mobile`, `-sec-ch-ua-model`, `-sec-ch-ua-arch`
- `-sec-ch-ua-bitness`, `-sec-ch-ua-form-factors`, `-sec-ch-ua-wow64`

Hint flag values should match wire format (often quoted strings / SF lists),
e.g. `-sec-ch-ua-platform '"Windows"'`.

## Workflow

1. Run once without `-json` for a readable summary.
2. Run again with `-json` when comparing fields or pasting into a bug report.
3. For header-based repros, pipe a mini HTTP header block ending with a blank
   line into `parse -`.
4. Compare UA-only vs UA+hints when investigating `ClientHintsMismatch` or
   frozen Chromium versions.

## Examples

```bash
go run ./cmd/useragent parse \
  'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36'

go run ./cmd/useragent parse -json \
  -sec-ch-ua-platform '"Windows"' \
  -sec-ch-ua-platform-version '"15.0.0"' \
  'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36'

printf 'User-Agent: Mozilla/5.0\nSec-CH-UA-Platform: "Android"\n\n' \
  | go run ./cmd/useragent parse -json -
```
