package opsys

import (
	"testing"

	"github.com/deangrant/user-agent/internal/detect"
	"github.com/deangrant/user-agent/internal/tokenize"
)

func TestDetectWindows10(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.OSName != "Windows" || state.OSVersion != "10" {
		t.Fatalf("os = %s %s, want Windows 10", state.OSName, state.OSVersion)
	}
	if state.OSClass != "Desktop" {
		t.Fatalf("class = %q, want Desktop", state.OSClass)
	}
	if state.DeviceCPU != "x86_64" {
		t.Fatalf("cpu = %q, want x86_64", state.DeviceCPU)
	}
}

func TestDetectWindowsPhone(t *testing.T) {
	ua := "Mozilla/5.0 (Windows Phone 10.0; Android 6.0.1; " +
		"Microsoft; Lumia 950) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/52.0.2743.116 Mobile Safari/537.36 " +
		"Edge/15.15063"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.OSName != "Windows Phone" {
		t.Fatalf("os = %q, want Windows Phone", state.OSName)
	}
	if state.OSVersion != "10.0" {
		t.Fatalf("version = %q, want 10.0", state.OSVersion)
	}
	if state.OSClass != "Mobile" {
		t.Fatalf("class = %q, want Mobile", state.OSClass)
	}
}

func TestDetectAndroidBuild(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 7.0; Nexus 6 Build/NBD90Z) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/53.0.2785.124 Mobile Safari/537.36"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.OSName != "Android" || state.OSVersion != "7.0" {
		t.Fatalf("os = %s %s, want Android 7.0", state.OSName, state.OSVersion)
	}
	if state.OSVersionBuild != "NBD90Z" {
		t.Fatalf("build = %q, want NBD90Z", state.OSVersionBuild)
	}
}

func TestDetectIOS(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 17_2 like Mac OS X) " +
		"AppleWebKit/605.1.15 (KHTML, like Gecko) " +
		"Version/17.2 Mobile/15E148 Safari/604.1"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.OSName != "iOS" {
		t.Fatalf("os = %q, want iOS", state.OSName)
	}
	if state.OSVersion != "17.2" {
		t.Fatalf("version = %q, want 17.2", state.OSVersion)
	}
}

func TestDetectIPadOS(t *testing.T) {
	ua := "Mozilla/5.0 (iPad; CPU OS 17_2 like Mac OS X) " +
		"AppleWebKit/605.1.15 (KHTML, like Gecko) " +
		"Version/17.2 Mobile/15E148 Safari/604.1"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.OSName != "iPadOS" {
		t.Fatalf("os = %q, want iPadOS", state.OSName)
	}
	if state.OSVersion != "17.2" {
		t.Fatalf("version = %q, want 17.2", state.OSVersion)
	}
}

func TestDetectMacOSX(t *testing.T) {
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) " +
		"AppleWebKit/605.1.15 (KHTML, like Gecko) " +
		"Version/18.0 Safari/605.1.15"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.OSName != "Mac OS X" {
		t.Fatalf("os = %q, want Mac OS X", state.OSName)
	}
	if state.OSVersion != "10.15.7" {
		t.Fatalf("version = %q, want 10.15.7", state.OSVersion)
	}
}

func TestDetectPresetOSNameSkip(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	state := &detect.State{
		UA:      ua,
		Tokens:  tokenize.Parse(ua),
		OSName:  "Linux",
		OSClass: "Desktop",
	}
	New().Detect(state)
	if state.OSName != "Linux" {
		t.Fatalf("os = %q, want Linux unchanged", state.OSName)
	}
}
