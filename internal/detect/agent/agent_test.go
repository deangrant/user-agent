package agent

import (
	"testing"

	"github.com/deangrant/user-agent/internal/detect"
	"github.com/deangrant/user-agent/internal/tokenize"
)

func TestDetectChromeDesktop(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.AgentName != "Chrome" {
		t.Fatalf("agent = %q, want Chrome", state.AgentName)
	}
	if state.AgentVersion != "120.0.0.0" {
		t.Fatalf("version = %q, want 120.0.0.0", state.AgentVersion)
	}
	if state.AgentClass != "Browser" {
		t.Fatalf("class = %q, want Browser", state.AgentClass)
	}
}

func TestDetectSafari(t *testing.T) {
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) " +
		"AppleWebKit/605.1.15 (KHTML, like Gecko) " +
		"Version/18.0 Safari/605.1.15"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.AgentName != "Safari" {
		t.Fatalf("agent = %q, want Safari", state.AgentName)
	}
	if state.AgentVersion != "18.0" {
		t.Fatalf("version = %q, want 18.0", state.AgentVersion)
	}
}

func TestDetectEdgeBeforeChrome(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.AgentName != "Edge" {
		t.Fatalf("agent = %q, want Edge", state.AgentName)
	}
}

func TestDetectIE11(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) " +
		"like Gecko"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.AgentName != "IE" {
		t.Fatalf("agent = %q, want IE", state.AgentName)
	}
	if state.AgentVersion != "11.0" {
		t.Fatalf("version = %q, want 11.0", state.AgentVersion)
	}
}

func TestDetectCriOS(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 17_2 like Mac OS X) " +
		"AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/120.0.6099.119 " +
		"Mobile/15E148 Safari/604.1"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.AgentName != "Chrome" {
		t.Fatalf("agent = %q, want Chrome", state.AgentName)
	}
}

func TestDetectMobileAppNoOp(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) " +
		"AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 " +
		"Instagram 269.0.0.18.75"
	state := &detect.State{
		UA:         ua,
		Tokens:     tokenize.Parse(ua),
		AppMatched: true,
		AgentName:  "Instagram",
		AgentClass: "Mobile App",
	}
	New().Detect(state)
	if state.AgentName != "Instagram" {
		t.Fatalf("agent = %q, want Instagram unchanged", state.AgentName)
	}
	if state.AgentClass != "Mobile App" {
		t.Fatalf("class = %q, want Mobile App", state.AgentClass)
	}
}

func TestDetectWebviewKeepsClassFillsChrome(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 13; Pixel 7; wv) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Version/4.0 Chrome/120.0.0.0 Mobile Safari/537.36"
	state := &detect.State{
		UA:         ua,
		Tokens:     tokenize.Parse(ua),
		AppMatched: true,
		AgentName:  "Android WebView",
		AgentClass: "Browser Webview",
	}
	New().Detect(state)
	if state.AgentClass != "Browser Webview" {
		t.Fatalf("class = %q, want Browser Webview", state.AgentClass)
	}
	if state.AgentName != "Chrome" {
		t.Fatalf("agent = %q, want Chrome", state.AgentName)
	}
}

func TestDetectSecurityStrong(t *testing.T) {
	ua := "Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9.0.4) " +
		"Gecko/20100101 Firefox/3.0.4"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.AgentSecurity != "Strong security" {
		t.Fatalf("security = %q, want Strong security", state.AgentSecurity)
	}
}

func TestDetectSecurityWeak(t *testing.T) {
	ua := "Mozilla/5.0 (X11; I; Linux i686; en-US; rv:1.8.1) " +
		"Gecko/20061010 Firefox/2.0"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.AgentSecurity != "Weak security" {
		t.Fatalf("security = %q, want Weak security", state.AgentSecurity)
	}
}

func TestDetectSecurityNone(t *testing.T) {
	ua := "Mozilla/5.0 (Windows; N; Windows NT 5.1; en-US; rv:1.7.5) " +
		"Gecko/20041107 Firefox/1.0"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.AgentSecurity != "No security" {
		t.Fatalf("security = %q, want No security", state.AgentSecurity)
	}
}

func TestDetectSecurityUnsetWithoutToken(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.AgentSecurity != "" {
		t.Fatalf("security = %q, want empty before finalize",
			state.AgentSecurity)
	}
}

func TestDetectSecurityDoesNotOverwrite(t *testing.T) {
	ua := "Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9.0.4) " +
		"Gecko/20100101 Firefox/3.0.4"
	state := &detect.State{
		UA:            ua,
		Tokens:        tokenize.Parse(ua),
		AgentSecurity: "Hacker",
	}
	New().Detect(state)
	if state.AgentSecurity != "Hacker" {
		t.Fatalf("security = %q, want Hacker preserved", state.AgentSecurity)
	}
}
