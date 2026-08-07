package merge

import (
	"testing"

	"github.com/deangrant/user-agent/internal/detect"
	"github.com/deangrant/user-agent/internal/hintparse"
)

func boolPtr(v bool) *bool { return &v }

func TestApplyUpgradesFrozenChromeVersion(t *testing.T) {
	state := &detect.State{
		AgentName:    "Chrome",
		AgentVersion: "120.0.0.0",
		AgentClass:   "Browser",
		Hints: hintparse.Hints{
			FullVersionList: []hintparse.BrandVersion{
				{Brand: "Not A;Brand", Version: "99.0.0.0"},
				{Brand: "Chromium", Version: "120.0.6099.109"},
				{Brand: "Google Chrome", Version: "120.0.6099.109"},
			},
		},
	}
	Apply(state)
	if state.AgentName != "Chrome" {
		t.Fatalf("agent = %q, want Chrome", state.AgentName)
	}
	if state.AgentVersion != "120.0.6099.109" {
		t.Fatalf("version = %q, want 120.0.6099.109", state.AgentVersion)
	}
}

func TestApplyKeepsNumericVersionAgainstNonNumericHints(t *testing.T) {
	state := &detect.State{
		AgentName:    "Chrome",
		AgentVersion: "120.0.6099.109",
		AgentClass:   "Browser",
		Hints: hintparse.Hints{
			FullVersionList: []hintparse.BrandVersion{
				{Brand: "Google Chrome", Version: "N/A"},
				{Brand: "Chromium", Version: "N/A"},
			},
		},
	}
	Apply(state)
	if state.AgentVersion != "120.0.6099.109" {
		t.Fatalf("version = %q, want 120.0.6099.109", state.AgentVersion)
	}
}

func TestApplyKeepsSafariAgainstChromeHints(t *testing.T) {
	state := &detect.State{
		AgentName:  "Safari",
		AgentClass: "Browser",
		Hints: hintparse.Hints{
			FullVersionList: []hintparse.BrandVersion{
				{Brand: "Google Chrome", Version: "120.0.6099.109"},
				{Brand: "Chromium", Version: "120.0.6099.109"},
			},
		},
	}
	Apply(state)
	if state.AgentName != "Safari" {
		t.Fatalf("agent = %q, want Safari", state.AgentName)
	}
}

func TestApplyWindows11FromPlatformVersion(t *testing.T) {
	state := &detect.State{
		OSName:    "Windows",
		OSVersion: "10",
		OSClass:   "Desktop",
		Hints: hintparse.Hints{
			Platform:        "Windows",
			PlatformVersion: "15.0.0",
		},
	}
	Apply(state)
	if state.OSVersion != "11" {
		t.Fatalf("os version = %q, want 11", state.OSVersion)
	}
	if state.ClientHintsMismatch {
		t.Fatal("ClientHintsMismatch = true, want false")
	}
}

func TestApplyKeepsUAOSOnFamilyConflict(t *testing.T) {
	state := &detect.State{
		OSName:    "Android",
		OSVersion: "13",
		OSClass:   "Mobile",
		Hints: hintparse.Hints{
			Platform:        "Windows",
			PlatformVersion: "15.0.0",
		},
	}
	Apply(state)
	if state.OSName != "Android" {
		t.Fatalf("os name = %q, want Android", state.OSName)
	}
	if state.OSVersion != "13" {
		t.Fatalf("os version = %q, want 13", state.OSVersion)
	}
	if !state.ClientHintsMismatch {
		t.Fatal("ClientHintsMismatch = false, want true")
	}
}

func TestApplyEmptyUAOSUsesCH(t *testing.T) {
	state := &detect.State{
		Hints: hintparse.Hints{
			Platform:        "Windows",
			PlatformVersion: "15.0.0",
		},
	}
	Apply(state)
	if state.OSName != "Windows" {
		t.Fatalf("os name = %q, want Windows", state.OSName)
	}
	if state.OSVersion != "11" {
		t.Fatalf("os version = %q, want 11", state.OSVersion)
	}
	if state.ClientHintsMismatch {
		t.Fatal("ClientHintsMismatch = true, want false")
	}
}

func TestApplyModelAndBrand(t *testing.T) {
	state := &detect.State{
		Hints: hintparse.Hints{
			Model: "Pixel 7",
		},
	}
	Apply(state)
	if state.DeviceName != "Pixel 7" {
		t.Fatalf("name = %q, want Pixel 7", state.DeviceName)
	}
	if state.DeviceBrand != "Google" {
		t.Fatalf("brand = %q, want Google", state.DeviceBrand)
	}
}

func TestApplyFormFactorsTablet(t *testing.T) {
	state := &detect.State{
		DeviceClass: "Phone",
		Hints: hintparse.Hints{
			FormFactors: []string{"Tablet"},
		},
	}
	Apply(state)
	if state.DeviceClass != "Tablet" {
		t.Fatalf("class = %q, want Tablet", state.DeviceClass)
	}
}

func TestApplyFormFactorsWatchOverMobile(t *testing.T) {
	state := &detect.State{
		DeviceClass: "Phone",
		Hints: hintparse.Hints{
			FormFactors: []string{"Mobile", "Watch"},
		},
	}
	Apply(state)
	if state.DeviceClass != "Watch" {
		t.Fatalf("class = %q, want Watch", state.DeviceClass)
	}
}

func TestApplyMobileOverridesTablet(t *testing.T) {
	state := &detect.State{
		DeviceClass: "Tablet",
		Hints: hintparse.Hints{
			Mobile: boolPtr(true),
		},
	}
	Apply(state)
	if state.DeviceClass != "Phone" {
		t.Fatalf("class = %q, want Phone", state.DeviceClass)
	}
}

func TestApplyCPUHints(t *testing.T) {
	state := &detect.State{
		Hints: hintparse.Hints{
			Arch:    "x86",
			Bitness: "64",
		},
	}
	Apply(state)
	if state.DeviceCPU != "x86 64-bit" {
		t.Fatalf("cpu = %q, want x86 64-bit", state.DeviceCPU)
	}
}

func TestApplyFinalizeUnknownAndDerived(t *testing.T) {
	state := &detect.State{
		AgentName:    "Chrome",
		AgentVersion: "120.0.6099.109",
	}
	Apply(state)
	if state.DeviceClass != "Unknown" {
		t.Fatalf("device class = %q, want Unknown", state.DeviceClass)
	}
	if state.OSClass != "Unknown" {
		t.Fatalf("os class = %q, want Unknown", state.OSClass)
	}
	if state.EngineClass != "Unknown" {
		t.Fatalf("engine class = %q, want Unknown", state.EngineClass)
	}
	if state.AgentClass != "Unknown" {
		t.Fatalf("agent class = %q, want Unknown", state.AgentClass)
	}
	if state.AgentVersionMajor != "120" {
		t.Fatalf("major = %q, want 120", state.AgentVersionMajor)
	}
	if state.AgentNameVersion == "" {
		t.Fatal("AgentNameVersion empty")
	}
	if state.AgentNameVersionMajor == "" {
		t.Fatal("AgentNameVersionMajor empty")
	}
}

func TestShouldOverrideAgent(t *testing.T) {
	tests := []struct {
		existing, fromCH string
		want             bool
	}{
		{"Chrome", "Chrome", false},
		{"Chrome", "Google Chrome", true},
		{"Safari", "Chrome", false},
		{"Firefox", "Chrome", false},
		{"", "Chrome", true},
		{"Unknown", "Edge", true},
	}
	for _, tt := range tests {
		got := shouldOverrideAgent(tt.existing, tt.fromCH)
		if got != tt.want {
			t.Fatalf("shouldOverrideAgent(%q, %q) = %v, want %v",
				tt.existing, tt.fromCH, got, tt.want)
		}
	}
}

func TestApplyPreservesIPadOSAgainstIOSPlatform(t *testing.T) {
	state := &detect.State{
		OSName:    "iPadOS",
		OSVersion: "17.0",
		OSClass:   "Mobile",
		Hints: hintparse.Hints{
			Platform:        "iOS",
			PlatformVersion: "17.2.0",
		},
	}
	Apply(state)
	if state.OSName != "iPadOS" {
		t.Fatalf("os = %q, want iPadOS", state.OSName)
	}
	if state.OSVersion != "17.2.0" {
		t.Fatalf("version = %q, want 17.2.0", state.OSVersion)
	}
}

func TestApplyPlatformIPadOS(t *testing.T) {
	state := &detect.State{
		Hints: hintparse.Hints{
			Platform:        "iPadOS",
			PlatformVersion: "17.2.0",
		},
	}
	Apply(state)
	if state.OSName != "iPadOS" {
		t.Fatalf("os = %q, want iPadOS", state.OSName)
	}
}

func TestApplyIOSPlatformWithIPadModel(t *testing.T) {
	state := &detect.State{
		Hints: hintparse.Hints{
			Platform:        "iOS",
			PlatformVersion: "17.2.0",
			Model:           "iPad",
		},
	}
	Apply(state)
	if state.OSName != "iPadOS" {
		t.Fatalf("os = %q, want iPadOS", state.OSName)
	}
	if state.DeviceName != "iPad" {
		t.Fatalf("device = %q, want iPad", state.DeviceName)
	}
}
