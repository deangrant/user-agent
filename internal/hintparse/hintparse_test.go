package hintparse_test

import (
	"testing"

	"github.com/deangrant/user-agent/internal/hintparse"
)

func TestParseBrandList(t *testing.T) {
	list := hintparse.ParseBrandList(
		`"Not A;Brand";v="99", "Chromium";v="120", "Google Chrome";v="120"`,
	)
	if len(list) != 3 {
		t.Fatalf("len=%d list=%v", len(list), list)
	}
	if list[0].Brand != "Not A;Brand" || list[0].Version != "99" {
		t.Fatalf("grease brand = %#v", list[0])
	}
	brand, ok := hintparse.SignificantBrand(list)
	if !ok || brand.Brand != "Google Chrome" || brand.Version != "120" {
		t.Fatalf("significant = %#v ok=%v", brand, ok)
	}
}

func TestParseBrandListExtraParams(t *testing.T) {
	list := hintparse.ParseBrandList(
		`"Chromium";v="120";u=1, "Google Chrome";v="120";foo="bar"`,
	)
	if len(list) != 2 {
		t.Fatalf("len=%d list=%v", len(list), list)
	}
	if list[0].Brand != "Chromium" || list[0].Version != "120" {
		t.Fatalf("chromium = %#v", list[0])
	}
	if list[1].Brand != "Google Chrome" || list[1].Version != "120" {
		t.Fatalf("chrome = %#v", list[1])
	}
}

func TestParseBrandListEscapedQuotes(t *testing.T) {
	list := hintparse.ParseBrandList(`"Foo\"Bar";v="1"`)
	if len(list) != 1 || list[0].Brand != `Foo"Bar` || list[0].Version != "1" {
		t.Fatalf("list = %#v", list)
	}
}

func TestUnquoteEscapes(t *testing.T) {
	if got := hintparse.Unquote(`"a\"b\\c"`); got != `a"b\c` {
		t.Fatalf("Unquote escapes = %q", got)
	}
	if got := hintparse.Unquote(`"Windows"`); got != "Windows" {
		t.Fatalf("Unquote = %q", got)
	}
}

func TestParseFormFactors(t *testing.T) {
	got := hintparse.ParseFormFactors(`"Desktop", "Mobile"`)
	if len(got) != 2 || got[0] != "Desktop" || got[1] != "Mobile" {
		t.Fatalf("form factors = %#v", got)
	}
	single := hintparse.ParseFormFactors(`"XRHeadset"`)
	if len(single) != 1 || single[0] != "XRHeadset" {
		t.Fatalf("single = %#v", single)
	}
}

func TestParseBoolean(t *testing.T) {
	if b := hintparse.ParseBoolean("?1"); b == nil || !*b {
		t.Fatalf("?1 => %v", b)
	}
	if b := hintparse.ParseBoolean("?0"); b == nil || *b {
		t.Fatalf("?0 => %v", b)
	}
	if b := hintparse.ParseBoolean("true"); b == nil || !*b {
		t.Fatalf("true => %v", b)
	}
	if b := hintparse.ParseBoolean("false"); b == nil || *b {
		t.Fatalf("false => %v", b)
	}
	if b := hintparse.ParseBoolean("1"); b == nil || !*b {
		t.Fatalf("1 => %v", b)
	}
	if b := hintparse.ParseBoolean("0"); b == nil || *b {
		t.Fatalf("0 => %v", b)
	}
	if b := hintparse.ParseBoolean("nope"); b != nil {
		t.Fatalf("nope => %v", b)
	}
}

func TestParseHeaders(t *testing.T) {
	h := hintparse.Parse(map[string]string{
		"sec-ch-ua-platform":         `"Android"`,
		"sec-ch-ua-platform-version": `"13.0.0"`,
		"sec-ch-ua-mobile":           "?1",
		"sec-ch-ua-model":            `"Pixel 7"`,
		"sec-ch-ua-full-version":     `"120.0.6099.109"`,
		"sec-ch-ua-wow64":            "?0",
		"sec-ch-ua-form-factors":     `"Mobile", "Tablet"`,
	})
	if h.Platform != "Android" || h.Model != "Pixel 7" {
		t.Fatalf("hints = %#v", h)
	}
	if h.Mobile == nil || !*h.Mobile {
		t.Fatalf("mobile = %v", h.Mobile)
	}
	if h.FullVersion != "120.0.6099.109" {
		t.Fatalf("full version = %q", h.FullVersion)
	}
	if h.WoW64 == nil || *h.WoW64 {
		t.Fatalf("wow64 = %v", h.WoW64)
	}
	if len(h.FormFactors) != 2 {
		t.Fatalf("form factors = %#v", h.FormFactors)
	}
}

func TestParseStringWithIgnoredParams(t *testing.T) {
	h := hintparse.Parse(map[string]string{
		"sec-ch-ua-platform": `"Windows";u=1`,
	})
	if h.Platform != "Windows" {
		t.Fatalf("platform = %q", h.Platform)
	}
}
