package useragent

import (
	"net/http"
	"strings"

	"github.com/deangrant/user-agent/internal/hintparse"
)

// BrandVersion is a brand and version pair from Client Hints.
type BrandVersion struct {
	Brand   string
	Version string
}

// ClientHints holds parsed User-Agent Client Hints header values.
type ClientHints struct {
	// Brands comes from Sec-CH-UA (low-entropy brand list).
	Brands []BrandVersion
	// FullVersionList comes from Sec-CH-UA-Full-Version-List.
	FullVersionList []BrandVersion
	// FullVersion is the deprecated Sec-CH-UA-Full-Version value.
	FullVersion string
	// Platform comes from Sec-CH-UA-Platform.
	Platform string
	// PlatformVersion comes from Sec-CH-UA-Platform-Version.
	PlatformVersion string
	// Mobile is set when Sec-CH-UA-Mobile is present.
	Mobile *bool
	// Model comes from Sec-CH-UA-Model.
	Model string
	// Arch comes from Sec-CH-UA-Arch.
	Arch string
	// Bitness comes from Sec-CH-UA-Bitness.
	Bitness string
	// FormFactors comes from Sec-CH-UA-Form-Factors.
	FormFactors []string
	// WoW64 is set when Sec-CH-UA-WoW64 is present.
	WoW64 *bool
}

// Empty reports whether no Client Hints fields are set.
func (h ClientHints) Empty() bool {
	return len(h.Brands) == 0 &&
		len(h.FullVersionList) == 0 &&
		h.FullVersion == "" &&
		h.Platform == "" &&
		h.PlatformVersion == "" &&
		h.Mobile == nil &&
		h.Model == "" &&
		h.Arch == "" &&
		h.Bitness == "" &&
		len(h.FormFactors) == 0 &&
		h.WoW64 == nil
}

// ClientHintsFromHeader extracts Client Hints from HTTP headers.
func ClientHintsFromHeader(h http.Header) ClientHints {
	if h == nil {
		return ClientHints{}
	}
	m := make(map[string]string, 12)
	for _, name := range hintHeaderNames {
		if v := headerGet(h, name); v != "" {
			m[name] = v
		}
	}
	return ClientHintsFromMap(m)
}

// ClientHintsFromMap extracts Client Hints from a header name→value map.
// Names are matched case-insensitively.
func ClientHintsFromMap(m map[string]string) ClientHints {
	if len(m) == 0 {
		return ClientHints{}
	}
	raw := make(map[string]string, len(m))
	for k, v := range m {
		raw[strings.ToLower(strings.TrimSpace(k))] = v
	}
	parsed := hintparse.Parse(raw)
	return clientHintsFromParsed(parsed)
}

func clientHintsFromParsed(p hintparse.Hints) ClientHints {
	out := ClientHints{
		FullVersion:     p.FullVersion,
		Platform:        p.Platform,
		PlatformVersion: p.PlatformVersion,
		Model:           p.Model,
		Arch:            p.Arch,
		Bitness:         p.Bitness,
		FormFactors:     append([]string(nil), p.FormFactors...),
		Mobile:          p.Mobile,
		WoW64:           p.WoW64,
	}
	for _, b := range p.Brands {
		out.Brands = append(out.Brands, BrandVersion{
			Brand:   b.Brand,
			Version: b.Version,
		})
	}
	for _, b := range p.FullVersionList {
		out.FullVersionList = append(out.FullVersionList, BrandVersion{
			Brand:   b.Brand,
			Version: b.Version,
		})
	}
	return out
}

func hintsToInternal(h ClientHints) hintparse.Hints {
	out := hintparse.Hints{
		FullVersion:     h.FullVersion,
		Platform:        h.Platform,
		PlatformVersion: h.PlatformVersion,
		Model:           h.Model,
		Arch:            h.Arch,
		Bitness:         h.Bitness,
		FormFactors:     append([]string(nil), h.FormFactors...),
		Mobile:          h.Mobile,
		WoW64:           h.WoW64,
	}
	for _, b := range h.Brands {
		out.Brands = append(out.Brands, hintparse.BrandVersion{
			Brand:   b.Brand,
			Version: b.Version,
		})
	}
	for _, b := range h.FullVersionList {
		bv := hintparse.BrandVersion{
			Brand:   b.Brand,
			Version: b.Version,
		}
		out.FullVersionList = append(out.FullVersionList, bv)
	}
	return out
}

var hintHeaderNames = []string{
	"Sec-CH-UA",
	"Sec-CH-UA-Arch",
	"Sec-CH-UA-Bitness",
	"Sec-CH-UA-Form-Factors",
	"Sec-CH-UA-Full-Version",
	"Sec-CH-UA-Full-Version-List",
	"Sec-CH-UA-Mobile",
	"Sec-CH-UA-Model",
	"Sec-CH-UA-Platform",
	"Sec-CH-UA-Platform-Version",
	"Sec-CH-UA-WoW64",
}

func headerGet(h http.Header, name string) string {
	return h.Get(name)
}
