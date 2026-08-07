package tokenize_test

import (
	"testing"

	"github.com/deangrant/user-agent/internal/tokenize"
)

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
	if !tok.HasCommentFragment("Windows NT 10.0") {
		t.Fatalf("comments = %v, want Windows NT", tok.Comments)
	}
}

func TestParseEmpty(t *testing.T) {
	tok := tokenize.Parse("")
	if len(tok.Products) != 0 || len(tok.Comments) != 0 {
		t.Fatalf("got %+v", tok)
	}
}
