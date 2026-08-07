package tokenize_test

import (
	"strings"
	"testing"

	"github.com/deangrant/user-agent/internal/tokenize"
)

func commentsContain(comments []string, substr string) bool {
	lower := strings.ToLower(substr)
	for _, c := range comments {
		if strings.Contains(strings.ToLower(c), lower) {
			return true
		}
	}
	return false
}

func TestParseProductsAndComments(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	tok := tokenize.Parse(ua)
	if len(tok.Products) < 4 {
		t.Fatalf("products = %v, want at least 4", tok.Products)
	}
	if p, ok := tok.FindProduct("Chrome"); !ok || p.Version != "120.0.0.0" {
		t.Fatalf("Chrome product = %v ok=%v", p, ok)
	}
	if !commentsContain(tok.Comments, "Windows NT 10.0") {
		t.Fatalf("comments = %v, want Windows NT", tok.Comments)
	}
}

func TestParseEmpty(t *testing.T) {
	tok := tokenize.Parse("")
	if len(tok.Products) != 0 || len(tok.Comments) != 0 {
		t.Fatalf("got %+v", tok)
	}
}

func TestParseNestedParens(t *testing.T) {
	ua := "Mozilla/5.0 ((Linux; Android 13; Pixel 7) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	tok := tokenize.Parse(ua)
	for _, c := range tok.Comments {
		if c == "(Linux" || c == "(linux" {
			t.Fatalf("comments = %v, contains garbage %q", tok.Comments, c)
		}
	}
	if !commentsContain(tok.Comments, "Linux") {
		t.Fatalf("comments = %v, want Linux", tok.Comments)
	}
	if _, ok := tok.FindProduct("Chrome"); !ok {
		t.Fatalf("products = %v, want Chrome", tok.Products)
	}
}

func TestParseUnclosedParenKeepsProducts(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 13; Pixel 7 AppleWebKit/537.36 " +
		"Chrome/120.0.0.0 Safari/537.36"
	tok := tokenize.Parse(ua)
	if p, ok := tok.FindProduct("Chrome"); !ok || p.Version != "120.0.0.0" {
		t.Fatalf("Chrome product = %v ok=%v", p, ok)
	}
	if _, ok := tok.FindProduct("Safari"); !ok {
		t.Fatalf("products = %v, want Safari", tok.Products)
	}
	if !commentsContain(tok.Comments, "Linux") {
		t.Fatalf("comments = %v, want Linux", tok.Comments)
	}
	if !commentsContain(tok.Comments, "Android") {
		t.Fatalf("comments = %v, want Android", tok.Comments)
	}
}

func TestParseSpaceSeparatedVersion(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) " +
		"AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 " +
		"Instagram 269.0.0.18.75"
	tok := tokenize.Parse(ua)
	p, ok := tok.FindProduct("Instagram")
	if !ok {
		t.Fatalf("products = %v, want Instagram", tok.Products)
	}
	if p.Version != "269.0.0.18.75" {
		t.Fatalf("Instagram version = %q, want 269.0.0.18.75", p.Version)
	}
	for _, prod := range tok.Products {
		if prod.Name == "269.0.0.18.75" {
			t.Fatalf("standalone version product present: %v", tok.Products)
		}
	}
	mobile, ok := tok.FindProduct("Mobile")
	if !ok || mobile.Version != "15E148" {
		t.Fatalf("Mobile slash product = %v ok=%v", mobile, ok)
	}
}
