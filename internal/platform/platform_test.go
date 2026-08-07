package platform_test

import (
	"testing"

	"github.com/deangrant/user-agent/internal/platform"
)

func TestResolveWindows(t *testing.T) {
	tests := []struct {
		in, ver string
	}{
		{"15.0.0", "11"},
		{"13.0.0", "11"},
		{"10.0.0", "10"},
		{"0.1.0", "7"},
		{"0.2.0", "8"},
		{"0.3.0", "8.1"},
	}
	for _, tt := range tests {
		name, ver := platform.ResolveWindows(tt.in)
		if name != "Windows" || ver != tt.ver {
			t.Fatalf("%q => %s %s, want Windows %s",
				tt.in, name, ver, tt.ver)
		}
	}
}

func TestMapWindowsNT(t *testing.T) {
	if got := platform.MapWindowsNT("6.1"); got != "7" {
		t.Fatalf("got %q", got)
	}
	if got := platform.MapWindowsNT("10.0"); got != "10" {
		t.Fatalf("got %q", got)
	}
}

func TestResolveFromCHIPadOS(t *testing.T) {
	name, ver, class := platform.ResolveFromCH("iPadOS", "17.2.0")
	if name != "iPadOS" {
		t.Fatalf("name = %q, want iPadOS", name)
	}
	if ver != "17.2.0" {
		t.Fatalf("ver = %q, want 17.2.0", ver)
	}
	if class != "Mobile" {
		t.Fatalf("class = %q, want Mobile", class)
	}
}

func TestResolveFromCHIOS(t *testing.T) {
	name, ver, class := platform.ResolveFromCH("iOS", "17.2.0")
	if name != "iOS" {
		t.Fatalf("name = %q, want iOS", name)
	}
	if ver != "17.2.0" {
		t.Fatalf("ver = %q, want 17.2.0", ver)
	}
	if class != "Mobile" {
		t.Fatalf("class = %q, want Mobile", class)
	}
}
