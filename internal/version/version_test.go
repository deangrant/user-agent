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
