---
name: patterns-catalog
description: >-
  Author and match bots, apps, and brands in the embedded patterns.json catalog.
  Use when adding detection coverage, debugging false positives, or editing
  bot/app/device pattern matching.
trigger: >-
  patterns.json, BotPattern, AppPattern, BrandPattern, heuristic bot,
  looksLikeBot, add pattern, false positive UA match, device brand prefix
---

# Patterns catalog

Embedded catalog: `internal/data/patterns.json`, loaded by `internal/data`.

## Schema

### `bots[]`

| Field | Role |
| ----- | ---- |
| `pattern` | Case-insensitive substring of the UA |
| `name` | Agent display name |
| `deviceClass` / `agentClass` | Classification (`Robot`, `Cloud`, …) |
| `heuristic` | If true, match only when `looksLikeBot` also agrees |

First **non-heuristic** contains-match wins. Heuristic matches are deferred and
applied only if no hard match fired and the UA still looks like a bot.

### `apps[]`

| Field | Role |
| ----- | ---- |
| `pattern` | Substring needle, or regexp if it contains `.*` |
| `name` / `agentClass` / `deviceClass` / `product` | Identity |

Regexp entries compile as `(?i)` + pattern. Browser-only rows (Browser class
with empty device class) are skipped so they do not steal browser detection.

### `brands[]`

| Field | Role |
| ----- | ---- |
| `prefix` or `pattern` | Token-boundary key (prefix preferred) |
| `brand` / `name` / `class` | Device brand, model name, optional class |

Matching uses token boundaries so `sm-r` can win over `sm-`. Put more specific
overlapping keys **before** shorter ones in the JSON array.

## Safe authoring

- Prefer distinctive tokens (`googlebot`, `headlesschrome`) over short English
  words (`bot`, `java/`).
- Mark weak tokens `heuristic: true` and rely on `looksLikeBot`.
- Always add a table-driven test with a realistic UA (and a negative case for
  risky patterns).
- After edits, JSON must parse; the `validate-patterns` hook checks this.
- Use `/add-pattern` for the end-to-end checklist and `/verify` before finishing.

## Ownership

| Catalog section | Detector package |
| --------------- | ---------------- |
| `bots` | `internal/detect/bot` |
| `apps` | `internal/detect/app` |
| `brands` | `internal/detect/device` |
