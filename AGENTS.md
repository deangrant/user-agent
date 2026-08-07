# Agent and contributor guidance

Structured conventions for AI agents and humans working in this repository. For
fuller context, see [README.md](README.md).

## Docs

- [`.agents/docs/ARCHITECTURE.md`](.agents/docs/ARCHITECTURE.md) — high-level system architecture and diagrams

## Rules

- [`.agents/rules/`](.agents/rules/) (symlinked from [`.cursor/rules`](.cursor/rules))
- [`.agents/rules/stdlib-only.mdc`](.agents/rules/stdlib-only.mdc) — always-on stdlib-only and testing norms
- [`.agents/rules/go-style-solid.mdc`](.agents/rules/go-style-solid.mdc) — Google Go style and SOLID Go design
- [`.agents/rules/patterns-catalog.mdc`](.agents/rules/patterns-catalog.mdc) — pattern catalog and detector ownership
- [`.agents/rules/client-hints-merge.mdc`](.agents/rules/client-hints-merge.mdc) — Client Hints parse and merge ownership

## Skills

- [`.agents/skills/`](.agents/skills/)
- [`.agents/skills/ua-detection-pipeline/`](.agents/skills/ua-detection-pipeline/) — detector order, State fill rules, composition root
- [`.agents/skills/patterns-catalog/`](.agents/skills/patterns-catalog/) — bots/apps/brands in `patterns.json`
- [`.agents/skills/client-hints/`](.agents/skills/client-hints/) — Sec-CH-UA headers, merge, CLI flags
- [`.agents/skills/google-go-style-guide/`](.agents/skills/google-go-style-guide/) — Google Go style
- [`.agents/skills/solid-go-design/`](.agents/skills/solid-go-design/) — SOLID design in Go

## Commands

- [`.agents/commands/`](.agents/commands/) (symlinked from [`.cursor/commands`](.cursor/commands))
- `/verify` — local CI checklist (`go test`, `golangci-lint`)
- `/parse-ua` — reproduce parsing via `cmd/useragent`
- `/add-pattern` — guided pattern catalog addition

## Hooks

- Config: [`.cursor/hooks.json`](.cursor/hooks.json)
- `afterFileEdit` → [`.agents/hooks/gofmt.sh`](.agents/hooks/gofmt.sh) formats edited `*.go` files
- `afterFileEdit` → [`.agents/hooks/golines.sh`](.agents/hooks/golines.sh) wraps Go lines to max-len 80
- `afterFileEdit` → [`.agents/hooks/validate-patterns.sh`](.agents/hooks/validate-patterns.sh) validates `patterns.json`
- `beforeShellExecution` → [`.agents/hooks/guard-stdlib-only.sh`](.agents/hooks/guard-stdlib-only.sh) asks before adding third-party Go modules
