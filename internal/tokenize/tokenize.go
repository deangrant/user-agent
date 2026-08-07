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
				scanProducts(&t, rest)
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
		i = attachOrAppendProduct(&t, ua, start, i)
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

func scanProducts(t *Tokens, s string) {
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
		i = attachOrAppendProduct(t, s, start, i)
	}
}

// attachOrAppendProduct adds ua[start:end] as a product. When the product
// has no slash version, a following dotted numeric token is merged in.
// Returns the index after any consumed version token.
func attachOrAppendProduct(t *Tokens, ua string, start, end int) int {
	name, version := splitProduct(ua[start:end])
	if name == "" {
		return end
	}
	next := end
	if version == "" {
		if ver, after, ok := peekDottedVersion(ua, end); ok {
			version = ver
			next = after
		}
	}
	t.Products = append(t.Products, Product{
		Name:    name,
		Version: version,
	})
	return next
}

// peekDottedVersion looks past whitespace for a digit/dot token that
// contains at least one '.' (e.g. "269.0.0.18.75").
func peekDottedVersion(ua string, i int) (ver string, after int, ok bool) {
	j := i
	for j < len(ua) && isSep(ua[j]) {
		j++
	}
	if j >= len(ua) || ua[j] == '(' {
		return "", i, false
	}
	start := j
	for j < len(ua) && !isSep(ua[j]) && ua[j] != '(' {
		j++
	}
	tok := ua[start:j]
	if !isDottedVersion(tok) {
		return "", i, false
	}
	return tok, j, true
}

func isDottedVersion(s string) bool {
	return isDigitRun(s) && strings.Contains(s, ".")
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

func isDigitRun(s string) bool {
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
