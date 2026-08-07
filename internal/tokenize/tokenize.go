// Package tokenize splits a User-Agent string into products and comments.
package tokenize

import (
	"strings"
	"unicode"
)

// Product is a name/version token such as "Chrome/120.0.0.0".
type Product struct {
	Name    string
	Version string
}

// Tokens holds parsed products and parenthetical comment fragments.
type Tokens struct {
	Products []Product
	Comments []string
	Raw      string
}

// Parse tokenizes a User-Agent string.
func Parse(ua string) Tokens {
	t := Tokens{Raw: ua}
	if ua == "" {
		return t
	}
	i := 0
	for i < len(ua) {
		for i < len(ua) && isSep(ua[i]) {
			i++
		}
		if i >= len(ua) {
			break
		}
		if ua[i] == '(' {
			end := strings.IndexByte(ua[i:], ')')
			if end < 0 {
				comment := strings.TrimSpace(ua[i+1:])
				if comment != "" {
					t.Comments = append(t.Comments, splitComment(comment)...)
				}
				break
			}
			comment := strings.TrimSpace(ua[i+1 : i+end])
			if comment != "" {
				t.Comments = append(t.Comments, splitComment(comment)...)
			}
			i += end + 1
			continue
		}
		start := i
		for i < len(ua) && !isSep(ua[i]) && ua[i] != '(' {
			i++
		}
		token := ua[start:i]
		if token == "" {
			continue
		}
		name, version := splitProduct(token)
		if name != "" {
			t.Products = append(t.Products, Product{
				Name:    name,
				Version: version,
			})
		}
	}
	return t
}

// FindProduct returns the first product with the given name
// (case-insensitive).
func (t Tokens) FindProduct(name string) (Product, bool) {
	for _, p := range t.Products {
		if strings.EqualFold(p.Name, name) {
			return p, true
		}
	}
	return Product{}, false
}

// HasCommentFragment reports whether any comment contains substr
// (case-insensitive).
func (t Tokens) HasCommentFragment(substr string) bool {
	lower := strings.ToLower(substr)
	for _, c := range t.Comments {
		if strings.Contains(strings.ToLower(c), lower) {
			return true
		}
	}
	return false
}

// CommentContains reports whether any comment equals value
// (case-insensitive) or the joined comments contain it.
func (t Tokens) CommentContains(substr string) bool {
	return t.HasCommentFragment(substr)
}

func splitProduct(token string) (name, version string) {
	token = strings.TrimSpace(token)
	if token == "" {
		return "", ""
	}
	if i := strings.IndexByte(token, '/'); i >= 0 {
		return token[:i], token[i+1:]
	}
	return token, ""
}

func splitComment(comment string) []string {
	parts := strings.Split(comment, ";")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func isSep(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

// LowerContains is a small helper for detectors.
func LowerContains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// IsDigitRun reports whether s is non-empty and all digits/dots.
func IsDigitRun(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r != '.' && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
