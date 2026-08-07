package engine

import (
	"testing"

	"github.com/deangrant/user-agent/internal/detect"
	"github.com/deangrant/user-agent/internal/tokenize"
)

func TestDetectBlinkChrome(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	state := &detect.State{
		UA:         ua,
		Tokens:     tokenize.Parse(ua),
		AgentName:  "Chrome",
		AgentClass: "Browser",
	}
	New().Detect(state)
	if state.EngineName != "Blink" {
		t.Fatalf("engine = %q, want Blink", state.EngineName)
	}
	if state.EngineVersion != "120.0.0.0" {
		t.Fatalf("version = %q, want 120.0.0.0", state.EngineVersion)
	}
}

func TestDetectWebKitSafari(t *testing.T) {
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) " +
		"AppleWebKit/605.1.15 (KHTML, like Gecko) " +
		"Version/18.0 Safari/605.1.15"
	state := &detect.State{
		UA:         ua,
		Tokens:     tokenize.Parse(ua),
		AgentName:  "Safari",
		AgentClass: "Browser",
	}
	New().Detect(state)
	if state.EngineName != "WebKit" {
		t.Fatalf("engine = %q, want WebKit", state.EngineName)
	}
	if state.EngineVersion != "605.1.15" {
		t.Fatalf("version = %q, want 605.1.15", state.EngineVersion)
	}
}

func TestDetectGeckoFirefox(t *testing.T) {
	ua := "Mozilla/5.0 (X11; Linux x86_64; rv:121.0) " +
		"Gecko/20100101 Firefox/121.0"
	state := &detect.State{
		UA:         ua,
		Tokens:     tokenize.Parse(ua),
		AgentName:  "Firefox",
		AgentClass: "Browser",
	}
	New().Detect(state)
	if state.EngineName != "Gecko" {
		t.Fatalf("engine = %q, want Gecko", state.EngineName)
	}
	if state.EngineVersion != "121.0" {
		t.Fatalf("version = %q, want 121.0", state.EngineVersion)
	}
}

func TestDetectLegacyEdgeHTML(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36 " +
		"Edge/15.15063"
	state := &detect.State{
		UA:         ua,
		Tokens:     tokenize.Parse(ua),
		AgentName:  "Edge",
		AgentClass: "Browser",
	}
	New().Detect(state)
	if state.EngineName != "EdgeHTML" {
		t.Fatalf("engine = %q, want EdgeHTML", state.EngineName)
	}
}

func TestDetectChromiumEdgeBlink(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0"
	state := &detect.State{
		UA:         ua,
		Tokens:     tokenize.Parse(ua),
		AgentName:  "Edge",
		AgentClass: "Browser",
	}
	New().Detect(state)
	if state.EngineName != "Blink" {
		t.Fatalf("engine = %q, want Blink", state.EngineName)
	}
}

func TestDetectPresetEngineNameNoOverwrite(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	state := &detect.State{
		UA:          ua,
		Tokens:      tokenize.Parse(ua),
		EngineName:  "Custom",
		EngineClass: "Special",
	}
	New().Detect(state)
	if state.EngineName != "Custom" {
		t.Fatalf("engine = %q, want Custom", state.EngineName)
	}
}

func TestDetectMobileAppEngineClass(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) " +
		"AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 " +
		"Instagram 269.0.0.18.75"
	state := &detect.State{
		UA:         ua,
		Tokens:     tokenize.Parse(ua),
		AgentName:  "Instagram",
		AgentClass: "Mobile App",
	}
	New().Detect(state)
	if state.EngineClass != "Mobile App" {
		t.Fatalf("engine class = %q, want Mobile App", state.EngineClass)
	}
	if state.EngineName != "WebKit" {
		t.Fatalf("engine = %q, want WebKit", state.EngineName)
	}
}
