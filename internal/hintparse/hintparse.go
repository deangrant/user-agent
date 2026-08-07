// Package hintparse parses Sec-CH-UA* Client Hints header values.
package hintparse

import (
	"strconv"
	"strings"
)

// BrandVersion is a brand/version pair from a CH brand list.
type BrandVersion struct {
	Brand   string
	Version string
}

// Hints is the structured form of Client Hints headers.
type Hints struct {
	Brands          []BrandVersion
	FullVersionList []BrandVersion
	FullVersion     string
	Platform        string
	PlatformVersion string
	Mobile          *bool
	Model           string
	Arch            string
	Bitness         string
	FormFactors     []string
	WoW64           *bool
}

// Empty reports whether no fields are set.
func (h Hints) Empty() bool {
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

// Parse builds Hints from a lower-cased header map.
func Parse(headers map[string]string) Hints {
	var h Hints
	if headers == nil {
		return h
	}
	get := func(name string) string {
		return strings.TrimSpace(headers[strings.ToLower(name)])
	}

	if v := get("sec-ch-ua"); v != "" {
		h.Brands = ParseBrandList(v)
	}
	if v := get("sec-ch-ua-full-version-list"); v != "" {
		h.FullVersionList = ParseBrandList(v)
	}
	if v := get("sec-ch-ua-full-version"); v != "" {
		h.FullVersion = Unquote(v)
	}
	if v := get("sec-ch-ua-platform"); v != "" {
		h.Platform = Unquote(v)
	}
	if v := get("sec-ch-ua-platform-version"); v != "" {
		h.PlatformVersion = Unquote(v)
	}
	if v := get("sec-ch-ua-model"); v != "" {
		h.Model = Unquote(v)
	}
	if v := get("sec-ch-ua-arch"); v != "" {
		h.Arch = Unquote(v)
	}
	if v := get("sec-ch-ua-bitness"); v != "" {
		h.Bitness = Unquote(v)
	}
	if v := get("sec-ch-ua-form-factors"); v != "" {
		h.FormFactors = ParseFormFactors(v)
	}
	if v := get("sec-ch-ua-mobile"); v != "" {
		h.Mobile = ParseBoolean(v)
	}
	if v := get("sec-ch-ua-wow64"); v != "" {
		h.WoW64 = ParseBoolean(v)
	}
	return h
}

// ParseBrandList parses Sec-CH-UA brand list header values.
func ParseBrandList(s string) []BrandVersion {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	var out []BrandVersion
	for _, part := range splitList(s) {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		brand, version := parseBrandPart(part)
		if brand == "" {
			continue
		}
		out = append(out, BrandVersion{Brand: brand, Version: version})
	}
	return out
}

func parseBrandPart(part string) (brand, version string) {
	// brand;v="version" — brand may contain ';' inside quotes.
	semi := indexSeparatorSemi(part)
	if semi < 0 {
		return Unquote(part), ""
	}
	brand = Unquote(strings.TrimSpace(part[:semi]))
	rest := strings.TrimSpace(part[semi+1:])
	restLower := strings.ToLower(rest)
	if strings.HasPrefix(restLower, "v=") {
		version = Unquote(strings.TrimSpace(rest[2:]))
	}
	return brand, version
}

// indexSeparatorSemi finds the ';' that separates brand from v=,
// ignoring semicolons inside double quotes.
func indexSeparatorSemi(s string) int {
	inQuotes := false
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"':
			inQuotes = !inQuotes
		case '\\':
			if inQuotes && i+1 < len(s) {
				i++
			}
		case ';':
			if !inQuotes {
				return i
			}
		}
	}
	return -1
}

// ParseFormFactors parses a quoted list of form factors.
func ParseFormFactors(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	// May be a single quoted value or a list.
	if !strings.Contains(s, ",") {
		v := Unquote(s)
		if v == "" {
			return nil
		}
		return []string{v}
	}
	var out []string
	for _, part := range splitList(s) {
		v := Unquote(strings.TrimSpace(part))
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

// ParseBoolean parses Client Hints booleans (?0 / ?1 / true / false).
func ParseBoolean(s string) *bool {
	s = strings.TrimSpace(Unquote(s))
	switch strings.ToLower(s) {
	case "?1", "1", "true":
		v := true
		return &v
	case "?0", "0", "false":
		v := false
		return &v
	default:
		if b, err := strconv.ParseBool(s); err == nil {
			return &b
		}
		return nil
	}
}

// Unquote strips surrounding double quotes and unescapes \" .
func Unquote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
		s = strings.ReplaceAll(s, `\"`, `"`)
	}
	return s
}

// SignificantBrand returns the most meaningful brand from a list,
// skipping greasing brands like "Not A;Brand".
func SignificantBrand(list []BrandVersion) (BrandVersion, bool) {
	var chromium BrandVersion
	var hasChromium bool
	for _, b := range list {
		if isGreaseBrand(b.Brand) {
			continue
		}
		lower := strings.ToLower(b.Brand)
		if lower == "chromium" {
			chromium = b
			hasChromium = true
			continue
		}
		return b, true
	}
	if hasChromium {
		return chromium, true
	}
	return BrandVersion{}, false
}

func isGreaseBrand(brand string) bool {
	b := strings.ToLower(strings.TrimSpace(brand))
	if b == "" {
		return true
	}
	if strings.Contains(b, "not") && strings.Contains(b, "brand") {
		return true
	}
	// Grease brands often contain punctuation noise.
	if strings.ContainsAny(b, ";:)(") {
		return true
	}
	return false
}

func splitList(s string) []string {
	var parts []string
	var b strings.Builder
	inQuotes := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '"':
			inQuotes = !inQuotes
			b.WriteByte(c)
		case '\\':
			if i+1 < len(s) {
				b.WriteByte(c)
				i++
				b.WriteByte(s[i])
			}
		case ',':
			if inQuotes {
				b.WriteByte(c)
			} else {
				parts = append(parts, b.String())
				b.Reset()
			}
		default:
			b.WriteByte(c)
		}
	}
	if b.Len() > 0 {
		parts = append(parts, b.String())
	}
	return parts
}
