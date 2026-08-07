# Add Pattern

Add or adjust an entry in the embedded pattern catalog and prove it with tests.

## Checklist

1. **Identify the owner**
   - Bots / crawlers / HTTP clients → `bots` in
     `internal/data/patterns.json`, tests in `internal/detect/bot`
   - In-app browsers / apps → `apps`, tests in `internal/detect/app`
   - Device brand/model → `brands`, tests in `internal/detect/device`

2. **Edit `internal/data/patterns.json`**
   - Keep JSON valid (trailing commas forbidden).
   - Prefer specific substrings; avoid broad tokens that match normal browsers.
   - For bots: set `heuristic: true` only when the pattern alone is too weak and
     must also pass `looksLikeBot` (see bot detector).
   - For apps: `.*` in `pattern` enables a case-insensitive regexp; otherwise
     matching is substring-based.
   - For brands: prefer a distinctive `prefix` or `pattern`; order matters for
     overlapping prefixes (e.g. more specific before shorter).

3. **Add a table-driven test** in the owning detector package with a realistic
   UA fixture. Failure messages: **got before want**.

4. **Verify**
   ```bash
   go test ./internal/detect/<owner>/...
   ```
   Then run `/verify` (`go test ./...` and `golangci-lint run ./...`).

5. **Optional**: reproduce with `/parse-ua` for the new fixture.

## Do not

- Add third-party dependencies for pattern matching.
- Put detector wiring or catalog loading outside the existing packages.
- Rely on Client Hints for catalog coverage; patterns are UA-side.
