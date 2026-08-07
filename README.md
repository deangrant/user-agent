# useragent

This package analyzes an HTTP User-Agent string. You can also pass User-Agent
Client Hints (`Sec-CH-UA*`). The package returns a `Result` with device,
operating system, layout engine, and agent fields.

## Requirements

- Go 1.26.1 or later (see `go.mod`).
- The module uses the Go standard library only. It does not add third-party
  modules.

## Install

```bash
go get github.com/deangrant/user-agent/useragent
```

## Quick start

```go
package main

import (
	"fmt"

	"github.com/deangrant/user-agent/useragent"
)

func main() {
	ua := "Mozilla/5.0 (Linux; Android 13; Pixel 7) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/120.0.0.0 Mobile Safari/537.36"

	r := useragent.Parse(ua)
	fmt.Println(r.Agent.NameVersion)             // Chrome 120.0.0.0
	fmt.Println(r.OperatingSystem.NameVersion)   // Android 13
	fmt.Println(r.Device.Class, r.Device.Name)   // Phone Pixel 7
	fmt.Println(r.IsMobile(), r.IsBot())         // true false
}
```

## Client Hints

Some browsers send a reduced User-Agent. Pass Client Hints to fill browser and
OS version fields from `Sec-CH-UA*` data.

```go
mobile := false
r := useragent.ParseWithHints(ua, useragent.ClientHints{
	Platform:        "Windows",
	PlatformVersion: "15.0.0", // Windows 11
	FullVersionList: []useragent.BrandVersion{
		{Brand: "Google Chrome", Version: "120.0.6099.109"},
		{Brand: "Chromium", Version: "120.0.6099.109"},
	},
	Mobile: &mobile,
})
```

Client Hints can update the OS name and version in the same OS family. Example:
Windows 10 from the User-Agent can become Windows 11 from
`PlatformVersion`. If the Client Hints platform is in a different OS family
than the User-Agent OS, the package keeps the User-Agent OS. It also sets
`Result.ClientHintsMismatch` to true.

The package parses header values as an RFC 8941 structured-fields subset. It
supports the Client Hints shapes in use here (lists with parameters, strings,
tokens, and booleans). It is not a full structured-fields library.

### ClientHints fields

| Field | Source header |
| --- | --- |
| `Brands` | `Sec-CH-UA` |
| `FullVersionList` | `Sec-CH-UA-Full-Version-List` |
| `FullVersion` | `Sec-CH-UA-Full-Version` (deprecated) |
| `Platform` | `Sec-CH-UA-Platform` |
| `PlatformVersion` | `Sec-CH-UA-Platform-Version` |
| `Mobile` | `Sec-CH-UA-Mobile` |
| `Model` | `Sec-CH-UA-Model` |
| `Arch` | `Sec-CH-UA-Arch` |
| `Bitness` | `Sec-CH-UA-Bitness` |
| `FormFactors` | `Sec-CH-UA-Form-Factors` |
| `WoW64` | `Sec-CH-UA-WoW64` |

Use `ClientHints.Empty` to test for an empty value set. Use
`ClientHintsFromHeader` or `ClientHintsFromMap` to build `ClientHints` from
HTTP headers or a name-to-value map. Map keys are case-insensitive.

### From an HTTP request

```go
r := useragent.ParseRequest(req)
// or
r := useragent.ParseHeaders(req.Header)
```

To request high-entropy hints from the browser, send an `Accept-CH` response
header. Example:

```
Accept-CH: Sec-CH-UA, Sec-CH-UA-Arch, Sec-CH-UA-Bitness, Sec-CH-UA-Form-Factors, Sec-CH-UA-Full-Version, Sec-CH-UA-Full-Version-List, Sec-CH-UA-Mobile, Sec-CH-UA-Model, Sec-CH-UA-Platform, Sec-CH-UA-Platform-Version, Sec-CH-UA-WoW64
```

## Result

| Section | Fields |
| --- | --- |
| Device | Class, Name, Brand, CPU |
| OperatingSystem | Class, Name, Version, VersionBuild, NameVersion |
| LayoutEngine | Class, Name, Version, VersionMajor, NameVersion, NameVersionMajor |
| Agent | Class, Name, Version, VersionMajor, NameVersion, NameVersionMajor |
| (top-level) | UserAgent, AgentSecurity, ClientHintsMismatch |

`LayoutEngine.Name` is the rendering engine (`Blink`, `WebKit`, `Gecko`, and
related names) when the User-Agent has that data. For bots, apps, and HTTP
libraries without a layout engine token, `Name` can be empty. Use
`LayoutEngine.Class` and `Agent` in that case.

Class fields use typed `*Unknown` constants when unset. Name and version
string fields stay empty (`""`) when unset.

Helpers:

- `IsMobile`
- `IsTablet`
- `IsPhone`
- `IsDesktop`
- `IsComputer` (same as `IsDesktop`)
- `IsBot`

`IsBot` is true for crawlers and other non-browser clients. Examples: HTTP
libraries (curl, okhttp), cloud apps, hacker tools, and test clients.

The library parse functions do not return an error. A partial `Result` is
normal for sparse input. The embedded pattern catalog loads on first use of
an analyzer (`NewAnalyzer` or the first package-level `Parse*`). Invalid
embedded catalog data can panic at that time. Released builds ship valid
data.

## Analyzer

Create an `Analyzer` when you parse many requests:

```go
a := useragent.NewAnalyzer()
r := a.Parse(ua)
```

`Analyzer` methods:

- `Parse`
- `ParseWithHints`
- `ParseHeaders`
- `ParseRequest`

A nil `Analyzer` receiver on `ParseWithHints` uses `NewAnalyzer()`.
Package-level `Parse*` functions use one shared default analyzer.

## CLI

Build and run the command-line tool:

```bash
go run ./cmd/useragent parse 'Mozilla/5.0 (Linux; Android 13; Pixel 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36'

go install ./cmd/useragent
useragent parse 'Mozilla/5.0 ...'
```

The binary uses the `parse` subcommand. By default, it prints field lines
(`[+] Device.Class: ...`). On a terminal, the `[+]` mark and key use green.
Empty values and `Unknown` class values are omitted. Use `-json` for indented
JSON. Use `-compact` for one-line JSON. Set `NO_COLOR` to disable ANSI colors.

Client Hints flags (wire values; quote strings as on the wire):

- `-sec-ch-ua`
- `-sec-ch-ua-full-version-list`
- `-sec-ch-ua-full-version`
- `-sec-ch-ua-platform`
- `-sec-ch-ua-platform-version`
- `-sec-ch-ua-mobile`
- `-sec-ch-ua-model`
- `-sec-ch-ua-arch`
- `-sec-ch-ua-bitness`
- `-sec-ch-ua-form-factors`
- `-sec-ch-ua-wow64`

Example:

```bash
useragent parse \
  -sec-ch-ua-platform '"Windows"' \
  -sec-ch-ua-platform-version '"15.0.0"' \
  -sec-ch-ua-full-version-list '"Google Chrome";v="120.0.6099.109"' \
  'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36'
```

Read HTTP headers from stdin with `-`:

```bash
printf '%s\n' \
  'User-Agent: Mozilla/5.0 (Linux; Android 13; Pixel 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36' \
  'Sec-CH-UA-Platform: "Android"' \
  'Sec-CH-UA-Platform-Version: "13.0.0"' \
  'Sec-CH-UA-Model: "Pixel 7"' \
  'Sec-CH-UA-Mobile: ?1' \
  | useragent parse -
```

Limits:

- User-Agent string: 8192 bytes maximum.
- Stdin header block: 65536 bytes maximum.

Exit codes:

| Code | Meaning |
| --- | --- |
| `0` | Success or help |
| `2` | Usage error |
| `1` | Other failure |

## Related docs

- [AGENTS.md](AGENTS.md) — agent and contributor index
- [.agents/docs/ARCHITECTURE.md](.agents/docs/ARCHITECTURE.md) — system architecture
- [SECURITY.md](SECURITY.md) — security reports

## License

MIT
