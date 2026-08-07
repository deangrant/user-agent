package version_test

import (
	"testing"

	"github.com/deangrant/user-agent/internal/version"
)

func TestMajorAndNameVersion(t *testing.T) {
	if got := version.Major("120.0.4896.75"); got != "120" {
		t.Fatalf("Major = %q", got)
	}
	if got := version.NameVersion("Chrome", "120.0"); got != "Chrome 120.0" {
		t.Fatalf("NameVersion = %q", got)
	}
	if got := version.NormalizeSeparators("10_15_7"); got != "10.15.7" {
		t.Fatalf("Normalize = %q", got)
	}
}

func TestCompareMajor(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"120.0.0.0", "120.0.6099.109", 0},
		{"119.0.0.0", "120.0.0.0", -1},
		{"121.0.0.0", "120.0.0.0", 1},
		{"N/A", "120.0.0.0", 0},
		{"120.0.0.0", "N/A", 0},
		{"abc", "def", 0},
		{"", "120", 0},
		{"120", "", 0},
		{"", "", 0},
	}
	for _, tt := range tests {
		if got := version.CompareMajor(tt.a, tt.b); got != tt.want {
			t.Fatalf("CompareMajor(%q, %q) = %d, want %d",
				tt.a, tt.b, got, tt.want)
		}
	}
}
