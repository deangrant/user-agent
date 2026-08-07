package app

import (
	"testing"

	"github.com/deangrant/user-agent/internal/detect"
	"github.com/deangrant/user-agent/internal/tokenize"
)

func TestDetectInstagram(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) " +
		"AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 " +
		"Instagram 269.0.0.18.75"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if !state.AppMatched {
		t.Fatal("AppMatched = false, want true")
	}
	if state.AgentName != "Instagram" {
		t.Fatalf("agent = %q, want Instagram", state.AgentName)
	}
	if state.AgentClass != "Mobile App" {
		t.Fatalf("class = %q, want Mobile App", state.AgentClass)
	}
}

func TestDetectSkipsWhenBotMatched(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) " +
		"AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 " +
		"Instagram 269.0.0.18.75"
	state := &detect.State{
		UA:         ua,
		Tokens:     tokenize.Parse(ua),
		BotMatched: true,
		AgentName:  "Googlebot",
		AgentClass: "Robot",
	}
	New().Detect(state)
	if state.AppMatched {
		t.Fatal("AppMatched = true, want false when BotMatched")
	}
	if state.AgentName != "Googlebot" {
		t.Fatalf("agent overwritten to %q", state.AgentName)
	}
}

func TestDetectAndroidWebView(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 13; Pixel 7; wv) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Version/4.0 Chrome/120.0.0.0 Mobile Safari/537.36"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if !state.AppMatched {
		t.Fatal("AppMatched = false, want true")
	}
	if state.AgentName != "Android WebView" {
		t.Fatalf("agent = %q, want Android WebView", state.AgentName)
	}
	if state.AgentClass != "Browser Webview" {
		t.Fatalf("class = %q, want Browser Webview", state.AgentClass)
	}
}

func TestDetectBrowserOnlyEdgeNotApp(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.AppMatched {
		t.Fatalf("AppMatched = true for Edge browser UA; agent=%#v", state)
	}
}
