// Package version provides helpers for User-Agent version strings.
package version

import (
	"strings"
	"unicode"
)

// Major returns the first numeric component of a version string.
func Major(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	var b strings.Builder
	for _, r := range v {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
			continue
		}
		break
	}
	return b.String()
}

// NameVersion joins name and version with a space when both exist.
func NameVersion(name, ver string) string {
	name = strings.TrimSpace(name)
	ver = strings.TrimSpace(ver)
	switch {
	case name == "" && ver == "":
		return ""
	case ver == "":
		return name
	case name == "":
		return ver
	default:
		return name + " " + ver
	}
}

// NameVersionMajor joins name and major version.
func NameVersionMajor(name, ver string) string {
	return NameVersion(name, Major(ver))
}

// FirstNumber extracts a leading dotted number from s.
func FirstNumber(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	var b strings.Builder
	started := false
	for _, r := range s {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
			started = true
			continue
		}
		if started && r == '.' {
			b.WriteRune(r)
			continue
		}
		if started {
			break
		}
	}
	out := b.String()
	out = strings.Trim(out, ".")
	return out
}

// CompareMajor compares major version numbers as integers.
// Returns -1, 0, 1. Non-numeric majors compare as 0.
func CompareMajor(a, b string) int {
	ai := atoi(Major(a))
	bi := atoi(Major(b))
	switch {
	case ai < bi:
		return -1
	case ai > bi:
		return 1
	default:
		return 0
	}
}

func atoi(s string) int {
	n := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			break
		}
		n = n*10 + int(r-'0')
	}
	return n
}

// NormalizeSeparators replaces '_' with '.' in version-like strings.
func NormalizeSeparators(v string) string {
	return strings.ReplaceAll(v, "_", ".")
}
