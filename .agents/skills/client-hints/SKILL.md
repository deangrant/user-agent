---
name: client-hints
description: >-
  Parse and merge User-Agent Client Hints in this library. Use when editing
  hintparse, merge, ClientHints types, Windows 11 platform-version mapping,
  CLI -sec-ch-* flags, or ClientHintsMismatch behavior.
trigger: >-
  Sec-CH-UA, Client Hints, hintparse, Full-Version-List, SignificantBrand,
  Platform-Version, Windows 11, ClientHintsMismatch, ParseHeaders, sec-ch-ua
---

# Client Hints

## Public surface

- `useragent.ClientHints` and `BrandVersion`
- `ClientHintsFromHeader` / `ClientHintsFromMap`
- `ParseWithHints`, `ParseHeaders`, `ParseRequest`
- Result flag: `ClientHintsMismatch`

## Headers (wire names)

Supported via `hintparse` / CLI flags:

- `Sec-CH-UA`, `Sec-CH-UA-Full-Version-List`, `Sec-CH-UA-Full-Version`
- `Sec-CH-UA-Platform`, `Sec-CH-UA-Platform-Version`
- `Sec-CH-UA-Mobile`, `Sec-CH-UA-Model`, `Sec-CH-UA-Arch`
- `Sec-CH-UA-Bitness`, `Sec-CH-UA-Form-Factors`, `Sec-CH-UA-WoW64`

Parsing is an **RFC 8941 structured-field subset** in `internal/hintparse`
(stdlib only—no full SF dependency).

## Brand selection

- Prefer Full-Version-List; fall back to Sec-CH-UA brands.
- `SignificantBrand` skips GREASE / “Not A;Brand” noise; prefers a concrete
  browser brand over bare Chromium when both appear.
- Merge normalizes names (e.g. “Google Chrome” → “Chrome”).
- Frozen UA versions like `120.0.0.0` yield to richer CH versions.

## OS / platform

- `internal/platform.ResolveFromCH` maps platform + platform-version.
- Windows Client Hints: major platform-version **≥ 13** maps to Windows 11
  (see comments/docs in `platform`); UA alone often still says `Windows NT 10.0`.
- OS **family** conflict (e.g. Android UA + Windows CH) → set
  `ClientHintsMismatch` and **do not** overwrite UA OS fields.
- Keep UA `iPadOS` when CH reports broader `iOS`.

## Device / form factors

- Model, arch, bitness, WoW64, mobile, and form-factors refine device fields in
  `merge.Apply` after detectors run.
- Form-factors can prefer Watch/TV/etc. over generic Mobile when present.

## CLI repro

Use `/parse-ua`. Pass wire-shaped values, including quotes where the header
uses sf-string form:

```bash
go run ./cmd/useragent parse -json \
  -sec-ch-ua-platform '"Windows"' \
  -sec-ch-ua-platform-version '"15.0.0"' \
  '<ua>'
```

Or pipe headers into `parse -`.
