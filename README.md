# useragent

Go library that analyzes HTTP `User-Agent` strings and optional User-Agent
Client Hints (`Sec-CH-UA*`) into structured details about the device,
operating system, layout engine, and agent.

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
	fmt.Println(r.IsMobile(), r.IsBot())
}
```

## Client Hints

Modern Chromium browsers often send a reduced User-Agent. Pass Client Hints
for accurate browser and OS versions:

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

From an HTTP request:

```go
r := useragent.ParseRequest(req)
// or
r := useragent.ParseHeaders(req.Header)
```

To receive high-entropy hints, send an `Accept-CH` response header, for
example:

```
Accept-CH: Sec-CH-UA, Sec-CH-UA-Arch, Sec-CH-UA-Bitness, Sec-CH-UA-Form-Factors, Sec-CH-UA-Full-Version-List, Sec-CH-UA-Mobile, Sec-CH-UA-Model, Sec-CH-UA-Platform, Sec-CH-UA-Platform-Version
```

## Result fields

| Section | Fields |
| --- | --- |
| Device | Class, Name, Brand, CPU |
| OperatingSystem | Class, Name, Version, VersionBuild, NameVersion |
| LayoutEngine | Class, Name, Version, VersionMajor, NameVersion, NameVersionMajor |
| Agent | Class, Name, Version, VersionMajor, NameVersion, NameVersionMajor |

Convenience helpers: `IsMobile`, `IsTablet`, `IsPhone`, `IsDesktop`,
`IsComputer`, `IsBot`.

Unknown values use empty strings and `*Unknown` class constants. Parsing
never returns an error; partial results are normal for sparse inputs.

## Analyzer

Reuse an `Analyzer` when parsing many requests:

```go
a := useragent.NewAnalyzer()
r := a.Parse(ua)
```

Package-level `Parse*` functions use a shared default analyzer.

## CLI

Build and run the command-line tool:

```bash
go run ./cmd/useragent parse 'Mozilla/5.0 (Linux; Android 13; Pixel 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36'

go install ./cmd/useragent
useragent parse 'Mozilla/5.0 ...'
```

Pretty JSON is printed by default. Use `-compact` for one line.

With Client Hints flags:

```bash
useragent parse \
  -sec-ch-ua-platform '"Windows"' \
  -sec-ch-ua-platform-version '"15.0.0"' \
  -sec-ch-ua-full-version-list '"Google Chrome";v="120.0.6099.109"' \
  'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36'
```

Or pipe HTTP headers on stdin (`-`):

```bash
printf '%s\n' \
  'User-Agent: Mozilla/5.0 (Linux; Android 13; Pixel 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36' \
  'Sec-CH-UA-Platform: "Android"' \
  'Sec-CH-UA-Platform-Version: "13.0.0"' \
  'Sec-CH-UA-Model: "Pixel 7"' \
  'Sec-CH-UA-Mobile: ?1' \
  | useragent parse -
```

## License

MIT
