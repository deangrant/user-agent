// Package hintparse parses Sec-CH-UA* Client Hints header values.
//
// Parsing uses an RFC 8941 structured-fields subset for the shapes
// Client Hints actually send (lists of items with parameters, strings,
// tokens, and booleans). It is not a full structured-fields
// implementation: dictionaries, integers, decimals, byte sequences,
// and inner lists are not supported.
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
		h.FullVersion = parseSFStringValue(v)
	}
	if v := get("sec-ch-ua-platform"); v != "" {
		h.Platform = parseSFStringValue(v)
	}
	if v := get("sec-ch-ua-platform-version"); v != "" {
		h.PlatformVersion = parseSFStringValue(v)
	}
	if v := get("sec-ch-ua-model"); v != "" {
		h.Model = parseSFStringValue(v)
	}
	if v := get("sec-ch-ua-arch"); v != "" {
		h.Arch = parseSFStringValue(v)
	}
	if v := get("sec-ch-ua-bitness"); v != "" {
		h.Bitness = parseSFStringValue(v)
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

// ParseBrandList parses Sec-CH-UA / Full-Version-List as an sf-list of
// items with a string/token brand and optional v parameter.
func ParseBrandList(s string) []BrandVersion {
	items := parseSFList(s)
	if len(items) == 0 {
		return nil
	}
	var out []BrandVersion
	for _, item := range items {
		if item.Value == "" {
			continue
		}
		ver := ""
		if item.Params != nil {
			ver = item.Params["v"]
		}
		out = append(out, BrandVersion{Brand: item.Value, Version: ver})
	}
	return out
}

// ParseFormFactors parses Sec-CH-UA-Form-Factors as an sf-list of strings.
func ParseFormFactors(s string) []string {
	items := parseSFList(s)
	if len(items) == 0 {
		return nil
	}
	var out []string
	for _, item := range items {
		if item.Value != "" {
			out = append(out, item.Value)
		}
	}
	return out
}

// ParseBoolean parses Client Hints booleans.
// Prefers RFC 8941 sf-boolean (?0 / ?1); also accepts true/false/0/1
// for backward compatibility.
func ParseBoolean(s string) *bool {
	if b := parseSFBoolean(s); b != nil {
		return b
	}
	s = strings.TrimSpace(Unquote(s))
	switch strings.ToLower(s) {
	case "1", "true":
		v := true
		return &v
	case "0", "false":
		v := false
		return &v
	default:
		if b, err := strconv.ParseBool(s); err == nil {
			return &b
		}
		return nil
	}
}

// Unquote strips surrounding double quotes using sf-string rules
// (\" and \\), or returns a trimmed bare token.
func Unquote(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if s[0] == '"' {
		if v, _, ok := parseSFString(s); ok {
			return v
		}
		// Fallback for malformed quotes: strip outer quotes only.
		if len(s) >= 2 && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
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
