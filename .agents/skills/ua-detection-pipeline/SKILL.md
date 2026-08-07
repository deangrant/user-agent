---
name: ua-detection-pipeline
description: >-
  Explain and extend the User-Agent detection pipeline for this repo. Use when
  adding detectors, fixing fill/overwrite bugs, wiring Analyzer, or reasoning
  about bot/app/agent/engine/opsys/device order and merge.
trigger: >-
  detector pipeline, detect.State, Analyzer, defaultDetectors, bot before app,
  merge.Apply, composition root, Detect method, Client Hints overlay
---

# UA detection pipeline

Use this skill when changing how analysis is wired or how detectors cooperate.

## Composition root

- Public API: package `useragent` (`Parse`, `ParseWithHints`, `ParseHeaders`,
  `ParseRequest`, `Analyzer`).
- Pipeline is built only in `useragent.defaultDetectors` / `NewAnalyzer`.
- Do **not** create `detect` → subpackage → `detect` import cycles.
- `detect.Detector` is a single method: `Detect(state *State)`.

## Order (significant)

1. **tokenize** — `tokenize.Parse(ua)` into `State.Tokens` (before detectors).
2. **bot** — `internal/detect/bot`
3. **app** — `internal/detect/app` (bots win; bot runs first)
4. **agent** — `internal/detect/agent`
5. **engine** — `internal/detect/engine` (may read `AgentName`, e.g. Blink)
6. **opsys** — `internal/detect/opsys` (package name avoids stdlib `os`)
7. **device** — `internal/detect/device` (may read OS fields)
8. **merge** — `merge.Apply(state)` overlays Client Hints and finalizes
   derived strings (`NameVersion`, majors, etc.)

## Fill rules

- Detectors fill **empty** fields unless they own the concern.
- Later detectors should not casually clobber stronger earlier results
  (especially bot/app identity).
- Client Hints merge **may overwrite** UA-derived agent/OS/device values when
  CH is present and consistent; OS family conflicts set `ClientHintsMismatch`
  and keep UA OS.

## Package map

| Concern | Package |
| ------- | ------- |
| Pattern tables | `internal/data` (`patterns.json` embed) |
| Tokenize UA | `internal/tokenize` |
| CH structured fields | `internal/hintparse` |
| CH → OS mapping | `internal/platform` |
| Merge overlay | `internal/merge` |
| Class constants | `internal/uaclass` |
| Version compare | `internal/version` |
| CLI | `cmd/useragent` |

## When adding a detector

- Implement `detect.Detector` in its own package under `internal/detect/`.
- Wire it in `useragent.defaultDetectors` in the correct order.
- Document any new cross-field reads (like engine→agent or device→OS).
- Prefer constructor injection of catalog data over new package globals.
- Keep the module stdlib-only.
