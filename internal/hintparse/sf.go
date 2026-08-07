package hintparse

import (
	"strings"
	"unicode"
)

// sfItem is one RFC 8941 Item: a bare value plus parameters.
type sfItem struct {
	Value  string
	Params map[string]string
}

// parseSFList parses an sf-list into items (CH subset: no inner lists).
func parseSFList(s string) []sfItem {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	var out []sfItem
	for _, part := range splitSFList(s) {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		item, ok := parseSFItem(part)
		if !ok {
			continue
		}
		out = append(out, item)
	}
	return out
}

// parseSFItem parses a single sf-item (string or token + parameters).
func parseSFItem(s string) (sfItem, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return sfItem{}, false
	}
	value, rest, ok := parseSFBareItem(s)
	if !ok {
		return sfItem{}, false
	}
	params := parseSFParameters(rest)
	return sfItem{Value: value, Params: params}, true
}

// parseSFStringValue parses a header that is a single string/token item
// (optionally with ignored parameters) and returns its string value.
func parseSFStringValue(s string) string {
	item, ok := parseSFItem(s)
	if !ok {
		return ""
	}
	return item.Value
}

// parseSFBoolean parses an sf-boolean (?0 / ?1), optionally as a lone item.
func parseSFBoolean(s string) *bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	// Allow a bare boolean or an item whose bare value is the boolean.
	if b := parseSFBooleanBare(s); b != nil {
		return b
	}
	item, ok := parseSFItem(s)
	if !ok {
		return nil
	}
	return parseSFBooleanBare(item.Value)
}

func parseSFBooleanBare(s string) *bool {
	s = strings.TrimSpace(s)
	switch s {
	case "?1":
		v := true
		return &v
	case "?0":
		v := false
		return &v
	default:
		return nil
	}
}

// parseSFBareItem reads a quoted sf-string or bare token from the start of s.
func parseSFBareItem(s string) (value, rest string, ok bool) {
	s = strings.TrimLeftFunc(s, unicode.IsSpace)
	if s == "" {
		return "", "", false
	}
	if s[0] == '"' {
		val, n, ok := parseSFString(s)
		if !ok {
			return "", "", false
		}
		return val, strings.TrimSpace(s[n:]), true
	}
	tok, n, ok := parseSFToken(s)
	if !ok {
		return "", "", false
	}
	return tok, strings.TrimSpace(s[n:]), true
}

// parseSFString parses a quoted string per RFC 8941 §3.3.3.
// Returns the unescaped value and bytes consumed (including quotes).
func parseSFString(s string) (value string, n int, ok bool) {
	if len(s) < 2 || s[0] != '"' {
		return "", 0, false
	}
	var b strings.Builder
	i := 1
	for i < len(s) {
		c := s[i]
		switch c {
		case '\\':
			if i+1 >= len(s) {
				return "", 0, false
			}
			esc := s[i+1]
			if esc != '"' && esc != '\\' {
				return "", 0, false
			}
			b.WriteByte(esc)
			i += 2
		case '"':
			return b.String(), i + 1, true
		default:
			if c <= 0x1f || c == 0x7f {
				return "", 0, false
			}
			b.WriteByte(c)
			i++
		}
	}
	return "", 0, false
}

// parseSFToken parses an RFC 8941 token (tchar-like starting rules).
func parseSFToken(s string) (value string, n int, ok bool) {
	if s == "" {
		return "", 0, false
	}
	c0 := s[0]
	if c0 != '*' && !isAlpha(c0) {
		return "", 0, false
	}
	i := 1
	for i < len(s) {
		c := s[i]
		if !isTokenChar(c) {
			break
		}
		i++
	}
	return s[:i], i, true
}

func parseSFParameters(s string) map[string]string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	params := make(map[string]string)
	for s != "" {
		if s[0] != ';' {
			break
		}
		s = strings.TrimSpace(s[1:])
		if s == "" {
			break
		}
		key, n, ok := parseSFParamKey(s)
		if !ok {
			break
		}
		s = strings.TrimSpace(s[n:])
		val := ""
		if s != "" && s[0] == '=' {
			s = strings.TrimSpace(s[1:])
			v, rest, ok := parseSFBareItem(s)
			if !ok {
				break
			}
			val = v
			s = rest
		}
		params[strings.ToLower(key)] = val
		s = strings.TrimSpace(s)
	}
	if len(params) == 0 {
		return nil
	}
	return params
}

func parseSFParamKey(s string) (key string, n int, ok bool) {
	if s == "" {
		return "", 0, false
	}
	c0 := s[0]
	if c0 != '*' && !isLowerAlpha(c0) {
		return "", 0, false
	}
	i := 1
	for i < len(s) {
		c := s[i]
		if !isLowerAlpha(c) && !isDigit(c) && c != '_' && c != '-' &&
			c != '.' && c != '*' {
			break
		}
		i++
	}
	return s[:i], i, true
}

// splitSFList splits an sf-list on commas outside of quotes.
func splitSFList(s string) []string {
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
			b.WriteByte(c)
			if inQuotes && i+1 < len(s) {
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

func isAlpha(c byte) bool {
	return isLowerAlpha(c) || (c >= 'A' && c <= 'Z')
}

func isLowerAlpha(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isTokenChar(c byte) bool {
	// RFC 8941 token: ALPHA / DIGIT / ":" / "/" / and other tchars
	// used in practice for CH (e.g. version tokens are quoted strings).
	switch c {
	case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^',
		'_', '`', '|', '~', ':', '/':
		return true
	}
	return isAlpha(c) || isDigit(c)
}
