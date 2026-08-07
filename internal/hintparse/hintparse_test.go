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
	brand, ok := hintparse.SignificantBrand(list)
	if !ok || brand.Brand != "Google Chrome" || brand.Version != "120" {
		t.Fatalf("significant = %#v ok=%v", brand, ok)
	}
}

func TestParseBooleanAndUnquote(t *testing.T) {
	if b := hintparse.ParseBoolean("?1"); b == nil || !*b {
		t.Fatalf("?1 => %v", b)
	}
	if b := hintparse.ParseBoolean("?0"); b == nil || *b {
		t.Fatalf("?0 => %v", b)
	}
	if got := hintparse.Unquote(`"Windows"`); got != "Windows" {
		t.Fatalf("Unquote = %q", got)
	}
}

func TestParseHeaders(t *testing.T) {
	h := hintparse.Parse(map[string]string{
		"sec-ch-ua-platform":         `"Android"`,
		"sec-ch-ua-platform-version": `"13.0.0"`,
		"sec-ch-ua-mobile":           "?1",
		"sec-ch-ua-model":            `"Pixel 7"`,
	})
	if h.Platform != "Android" || h.Model != "Pixel 7" {
		t.Fatalf("hints = %#v", h)
	}
	if h.Mobile == nil || !*h.Mobile {
		t.Fatalf("mobile = %v", h.Mobile)
	}
}
