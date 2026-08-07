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
			closeIdx, ok := findMatchingParen(ua, i)
			if !ok {
				rest := ua[i+1:]
				comment := strings.TrimSpace(rest)
				if comment != "" {
					t.Comments = append(t.Comments, splitComment(comment)...)
				}
				appendProducts(&t, rest)
				break
			}
			comment := strings.TrimSpace(ua[i+1 : closeIdx])
			if comment != "" {
				t.Comments = append(t.Comments, splitComment(comment)...)
			}
			i = closeIdx + 1
			continue
		}
		start := i
		for i < len(ua) && !isSep(ua[i]) && ua[i] != '(' {
			i++
		}
		appendProductToken(&t, ua[start:i])
	}
	return t
}

// findMatchingParen returns the index of the ')' that closes ua[open]
// using depth counting. ok is false when the parenthesis is unclosed.
func findMatchingParen(ua string, open int) (closeIdx int, ok bool) {
	depth := 1
	for j := open + 1; j < len(ua); j++ {
		switch ua[j] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return j, true
			}
		}
	}
	return 0, false
}

func appendProducts(t *Tokens, s string) {
	i := 0
	for i < len(s) {
		for i < len(s) && isSep(s[i]) {
			i++
		}
		if i >= len(s) {
			break
		}
		if s[i] == '(' {
			i++
			continue
		}
		start := i
		for i < len(s) && !isSep(s[i]) && s[i] != '(' {
			i++
		}
		appendProductToken(t, s[start:i])
	}
}

func appendProductToken(t *Tokens, token string) {
	name, version := splitProduct(token)
	if name != "" {
		t.Products = append(t.Products, Product{
			Name:    name,
			Version: version,
		})
	}
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
		p = strings.TrimLeft(p, "(")
		p = strings.TrimRight(p, ")")
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
