# Verify

Run the local quality gate for this repository. Fix failures before considering
Go work done.

## Steps

1. From the repo root, run:

```bash
go test ./...
```

2. Then run:

```bash
golangci-lint run ./...
```

3. If `lll` (or golines-related) failures report lines over 80 characters, wrap
   or reformat with golines (`max-len: 80`, see `.golangci.yml`). Prefer
   `gofmt`/`goimports` style already enforced by the project formatters.

## Success criteria

- Both commands exit 0.
- Do not skip failing packages or disable linters to “pass.”
- Report any remaining failures with package/file context.
