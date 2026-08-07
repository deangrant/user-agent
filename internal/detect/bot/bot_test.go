package bot

import (
	"testing"

	"github.com/deangrant/user-agent/internal/detect"
	"github.com/deangrant/user-agent/internal/tokenize"
)

func TestDetectGooglebot(t *testing.T) {
	ua := "Mozilla/5.0 (compatible; Googlebot/2.1; " +
		"+http://www.google.com/bot.html)"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if !state.BotMatched {
		t.Fatal("BotMatched = false, want true")
	}
	if state.AgentName != "Googlebot" {
		t.Fatalf("agent = %q, want Googlebot", state.AgentName)
	}
	if state.AgentClass != "Robot" {
		t.Fatalf("agent class = %q, want Robot", state.AgentClass)
	}
	if state.DeviceClass != "Robot" {
		t.Fatalf("device class = %q, want Robot", state.DeviceClass)
	}
	if state.EngineClass != "Robot" {
		t.Fatalf("engine class = %q, want Robot", state.EngineClass)
	}
	if state.OSClass != "Cloud" {
		t.Fatalf("os class = %q, want Cloud", state.OSClass)
	}
}

func TestDetectBrowserWithBotSubstringNotBot(t *testing.T) {
	// "robot" contains heuristic "bot", but a normal Chrome UA must not match.
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 RobotLab/1.0"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.BotMatched {
		t.Fatalf("BotMatched = true, want false; agent=%#v", state)
	}
}

func TestDetectHeuristicGenericBot(t *testing.T) {
	ua := "SomethingBot/1.0"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if !state.BotMatched {
		t.Fatal("BotMatched = false, want true")
	}
	if state.AgentName != "Generic Bot" {
		t.Fatalf("agent = %q, want Generic Bot", state.AgentName)
	}
}

func TestDetectWhatsAppPreviewIsBot(t *testing.T) {
	ua := "WhatsApp/2.23.25.76 A"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if !state.BotMatched {
		t.Fatal("BotMatched = false, want true")
	}
	if state.AgentName != "WhatsApp" {
		t.Fatalf("agent = %q, want WhatsApp", state.AgentName)
	}
}

func TestDetectChromeWhatsAppNotHeuristicBot(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 13; Pixel 7) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/120.0.0.0 Mobile Safari/537.36 WhatsApp/2.23.25.76"
	state := &detect.State{
		UA:     ua,
		Tokens: tokenize.Parse(ua),
	}
	New().Detect(state)
	if state.BotMatched {
		t.Fatalf("BotMatched = true, want false; agent=%#v", state)
	}
}

func TestDetectEmptyUA(t *testing.T) {
	state := &detect.State{}
	New().Detect(state)
	if state.BotMatched {
		t.Fatal("BotMatched = true on empty UA")
	}
	New().Detect(nil) // must not panic
}
